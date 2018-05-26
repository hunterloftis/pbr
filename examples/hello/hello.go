package main

import (
	"fmt"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/surface"
)

func main() {
	floor := surface.UnitCube(material.Default).Move(0, -1, 0).Scale(100, 1, 100)
	wall := surface.UnitCube(material.Default).Move(0, 0, -2).Scale(100, 100, 1)
	halogen := material.Halogen(2500)
	light := surface.UnitSphere(halogen).Move(0, 30, 15).Scale(30, 30, 30)
	ball := surface.UnitSphere(material.Copper)
	ball2 := surface.UnitSphere(material.RedPlastic).Move(1, 0, 0)
	scene := pbr.NewScene(floor, light, ball, wall, ball2)
	cam := pbr.NewCamera(700, 500).MoveTo(0, 0.5, 5).LookAt(ball.Center(), ball.Center())
	render := pbr.NewRender(scene, cam)

	fmt.Println("rendering hello.png (3 minutes)...")
	render.Start()
	time.Sleep(time.Minute * 5)
	render.Stop()
	render.WritePngs("hello.png", "hello-heat.png", "hello-noise.png", 1)
}
