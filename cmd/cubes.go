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
	whitePlastic := pbr.Plastic(1, 1, 1, 0)
	bluePlastic := pbr.Plastic(0, 0, 1, 1)
	silver := pbr.Metal(0.972, 0.960, 0.915, 1)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	glass := pbr.Glass(0, 0, 0, 0, 1)

	scene.SetEnv("images/glacier.hdr", 100)

	scene.Add(
		pbr.UnitCube(pbr.Ident(), silver),
		pbr.UnitCube(pbr.Ident().Trans(0, 0, -4), gold),
		pbr.UnitCube(pbr.Ident().Trans(-2, -0.245, 1.25).Scale(1, 0.5, 1), glass),
		pbr.UnitCube(pbr.Ident().Trans(1.75, 0, 2).Rot(pbr.Vector3{0, 0.25 * math.Pi, 0}).Scale(0.5, 2, 1), redPlastic),
		pbr.UnitCube(pbr.Ident().Trans(0, -1, 0).Scale(1000, 1, 1000), whitePlastic),
		pbr.UnitSphere(pbr.Ident().Trans(-1.5, 0, 0), bluePlastic),
		pbr.UnitSphere(pbr.Ident().Trans(-4.5, 0, 1.5), gold),
		pbr.UnitSphere(pbr.Ident().Trans(2, 0, 0), silver),
		pbr.UnitSphere(pbr.Ident().Trans(50, 30, 0).Scale(40, 40, 40), light),
	)

	camera.MoveTo(-6, 1.5, 5)
	camera.LookAt(0, 0, 0)
	camera.Focus(0, 0, 0, 1.4)

	cli.Start()
}
