package main

// TODO: tweet for feedback
// https://arschles.svbtle.com/orthogonality-in-go

import (
	"math"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	// f, _ := os.Create("profile.pprof")
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	scene := pbr.Scene{}
	camera := pbr.Camera35mm(960, 540, 0.050)
	renderer := pbr.CamRenderer(&camera, 2)
	cli := pbr.NewCLI(&scene, &camera, renderer)

	light := pbr.Light(1000, 1000, 1000)
	redPlastic := pbr.Plastic(1, 0, 0, 1)
	whitePlastic := pbr.Plastic(1, 1, 1, 0)
	bluePlastic := pbr.Plastic(0, 0, 1, 1)
	silver := pbr.Metal(0.972, 0.960, 0.915, 1)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	glass := pbr.Glass(0, 0, 0, 0, 1)

	scene.SetEnv("images/glacier.hdr", 100)

	scene.Add(&pbr.Cube{
		Pos: pbr.Identity(),
		Mat: silver,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(0, 0, -4),
		Mat: gold,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(-2, -0.245, 1.25).Scale(1, 0.5, 1),
		Mat: glass,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(1.75, 0, 2).Rot(pbr.Vector3{0, 0.25 * math.Pi, 0}).Scale(0.5, 2, 1),
		Mat: redPlastic,
	})
	scene.Add(&pbr.Cube{
		Pos: pbr.Identity().Trans(0, -1, 0).Scale(1000, 1, 1000),
		Mat: whitePlastic,
	})
	scene.Add(&pbr.Sphere{pbr.Vector3{-1.5, 0, 0}, 0.5, bluePlastic})
	scene.Add(&pbr.Sphere{pbr.Vector3{-4.5, 0, 1.5}, 0.5, gold})
	scene.Add(&pbr.Sphere{pbr.Vector3{2, 0, 0}, 0.5, silver})
	scene.Add(&pbr.Sphere{pbr.Vector3{50, 30, 0}, 20, light})

	camera.MoveTo(-6, 1.5, 5)
	camera.LookAt(0, 0, 0)
	camera.Focus(0, 0, 0, 1.4)

	cli.Start()
}
