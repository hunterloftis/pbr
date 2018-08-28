package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hunterloftis/pbr2/pkg/camera"
	"github.com/hunterloftis/pbr2/pkg/env"
	"github.com/hunterloftis/pbr2/pkg/format/obj"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

var materials = map[string]surface.Material{
	"gold":    material.Gold(0.03, 0.9),
	"mirror":  material.Mirror(0.001),
	"glass":   material.Glass(0.03),
	"flat":    material.Plastic(1, 1, 1, 0.5),
	"plastic": material.Plastic(1, 1, 1, 0.1),
}

func main() {
	if err := run(options()); err != nil {
		printErr(err)
		os.Exit(1)
	}
}

func run(o *Options) error {
	mesh, err := obj.ReadFile(o.Scene, true)
	if err != nil {
		return err
	}

	if o.Scale != nil {
		mesh.Scale(*o.Scale)
	}
	if o.Rotate != nil {
		mesh.Rotate(*o.Rotate)
	}
	if o.Material != "" {
		m := materials[strings.ToLower(o.Material)]
		mesh.SetMaterial(m)
	}
	bounds, surfaces := mesh.Bounds()
	camera := camera.NewSLR()
	environment := render.Environment(env.NewGradient(rgb.Black, *o.Ambient, 3))

	o.SetDefaults(bounds)
	camera.MoveTo(*o.From).LookAt(*o.To)
	camera.Lens = o.Lens / 1000
	camera.FStop = o.FStop
	camera.Focus = o.Focus

	if o.Verbose || o.Info {
		printInfo(bounds, len(surfaces), camera)
		if o.Info {
			return nil
		}
	}

	if o.Env != "" {
		environment, err = env.ReadFile(o.Env, o.Rad)
		if err != nil {
			return err
		}
	}

	if o.Floor > 0 {
		floor := surface.UnitCube(material.Plastic(o.FloorColor.X, o.FloorColor.Y, o.FloorColor.Z, o.FloorRough))
		dims := bounds.Max.Minus(bounds.Min).Scaled(o.Floor)
		floor.Shift(geom.Vec{bounds.Center.X, bounds.Min.Y - dims.Y*0.5, bounds.Center.Z})
		floor.Scale(geom.Vec{dims.X, dims.Y, dims.Z})
		surfaces = append(surfaces, floor)
	}

	if o.Sun != nil {
		sun := surface.UnitSphere(material.Daylight(1000000))
		sun.Shift(*o.Sun).Scale(geom.Vec{o.SunSize, o.SunSize, o.SunSize})
		surfaces = append(surfaces, sun)
	}

	tree := surface.NewTree(surfaces...)
	scene := render.NewScene(camera, tree, environment)

	fmt.Println("Surfaces:", len(surfaces))
	return render.Iterative(scene, o.Out, o.Width, o.Height, o.Bounce, !o.Indirect)
}
