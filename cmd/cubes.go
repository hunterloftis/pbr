package main

// TODO: tweet for feedback
// https://arschles.svbtle.com/orthogonality-in-go

import (
	"flag"
	"fmt"

	"github.com/hunterloftis/trace/trace"
)

func main() {
	out := flag.String("out", "trace.png", "Output png filename.")
	frames := flag.Int("frames", 4, "Number of frames to combine.")
	samples := flag.Int("samples", 4, "Average per pixel samples to take.")
	heat := flag.String("heat", "", "Heatmap png filename.")
	flag.Parse()

	scene := trace.Scene{}
	camera := trace.NewCamera(960, 540, 0.050)
	sampler := trace.NewSampler(&camera, &scene, 10)
	renderer := trace.NewRenderer(&camera)

	light := trace.NewLight(1000, 1000, 1000)
	redPlastic := trace.NewPlastic(1, 0, 0, 1)
	// bluePlastic := trace.NewPlastic(0, 0, 1, 1)
	whitePlastic := trace.NewPlastic(1, 1, 1, 0)
	// silver := trace.NewMetal(0.972, 0.960, 0.915, 1)
	// gold := trace.NewMetal(1.022, 0.782, 0.344, 0.8)
	// glass := trace.NewGlass(0, 0, 0, 0, 1)
	// frostedGlass := trace.NewGlass(0, 1, 0, 0.05, 0.8)

	scene.Add(&trace.Cube{
		Pos: trace.Identity(), //.Trans(0.13, 0, -3).Scale(0.1, 0.1, 0.1),
		Mat: redPlastic,
	})
	scene.Add(&trace.Sphere{trace.Vector3{0.13, 0, -3}, 0.1, redPlastic})
	scene.Add(&trace.Sphere{trace.Vector3{15.0, 25.0, -10.0}, 15.0, light})
	scene.Add(&trace.Sphere{trace.Vector3{0, -10000.1, -4}, 10000, whitePlastic})
	camera.Move(0, 0.15, -1.5)
	camera.LookAt(0, 0, -4)
	camera.Focus(0.13, 0, -3, 4)

	frameSamples := (*samples) * sampler.Width * sampler.Height
	fmt.Printf("Collecting %v frames, taking %v samples/frame...\n", *frames, frameSamples)
	sampler.Collect(*frames, frameSamples)
	renderer.Write(sampler.Values(), *out)
	if len(*heat) > 0 {
		renderer.Write(sampler.Counts(), *heat)
	}
	fmt.Printf("Done: %v\n", *out)
}
