package main

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
	camera := trace.Camera{Width: 960, Height: 540}
	sampler := trace.NewSampler(&camera, &scene, 10)
	renderer := trace.NewRenderer(&camera)
	light := trace.NewLight(1000, 1000, 1000)
	redPlastic := trace.NewPlastic(1, 0, 0, 1)
	bluePlastic := trace.NewPlastic(0, 0, 1, 1)
	whitePlastic := trace.NewPlastic(1, 1, 1, 0)
	silver := trace.NewMetal(0.972, 0.960, 0.915, 1)
	gold := trace.NewMetal(1.022, 0.782, 0.344, 0.8)
	glass := trace.NewGlass(0, 0, 0, 0, 1)
	frostedGlass := trace.NewGlass(0, 1, 0, 0.05, 0.8)

	scene.SetEnv("images/glacier.hdr", 100)
	scene.Add(trace.Sphere{trace.Vector3{1.5, 0.02, -5}, 1, redPlastic})
	scene.Add(trace.Sphere{trace.Vector3{0, 0.02, -7}, 1, silver})
	scene.Add(trace.Sphere{trace.Vector3{5, 0.02, -7}, 1, gold})
	scene.Add(trace.Sphere{trace.Vector3{-2.25, 0.02, -8}, 1, frostedGlass})
	scene.Add(trace.Sphere{trace.Vector3{-8, 0.02, -9}, 1, bluePlastic})
	scene.Add(trace.Sphere{trace.Vector3{-2.5, 0.02, -4.5}, 1, glass})
	scene.Add(trace.Sphere{trace.Vector3{150, -250, -100}, 150, light})
	scene.Add(trace.Sphere{trace.Vector3{0, 10001, -6}, 10000, whitePlastic})
	// camera.Move(0, 2, 1)
	camera.LookAt(0, 0, -1)

	frameSamples := (*samples) * sampler.Width * sampler.Height
	fmt.Printf("Collecting %v frames, taking %v samples/frame...\n", *frames, frameSamples)
	sampler.Collect(*frames, frameSamples)
	renderer.Write(sampler.Values(), *out)
	if len(*heat) > 0 {
		renderer.Write(sampler.Counts(), *heat)
	}
	fmt.Printf("Done: %v\n", *out)
}
