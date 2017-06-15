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
	lambert := trace.Material{Color: trace.Vector3{1, 1, 1}, Gloss: 0}

	scene.SetEnv("images/ennis.hdr", 100)
	scene.Add(trace.Sphere{trace.Vector3{-1.5, 0, -5}, 1, light})
	scene.Add(trace.Sphere{trace.Vector3{1.5, 0, -5}, 1, lambert})

	for i := 0; i < 10; i++ {
		sampler.Sample()
	}
	renderer.Write(sampler.Values(), "test.png")
}
