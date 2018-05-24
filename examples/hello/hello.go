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
	floor := surface.UnitCube(material.Default).Move(0, -1, 0).Scale(5, 1, 5)
	wall := surface.UnitCube(material.Default).Move(0, 0, -1.5).Scale(5, 5, 1)
	wall2 := surface.UnitCube(material.Default).Move(1.5, 0, 0).Scale(1, 5, 5)
	halogen := material.Halogen(5000)
	light := surface.UnitSphere(halogen).Move(-5, 8, 6).Scale(5, 5, 5)
	ball := surface.UnitSphere(material.Default)
	scene := pbr.NewScene(floor, light, ball, wall, wall2)
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
