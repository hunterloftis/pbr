package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	scene := trace.Scene{}
	camera := trace.Camera{}
	sampler := trace.NewSampler(960, 540, camera, scene, 10)
	renderer := trace.Renderer{Width: 960, Height: 540}

	scene.Add(trace.Sphere{trace.Vector3{0, 0, -10}, 1})

	for i := 0; i < 10; i++ {
		sampler.Sample()
	}
	renderer.Write(sampler.Values(), "test.png")
}
