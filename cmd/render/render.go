package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	var position, target, focus *pbr.Vector3
	sky := pbr.Vector3{40, 50, 60}
	in := os.Args[1]
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	quality := flag.Float64("quality", math.Inf(1), "Minimum samples-per-pixel to reach before exiting")
	adapt := flag.Int("adapt", 4, "Adaptive sampling; 0=off, 3=medium, 5=high") // TODO: 0 is broken
	bounces := flag.Int("bounces", 10, "Maximum light bounces")
	profile := flag.Bool("profile", false, "Record performance into profile.pprof")
	pano := flag.String("pano", "", "Panoramic environment map hdr (radiosity) file")
	lens := flag.Float64("lens", 50, "Camera focal length in mm")
	exposure := flag.Float64("exposure", 1, "Exposure multiplier")
	fStop := flag.Float64("fstop", 4, "Camera f-stop")
	flag.Var(position, "position", "Camera position")
	flag.Var(target, "target", "Camera target location")
	flag.Var(focus, "focus", "Camera focus location")
	flag.Var(&sky, "sky", "Ambient sky lighting")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [options] <scene.dae>\n\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println()
	}
	flag.Parse()

	if *profile {
		f, _ := os.Create("profile.pprof")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	xml, _ := ioutil.ReadFile(in)
	scene := pbr.ColladaScene(xml)
	camera := pbr.NewCamera(1280, 720, pbr.CameraConfig{
		Lens:     (*lens) / 1000.0,
		Position: position,
		Target:   target,
		Focus:    focus,
		FStop:    *fStop,
	})
	sampler := pbr.NewSampler(camera, scene, pbr.SamplerConfig{
		Bounces: *bounces,
		Adapt:   *adapt,
	})
	renderer := pbr.NewRenderer(sampler, pbr.RenderConfig{
		Exposure: *exposure,
	})

	scene.SetSky(sky, pbr.Vector3{})
	if len(*pano) > 0 {
		hdr, _ := os.Open(*pano)
		defer hdr.Close()
		scene.SetPano(hdr, 100) // TODO: read radiosity info or allow it from the command line
	}

	// For debugging until we're actually parsing collada files
	scene.Add(pbr.UnitCube(pbr.Plastic(1, 0, 0, 1), pbr.Rot(pbr.Vector3{0, 1, 0}), pbr.Scale(0.5, 0.5, 0.5)))

	start := time.Now()
	running := true
	interrupt := make(chan os.Signal, 2)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // TODO: abstract this?
	go func() {
		<-interrupt
		running = false
		pbr.ShowProgress(sampler, start, running)
	}()

	pbr.ShowProgress(sampler, start, running)
	for running && sampler.PerPixel() < *quality {
		sampler.Sample()
		pbr.ShowProgress(sampler, start, running)
	}

	pbr.WritePNG(*out, renderer.Rgb())
	if len(*heat) > 0 {
		pbr.WritePNG(*heat, renderer.Heat())
	}
}
