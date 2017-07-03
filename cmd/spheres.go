package main

import (
	"flag"
	"fmt"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	out := flag.String("out", "render.png", "Output png filename.")
	frames := flag.Int("frames", 4, "Number of frames to combine.")
	samples := flag.Int("samples", 4, "Average per pixel samples to take.")
	heat := flag.String("heat", "", "Heatmap png filename.")
	flag.Parse()

	scene := pbr.Scene{}
	camera := pbr.NewCamera(960, 540, 0.050)
	sampler := pbr.NewSampler(&camera, &scene, 10)
	renderer := pbr.NewRenderer(&camera)
	light := pbr.NewLight(1000, 1000, 1000)
	redPlastic := pbr.NewPlastic(1, 0, 0, 1)
	bluePlastic := pbr.NewPlastic(0, 0, 1, 1)
	whitePlastic := pbr.NewPlastic(1, 1, 1, 0)
	silver := pbr.NewMetal(0.972, 0.960, 0.915, 1)
	gold := pbr.NewMetal(1.022, 0.782, 0.344, 0.8)
	glass := pbr.NewGlass(0, 0, 0, 0, 1)
	frostedGlass := pbr.NewGlass(0, 1, 0, 0.05, 0.8)

	// scene.SetEnv("images/glacier.hdr", 100)
	scene.Add(&pbr.Sphere{pbr.Vector3{-0.02, 0, -3.4}, 0.1, silver})
	scene.Add(&pbr.Sphere{pbr.Vector3{0.13, 0, -3}, 0.1, redPlastic})
	scene.Add(&pbr.Sphere{pbr.Vector3{0.42, 0, -3.5}, 0.1, gold})
	scene.Add(&pbr.Sphere{pbr.Vector3{-0.26, 0, -3.9}, 0.1, frostedGlass})
	scene.Add(&pbr.Sphere{pbr.Vector3{-0.27, 0, -2.9}, 0.1, glass})
	scene.Add(&pbr.Sphere{pbr.Vector3{-0.8, 0, -4}, 0.1, bluePlastic})
	scene.Add(&pbr.Sphere{pbr.Vector3{15.0, 25.0, -10.0}, 15.0, light})
	scene.Add(&pbr.Sphere{pbr.Vector3{0, -10000.1, -4}, 10000, whitePlastic})
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
