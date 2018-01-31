package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/obj"
	"github.com/hunterloftis/pbr/surface"
)

func main() {
	moses, err := obj.ReadFile("fixtures/models/moses/model.obj", false)
	if err != nil {
		panic(err)
	}
	key := surface.UnitSphere(material.Light(100000, 100000, 50000)).Move(-20, 10, 20).Scale(5, 5, 5)
	fill := surface.UnitSphere(material.Light(20000, 20000, 50000)).Move(30, 10, 5).Scale(5, 5, 5)
	back := surface.UnitSphere(material.Light(25000, 25000, 100000)).Move(-30, -5, -10).Scale(8, 8, 8)
	scene := pbr.NewScene(moses...)
	bounds, _ := scene.Info()
	target := bounds.Center
	scene.Add(key, fill, back)
	cam := pbr.NewCamera(888, 500).MoveTo(0, -10, 50).LookAt(target, target)
	render := pbr.NewRender(scene, cam)
	interrupt := make(chan os.Signal, 2)

	fmt.Println("rendering moses.png (press Ctrl+C to finish)...")
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	render.Start()
	<-interrupt
	render.Stop()
	render.WritePngs("moses.png", "moses-heat.png", "moses-noise.png", 1)
}
