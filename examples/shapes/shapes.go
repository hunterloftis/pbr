package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/rgb"
	"github.com/hunterloftis/pbr/surface"
)

func main() {
	light := material.Light(1500, 1500, 1500)
	redPlastic := material.Plastic(0.9, 0, 0, 0)
	whitePlastic := material.Plastic(1, 1, 1, 0.2)
	bluePlastic := material.Plastic(0, 0, 0.9, 0)
	greenPlastic := material.Plastic(0, 0.9, 0, 0)
	gold := material.Metal(1.022, 0.782, 0.344, 0.9, 0)
	greenGlass := material.Glass(0.2, 1, 0.1, 0.05)

	scene := pbr.NewScene()
	scene.SetAmbient(rgb.Energy{40, 50, 60})

	cam := pbr.NewCamera(888, 500).MoveTo(-0.6, 0.12, 0.8)
	cam.LookAt(geom.Vector3{}, geom.Vector3{0, -0.025, 0.2})

	render := pbr.NewRender(scene, cam)

	scene.Add(
		surface.UnitCube(redPlastic).Rotate(0, -0.25*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(gold).Move(0, 0, -0.4).Rotate(0, 0.1*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(greenGlass).Move(-0.3, 0, 0.3).Rotate(0, -0.1*math.Pi, 0).Scale(0.1, 0.1, 0.1),
		surface.UnitCube(greenGlass).Move(0.175, 0.05, 0.18).Rotate(0, 0.55*math.Pi, 0).Scale(0.02, 0.2, 0.2),
		surface.UnitCube(whitePlastic).Move(0, -0.55, 0).Scale(1000, 1, 1000),
		surface.UnitSphere(greenGlass).Move(-0.2, 0.001, -0.2).Scale(0.1, 0.1, 0.1),
		surface.UnitSphere(bluePlastic).Move(0.3, 0.05, 0).Scale(0.2, 0.2, 0.2),
		surface.UnitSphere(light).Move(7, 30, 6).Scale(30, 30, 30),
		surface.UnitSphere(greenPlastic).Move(0, -0.025, 0.2).Scale(0.1, 0.05, 0.1),
		surface.UnitSphere(gold).Move(0.45, 0.05, -0.4).Scale(0.2, 0.2, 0.2),
	)

	fmt.Println("rendering shapes.png (15 mins)...")
	render.Start()
	time.Sleep(time.Minute * 15)
	render.Stop()
	render.WritePngs("shapes.png", "shapes-heat.png", "shapes-noise.png", 1)
}
