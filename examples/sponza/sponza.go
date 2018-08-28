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
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	}
}

func run() error {
	mesh, err := obj.ReadFile("./fixtures/models/sponza/sponza.obj", true)
	if err != nil {
		return err
	}

	bounds, surfaces := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, rgb.Energy{4000, 4000, 4000}, 3))

	camera.MoveTo(geom.Vec{1140, 620, -160}).LookAt(geom.Vec{1090, 608, -150})
	camera.Lens = 0.028
	floor := surface.UnitCube(material.Plastic(0.9, 0.9, 0.9, 0.5))
	dims := bounds.Max.Minus(bounds.Min).Scaled(1.1)
	floor.Shift(geom.Vec{bounds.Center.X, bounds.Min.Y - dims.Y*0.25, bounds.Center.Z})
	floor.Scale(geom.Vec{dims.X, dims.Y * 0.5, dims.Z})
	surfaces = append(surfaces, floor)

	sun := surface.UnitSphere(material.Daylight(800000))
	sun.Shift(geom.Vec{1300, 5000, -600}).Scale(geom.Vec{400, 400, 400})
	surfaces = append(surfaces, sun)
	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	return render.Iterative(scene, "sponza.png", 1280, 720, 8, true)
}
