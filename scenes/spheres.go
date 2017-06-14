package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	scene := trace.Scene{}
	camera := trace.Camera{Width: 960, Height: 540}
	sampler := trace.NewSampler(&camera, &scene, 10)
	renderer := trace.NewRenderer(&camera)

	scene.Add(trace.Sphere{trace.Vector3{0, 0, -4}, 1})

	for i := 0; i < 10; i++ {
		sampler.Sample()
	}
	renderer.Write(sampler.Values(), "test.png")
}
