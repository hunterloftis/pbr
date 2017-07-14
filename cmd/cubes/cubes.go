package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "", "Heatmap png filename")
	profile := flag.String("profile", "", "Record performance into profile.pprof")
	workers := runtime.NumCPU()
	flag.Parse()

	// https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs
	switch *profile {
	case "block":
		f, _ := os.Create("profile.pprof")
		runtime.SetBlockProfileRate(1)
		defer pprof.Lookup("block").WriteTo(f, 10)
	case "cpu":
		f, _ := os.Create("profile.pprof")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	scene := pbr.EmptyScene()
	camera := pbr.NewCamera(1280, 720, pbr.CameraConfig{
		Position: pbr.Vector3{-0.6, 0.12, 0.8},
		Target:   &pbr.Vector3{0, 0, 0},
		Focus:    &pbr.Vector3{0, -0.025, 0.2},
		FStop:    4,
	})
	renderer := pbr.CamRenderer(camera)

	light := pbr.Light(1500, 1500, 1500)
	redPlastic := pbr.Plastic(1, 0, 0, 1)
	whitePlastic := pbr.Plastic(1, 1, 1, 0.8)
	// grayLambert := pbr.Lambert(0.2, 0.2, 0.2, 0.1)
	bluePlastic := pbr.Plastic(0, 0, 1, 1)
	greenPlastic := pbr.Plastic(0, 0.9, 0, 1)
	// silver := pbr.Metal(0.95, 0.93, 0.88, 1)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	// glass := pbr.Glass(1, 1, 1, 1)
	greenGlass := pbr.Glass(0.2, 1, 0.1, 0.95)

	scene.SetSky(pbr.Vector3{40, 50, 60}, pbr.Vector3{})

	scene.Add(
		pbr.UnitCube(pbr.Ident().Rot(pbr.Vector3{0, -0.25 * math.Pi, 0}).Scale(0.1, 0.1, 0.1), redPlastic),
		pbr.UnitCube(pbr.Ident().Trans(0, 0, -0.4).Rot(pbr.Vector3{0, 0.1 * math.Pi, 0}).Scale(0.1, 0.1, 0.1), gold),
		pbr.UnitCube(pbr.Ident().Trans(-0.3, 0, 0.3).Rot(pbr.Vector3{0, -0.1 * math.Pi, 0}).Scale(0.1, 0.1, 0.1), greenGlass),
		pbr.UnitCube(pbr.Ident().Trans(0.175, 0.05, 0.18).Rot(pbr.Vector3{0, 0.55 * math.Pi, 0}).Scale(0.02, 0.2, 0.2), greenGlass),
		pbr.UnitCube(pbr.Ident().Trans(0, -0.55, 0).Scale(1000, 1, 1000), whitePlastic).SetGrid(bluePlastic, 1.0/20.0),
		pbr.UnitSphere(pbr.Ident().Trans(-0.2, 0.001, -0.2).Scale(0.1, 0.1, 0.1), greenGlass),
		pbr.UnitSphere(pbr.Ident().Trans(0.3, 0.05, 0).Scale(0.2, 0.2, 0.2), bluePlastic),
		pbr.UnitSphere(pbr.Ident().Trans(7, 30, 6).Scale(30, 30, 30), light),
		pbr.UnitSphere(pbr.Ident().Trans(0, -0.025, 0.2).Scale(0.1, 0.05, 0.1), greenPlastic),
		pbr.UnitSphere(pbr.Ident().Trans(0.45, 0.05, -0.4).Scale(0.2, 0.2, 0.2), gold),
	)

	start := time.Now().UnixNano()
	m := pbr.NewMonitor()
	m.Interrupt()

	for i := 0; i < workers; i++ {
		m.AddSampler(pbr.NewSampler(camera, scene, pbr.SamplerConfig{
			Bounces: 10,
			Adapt:   5,
		}))
	}

	go func() {
		for m.Active() > 0 {
			samples := <-m.Progress
			ns := int(time.Now().UnixNano() - start)
			showProgress(samples, camera.Pixels(), ns, m.Stopped())
		}
	}()

	for i := 0; i < workers; i++ {
		renderer.Merge(<-m.Results)
	}

	writePNG(*out, renderer.Rgb())
	showFile(*out)
	if len(*heat) > 0 {
		writePNG(*heat, renderer.Heat())
		showFile(*heat)
	}
}

func writePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, i)
}

func showProgress(samples, pixels, ns int, stopped bool) {
	note := ""
	if stopped {
		note = " (wrapping up...)"
	}
	pp := samples / pixels
	pms := samples / (ns / 1e6)
	fmt.Printf("\r%v samples/pixel, %v samples/ms%v", pp, pms, note) // https://stackoverflow.com/a/15442704/1911432
}

func showFile(file string) {
	fmt.Printf("\n-> %v\n", file)
}
