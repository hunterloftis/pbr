package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/surface"
)

// The pathological case for indirect lighting: a small, very bright light at a large distance
func main() {
	floor := surface.UnitCube().Move(0, -1, 0).Scale(100, 1, 100)
	halogen := material.Light(10000000, 10000000, 10000000)
	light := surface.UnitSphere(halogen).Move(-50, 100, -10)
	box := surface.UnitSphere(material.Gold)
	scene := pbr.NewScene(floor, light, box)
	cam := pbr.NewCamera(888, 500).MoveTo(-3, 2, 5).LookAt(box.Center(), box.Center())
	render := pbr.NewRender(scene, cam)
	interrupt := make(chan os.Signal, 2)

	render.SetDirect(1)

	fmt.Println("rendering hello.png (press Ctrl+C to finish)...")
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	render.Start()
	<-interrupt
	render.Stop()
	render.WritePngs("hello.png", "hello-heat.png", "hello-noise.png", 1)
}
