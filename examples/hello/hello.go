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
	wall := surface.UnitCube(material.Default).Move(0, 0, -2).Scale(10, 10, 1)
	halogen := material.Light(4781, 4518, 4200)
	light := surface.UnitSphere(halogen).Move(-5, 8, 6).Scale(5, 5, 5)
	ball := surface.UnitSphere(material.Chrome)
	scene := pbr.NewScene(floor, light, ball, wall)
	cam := pbr.NewCamera(888, 500).MoveTo(0, 0.5, 5).LookAt(ball.Center(), ball.Center())
	render := pbr.NewRender(scene, cam)

	fmt.Println("rendering hello.png (3 minutes)...")
	render.Start()
	time.Sleep(time.Minute * 1)
	render.Stop()
	render.WritePngs("hello.png", "hello-heat.png", "hello-noise.png", 1)
}
