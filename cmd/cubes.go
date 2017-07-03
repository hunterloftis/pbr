package main

// TODO: tweet for feedback
// https://arschles.svbtle.com/orthogonality-in-go

import (
	"flag"
	"fmt"
	"math"

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
	whitePlastic := pbr.NewPlastic(1, 1, 1, 0)
	silver := pbr.NewMetal(0.972, 0.960, 0.915, 1)
	gold := pbr.NewMetal(1.022, 0.782, 0.344, 0.9)
	glass := pbr.NewGlass(0, 0, 0, 0, 1)

	scene.SetEnv("images/glacier.hdr", 100)
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity(),
		Mat: silver,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(2, 0, -4),
		Mat: gold,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(-1.5, -0.25, 1.25).Scale(1, 0.5, 1),
		Mat: glass,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(1.75, 0, 2).Rotate(0, 0.25*math.Pi, 0).Scale(0.5, 2, 1),
		Mat: redPlastic,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(0, -1, 0).Scale(1000, 1, 1000),
		Mat: whitePlastic,
	})
	scene.Add(&pbr.Sphere{pbr.Vector3{-1.5, 0, 0}, 0.5, gold})
	scene.Add(&pbr.Sphere{pbr.Vector3{2, 0, 0}, 0.5, silver})
	scene.Add(&pbr.Sphere{pbr.Vector3{50, 30, 0}, 20, light})
	camera.Move(-6, 2.5, 5)
	camera.LookAt(0, 0, 0)
	camera.Focus(0, 0, 0, 4)

	frameSamples := (*samples) * sampler.Width * sampler.Height
	fmt.Printf("Collecting %v frames, taking %v samples/frame...\n", *frames, frameSamples)
	sampler.Collect(*frames, frameSamples)
	renderer.Write(sampler.Values(), *out)
	if len(*heat) > 0 {
		renderer.Write(sampler.Counts(), *heat)
	}
	fmt.Printf("Done: %v\n", *out)
}
