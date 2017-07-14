package main

import (
	"runtime"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	workers := runtime.NumCPU()
	scene := pbr.EmptyScene()
	camera := pbr.NewCamera(1280, 720)
	renderer := pbr.CamRenderer(camera)
	monitor := pbr.NewMonitor(10)

	scene.SetSky(pbr.Vector3{40, 50, 60}, pbr.Vector3{})
	scene.Add(pbr.UnitSphere(pbr.Ident(), pbr.Plastic(1, 1, 1, 0.8)))

	for i := 0; i < workers; i++ {
		monitor.AddSampler(pbr.NewSampler(camera, scene))
	}
	for i := 0; i < workers; i++ {
		renderer.Merge(<-m.Results)
	}
	pbr.WritePNG("hello.png", renderer.Rgb())
}
