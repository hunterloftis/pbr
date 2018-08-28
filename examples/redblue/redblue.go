package main

import (
	"fmt"
	"os"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/format/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	}
}

func run() error {
	mesh, err := obj.ReadFile("./fixtures/models/mario/mario-sculpture.obj", true)
	if err != nil {
		return err
	}

	bounds, surfaces := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewFlat(0, 0, 0))

	camera.MoveTo(geom.Vec{100, 100, 400}).LookAt(geom.Vec{0, 0, 0})
	floor := surface.UnitCube(material.Plastic(0.9, 0.9, 0.9, 0.5))
	dims := bounds.Max.Minus(bounds.Min).Scaled(1.1)
	floor.Shift(geom.Vec{bounds.Center.X, bounds.Min.Y - dims.Y*0.25, bounds.Center.Z})
	floor.Scale(geom.Vec{dims.X, dims.Y * 0.5, dims.Z})
	surfaces = append(surfaces, floor)

	red := surface.UnitSphere(material.Light(200000, 10000, 10000))
	red.Shift(geom.Vec{-100, 0, 0}).Scale(geom.Vec{10, 10, 10})
	blue := surface.UnitSphere(material.Light(10000, 10000, 200000))
	blue.Shift(geom.Vec{100, 0, 0}).Scale(geom.Vec{10, 10, 10})
	surfaces = append(surfaces, red, blue)
	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "redblue.png", 1280, 720, 6, true)
}
