package main

import (
	"math"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	scene := pbr.EmptyScene()
	camera := pbr.Camera35mm(960, 540, 0.055)
	renderer := pbr.CamRenderer(camera, 1)
	cli := pbr.CliRunner(scene, camera, renderer)

	light := pbr.Light(1500, 1500, 1500)
	redPlastic := pbr.Plastic(1, 0, 0, 1)
	whitePlastic := pbr.Plastic(1, 1, 1, 0.8)
	// grayLambert := pbr.Lambert(0.2, 0.2, 0.2, 0.1)
	// bluePlastic := pbr.Plastic(0, 0, 1, 1)
	silver := pbr.Metal(0.95, 0.93, 0.88, 1)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	// glass := pbr.Glass(1, 1, 1, 1)
	greenGlass := pbr.Glass(0.8, 1, 0.7, 0.95)

	scene.SetSky(pbr.Vector3{30, 40, 50}, pbr.Vector3{})

	scene.Add(
		pbr.UnitCube(pbr.Ident().Rot(pbr.Vector3{0, -0.25 * math.Pi, 0}), redPlastic),
		pbr.UnitCube(pbr.Ident().Trans(0, 0, -4).Rot(pbr.Vector3{0, 0.1 * math.Pi, 0}), gold),
		pbr.UnitCube(pbr.Ident().Trans(-2, 0, 1.25).Rot(pbr.Vector3{0, -0.1 * math.Pi, 0}), greenGlass),
		pbr.UnitCube(pbr.Ident().Trans(1.75, 0.5, 1.75).Rot(pbr.Vector3{0, 0.55 * math.Pi, 0}).Scale(0.1, 2, 1), greenGlass),
		pbr.UnitCube(pbr.Ident().Trans(0, -1, 0).Scale(1000, 1, 1000), whitePlastic),
		pbr.UnitSphere(pbr.Ident().Trans(-2, 0.01, -2), greenGlass),
		pbr.UnitSphere(pbr.Ident().Trans(3, 0.5, 0).Scale(2, 2, 2), greenGlass),
		pbr.UnitSphere(pbr.Ident().Trans(70, 300, 60).Scale(300, 300, 300), light),
		pbr.UnitSphere(pbr.Ident().Trans(0, -0.25, 2).Scale(1, 0.5, 1), greenGlass),
		pbr.UnitSphere(pbr.Ident().Trans(4.5, 0.5, -4).Scale(2, 2, 2), silver),
	)

	camera.MoveTo(-6, 1.5, 8)
	camera.LookAt(0, 0, 0)
	camera.Focus(0, 0, 0, 1.4)

	cli.Render()
}
