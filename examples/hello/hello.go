package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/surface"
)

func main() {
	// for a := 0.0; a < math.Pi; a += 0.01 {
	// 	x := math.Cos(a)
	// 	y := math.Sin(a)
	// 	z := 0.0
	// 	in := geom.Vector3{x, y, z}.Unit()
	// 	s := fresnelSchlick(in, geom.Up, 0.04)
	// 	fmt.Println("angle:", math.Acos(in.Dot(geom.Up)), "schlick:", s)
	// }
	floor := surface.UnitCube(material.Default).Move(0, -1, 0).Scale(100, 1, 100)
	wall := surface.UnitCube(material.Default).Move(0, 0, -2).Scale(100, 100, 1)
	halogen := material.Halogen(1500)
	light := surface.UnitSphere(halogen).Move(0, 30, 15).Scale(30, 30, 30)
	ball := surface.UnitSphere(material.Copper)
	scene := pbr.NewScene(floor, light, ball, wall)
	cam := pbr.NewCamera(500, 500).MoveTo(0, 0.5, 5).LookAt(ball.Center(), ball.Center())
	render := pbr.NewRender(scene, cam)

	fmt.Println("rendering hello.png (3 minutes)...")
	render.Start()
	time.Sleep(time.Minute * 1)
	render.Stop()
	render.WritePngs("hello.png", "hello-heat.png", "hello-noise.png", 1)
}

func fresnelSchlick(in, normal geom.Direction, f0 float64) float64 {
	return f0 + (1-f0)*math.Pow(1-normal.Dot(in), 5)
}
