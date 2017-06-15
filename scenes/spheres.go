package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	scene := trace.Scene{}
	camera := trace.Camera{Width: 960, Height: 540}
	sampler := trace.NewSampler(&camera, &scene, 10)
	renderer := trace.NewRenderer(&camera)
	light := trace.Material{Light: trace.Vector3{500, 500, 500}}

	scene.SetEnv("images/ennis.hdr", 100)
	scene.Add(trace.Sphere{trace.Vector3{0, 0, -4}, 1, light})

	for i := 0; i < 10; i++ {
		sampler.Sample()
	}
	renderer.Write(sampler.Values(), "test.png")
}
