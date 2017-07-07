package main

import (
	"math"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	scene := pbr.EmptyScene()
	camera := pbr.Camera35mm(960, 540, 0.050)
	renderer := pbr.CamRenderer(camera, 2)
	cli := pbr.CliRunner(scene, camera, renderer)

	light := pbr.Light(1000, 1000, 1000)
	redPlastic := pbr.Plastic(1, 0, 0, 1)
	whitePlastic := pbr.Plastic(1, 1, 1, 1)
	whiteLambert := pbr.Plastic(1, 1, 1, 0)
	// bluePlastic := pbr.Plastic(0, 0, 1, 1)
	// silver := pbr.Metal(0.972, 0.960, 0.915, 1)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	// glass := pbr.Glass(1, 1, 1, 1)
	greenGlass := pbr.Glass(0.8, 1, 0.8, 0.95)

	scene.SetEnv("images/glacier.hdr", 100)

	scene.Add(
		pbr.UnitCube(pbr.Ident(), redPlastic),
		pbr.UnitCube(pbr.Ident().Trans(0, 0, -4), gold),
		pbr.UnitCube(pbr.Ident().Trans(-2, 0, 1.25), whitePlastic),
		pbr.UnitCube(pbr.Ident().Trans(1.75, 0.5, 1.75).Rot(pbr.Vector3{0, 0.4 * math.Pi, 0}).Scale(0.2, 2, 1), greenGlass),
		pbr.UnitCube(pbr.Ident().Trans(0, -1, 0).Scale(1000, 1, 1000), whiteLambert),
		pbr.UnitSphere(pbr.Ident().Trans(-2, 0.01, -2), greenGlass),
		pbr.UnitSphere(pbr.Ident().Trans(3, 0.5, 0).Scale(2, 2, 2), greenGlass),
		pbr.UnitSphere(pbr.Ident().Trans(20, 75, 20).Scale(50, 50, 50), light),
	)

	camera.MoveTo(-7, 3, 6)
	camera.LookAt(0, 0, 0)
	camera.Focus(0, 0, 0, 1.4)

	cli.Render()
}
