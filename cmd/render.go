package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"runtime"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	in := os.Args[1]
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	workers := flag.Int("workers", runtime.NumCPU(), "Concurrency level")
	samples := flag.Float64("samples", math.Inf(1), "Max samples per pixel")
	adapt := flag.Int("adapt", 4, "Adaptive sampling; 0=off, 3=medium, 5=high")
	bounces := flag.Int("bounces", 10, "Maximum light bounces")
	profile := flag.Bool("profile", false, "Record performance into profile.pprof")
	sky := flag.String("sky", "40,50,60", "Ambient sky lighting RGB")
	pano := flag.String("pano", "", "Panoramic environment map hdr (radiosity) file")
	flag.Parse()

	xml, _ := ioutil.ReadFile(in)
	scene := pbr.ColladaScene(xml)
	camera := pbr.Camera35mm(1280, 720, 0.050)
	renderer := pbr.CamRenderer(camera, 1)
	monitor := pbr.Monitor(renderer, workers, samples)

	camera.MoveTo(-scene.Scale(), scene.Scale(), -scene.Scale())
	camera.LookAt(0, 0, 0)
	camera.Focus(0, 0, 0, 4)
	monitor.Start()

	for pp := range <-monitor.C {
		fmt.Printf("%v samples / pixel\n", pp)
	}
	writePNG(*out, c.renderer.Rgb())
	fmt.Printf("-> %v\n", *out)
	if len(*heat) > 0 {
		writePNG(*heat, c.renderer.Heat())
		fmt.Printf("-> %v\n", *heat)
	}
}
