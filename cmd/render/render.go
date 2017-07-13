package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	var target, focus *pbr.Vector3
	position := pbr.Vector3{0, 0, 0}
	sky := pbr.Vector3{40, 50, 60}
	in := os.Args[1]
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	workers := flag.Int("workers", runtime.NumCPU(), "Concurrency level")
	quality := flag.Float64("quality", math.Inf(1), "Minimum samples-per-pixel to reach before exiting")
	adapt := flag.Int("adapt", 4, "Adaptive sampling; 0=off, 3=medium, 5=high")
	bounces := flag.Int("bounces", 10, "Maximum light bounces")
	profile := flag.Bool("profile", false, "Record performance into profile.pprof")
	pano := flag.String("pano", "", "Panoramic environment map hdr (radiosity) file")
	lens := flag.Float64("lens", 50, "Camera focal length in mm")
	exposure := flag.Float64("exposure", 1, "Exposure multiplier")
	fStop := flag.Float64("fstop", 4, "Camera f-stop")
	flag.Var(&position, "position", "Camera position")
	flag.Var(target, "target", "Camera target location")
	flag.Var(focus, "focus", "Camera focus location")
	flag.Var(&sky, "sky", "Ambient sky lighting")
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
	renderer := pbr.CamRenderer(camera, pbr.RenderConfig{
		Exposure: *exposure,
	})

	scene.SetSky(sky, pbr.Vector3{})
	if len(*pano) > 0 {
		hdr, _ := os.Open(*pano)
		defer hdr.Close()
		scene.SetPano(hdr, 100) // TODO: read radiosity info or allow it from the command line
	}

	redPlastic := pbr.Plastic(1, 0, 0, 1)
	scene.Add(
		pbr.UnitCube(pbr.Ident().Trans(0, 0, -3).Rot(pbr.Vector3{0, 1, 0}).Scale(0.25, 0.25, 0.25), redPlastic),
	)

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	m := pbr.NewMonitor()
	for i := 0; i < *workers; i++ {
		m.AddSampler(pbr.NewSampler(camera, scene, pbr.SamplerConfig{
			Bounces: *bounces,
			Adapt:   *adapt,
		}))
	}

	for m.Active() > 0 {
		select {
		case samples := <-m.Progress:
			perPixel := samples / camera.Pixels()
			showProgress(perPixel, m.Stopped())
			if float64(perPixel) >= *quality {
				m.Stop()
			}
		case r := <-m.Results: // TODO: move responsibility for merging results into Monitor?
			renderer.Merge(r)
		case <-interrupt:
			m.Stop()
		default:
		}
	}

	writePNG(*out, renderer.Rgb())
	showFile(*out)
	if len(*heat) > 0 {
		writePNG(*heat, renderer.Heat())
		showFile(*heat)
	}
}

func writePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, i)
}

func showProgress(perPixel int, stopped bool) {
	note := ""
	if stopped {
		note = " (wrapping up...)"
	}
	fmt.Printf("\r%v samples per pixel%v", perPixel, note) // https://stackoverflow.com/a/15442704/1911432
}

func showFile(file string) {
	fmt.Printf("\n-> %v\n", file)
}
