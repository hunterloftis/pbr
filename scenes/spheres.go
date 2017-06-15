package main

import (
	"github.com/hunterloftis/trace/trace"
)

func main() {
	scene := trace.Scene{}
	camera := trace.Camera{Width: 960, Height: 540}
	sampler := trace.NewSampler(&camera, &scene, 10)
	renderer := trace.NewRenderer(&camera)
	light := trace.NewLight(1000, 1000, 1000)
	redPlastic := trace.NewPlastic(1, 0, 0, 1)
	bluePlastic := trace.NewPlastic(0, 0, 1, 1)
	whitePlastic := trace.NewPlastic(1, 1, 1, 0)
	silver := trace.NewMetal(0.972, 0.960, 0.915, 1)
	glass := trace.NewGlass(0, 1, 0, 0.05, 1)

	scene.SetEnv("images/glacier.hdr", 100)
	scene.Add(trace.Sphere{trace.Vector3{1.5, 0, -5}, 1, redPlastic})
	scene.Add(trace.Sphere{trace.Vector3{0, 0, -7}, 1, silver})
	scene.Add(trace.Sphere{trace.Vector3{-2.5, 0, -9}, 1, bluePlastic})
	scene.Add(trace.Sphere{trace.Vector3{-2, 0, -4}, 1, glass})
	scene.Add(trace.Sphere{trace.Vector3{100, -150, -50}, 100, light})
	scene.Add(trace.Sphere{trace.Vector3{0, 10001, -6}, 10000, whitePlastic})

	for i := 0; i < 10; i++ {
		sampler.Sample()
	}
	renderer.Write(sampler.Values(), "test.png")
}
