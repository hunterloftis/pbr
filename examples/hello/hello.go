package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	floor := surface.UnitCube(material.Plastic(1, 1, 1, 0.05))
	floor.Shift(geom.Vec{0, -0.1, 0}).Scale(geom.Vec{10, 0.1, 10})
	ball := surface.UnitSphere(material.Gold(0.05, 1))
	ball.Scale(geom.Vec{0.1, 0.1, 0.1})

	c := camera.NewSLR().MoveTo(geom.Vec{0, 0, -0.5}).LookAt(geom.Vec{0, 0, 0})
	s := surface.NewList(ball, floor)
	e := env.NewGradient(rgb.Black, rgb.Energy{750, 750, 750}, 7)

	scene := render.NewScene(c, s, e)
	err := render.Iterative(scene, "hello.png", 898, 450, 8, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	}
}
