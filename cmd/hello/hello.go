package main

import (
	"fmt"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	scene := pbr.EmptyScene()
	camera := pbr.NewCamera(1280, 720)
	sampler := pbr.NewSampler(camera, scene)
	renderer := pbr.NewRenderer(sampler)

	scene.SetSky(pbr.Vector3{40, 50, 60}, pbr.Vector3{})
	scene.Add(pbr.UnitSphere(pbr.Ident().Trans(0, 0, -3), pbr.Plastic(1, 1, 1, 0.8))) // TODO: (mat *Material, transforms ...*Matrix)

	for sampler.PerPixel() < 16 {
		sampler.SampleFrame()
		fmt.Printf("%.1f samples / pixel\n", sampler.PerPixel()) // TODO: add newline to output.go equivalent
	}
	pbr.WritePNG("hello.png", renderer.Rgb())
}
