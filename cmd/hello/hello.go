package main

import (
	"fmt"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	scene := pbr.EmptyScene()
	camera := pbr.NewCamera(960, 540)
	sampler := pbr.NewSampler(camera, scene)
	renderer := pbr.NewRenderer(sampler)

	scene.SetSky(pbr.Vector3{256, 256, 256}, pbr.Vector3{})
	scene.Add(pbr.UnitSphere(pbr.Plastic(1, 0, 0, 1), pbr.Trans(0, 0, -3))) // TODO: (mat *Material, transforms ...*Matrix)

	for sampler.PerPixel() < 16 {
		sampler.Sample()
		fmt.Printf("\r%.1f samples / pixel", sampler.PerPixel())
	}
	pbr.WritePNG("hello.png", renderer.Rgb())
}
