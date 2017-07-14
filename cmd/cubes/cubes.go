package main

import (
	"flag"
	"math"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr/pbr"
)

func main() {
	out := flag.String("out", "render.png", "Output png filename")
	heat := flag.String("heat", "heat.png", "Heatmap png filename")
	profile := flag.String("profile", "", "Record performance into profile.pprof")
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
		Position: &pbr.Vector3{-0.6, 0.12, 0.8},
		Target:   &pbr.Vector3{0, 0, 0},
		Focus:    &pbr.Vector3{0, -0.025, 0.2},
		FStop:    4,
	})
	sampler := pbr.NewSampler(camera, scene)
	renderer := pbr.NewRenderer(sampler)

	light := pbr.Light(1500, 1500, 1500)
	redPlastic := pbr.Plastic(1, 0, 0, 1)
	whitePlastic := pbr.Plastic(1, 1, 1, 0.8)
	bluePlastic := pbr.Plastic(0, 0, 1, 1)
	greenPlastic := pbr.Plastic(0, 0.9, 0, 1)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	greenGlass := pbr.Glass(0.2, 1, 0.1, 0.95)

	scene.SetSky(pbr.Vector3{40, 50, 60}, pbr.Vector3{})
	scene.Add(
		pbr.UnitCube(redPlastic, pbr.Rot(pbr.Vector3{0, -0.25 * math.Pi, 0}), pbr.Scale(0.1, 0.1, 0.1)),
		pbr.UnitCube(gold, pbr.Trans(0, 0, -0.4), pbr.Rot(pbr.Vector3{0, 0.1 * math.Pi, 0}), pbr.Scale(0.1, 0.1, 0.1)),
		pbr.UnitCube(greenGlass, pbr.Trans(-0.3, 0, 0.3), pbr.Rot(pbr.Vector3{0, -0.1 * math.Pi, 0}), pbr.Scale(0.1, 0.1, 0.1)),
		pbr.UnitCube(greenGlass, pbr.Trans(0.175, 0.05, 0.18), pbr.Rot(pbr.Vector3{0, 0.55 * math.Pi, 0}), pbr.Scale(0.02, 0.2, 0.2)),
		pbr.UnitCube(whitePlastic, pbr.Trans(0, -0.55, 0), pbr.Scale(1000, 1, 1000)), // .SetGrid(bluePlastic, 1.0/20.0)
		pbr.UnitSphere(greenGlass, pbr.Trans(-0.2, 0.001, -0.2), pbr.Scale(0.1, 0.1, 0.1)),
		pbr.UnitSphere(bluePlastic, pbr.Trans(0.3, 0.05, 0), pbr.Scale(0.2, 0.2, 0.2)),
		pbr.UnitSphere(light, pbr.Trans(7, 30, 6), pbr.Scale(30, 30, 30)),
		pbr.UnitSphere(greenPlastic, pbr.Trans(0, -0.025, 0.2), pbr.Scale(0.1, 0.05, 0.1)),
		pbr.UnitSphere(gold, pbr.Trans(0.45, 0.05, -0.4), pbr.Scale(0.2, 0.2, 0.2)),
	)

	start := time.Now()
	running := true
	interrupt := make(chan os.Signal, 2)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interrupt
		running = false
		pbr.ShowProgress(sampler, start, running)
	}()

	pbr.ShowProgress(sampler, start, running)
	for running {
		sampler.Sample()
		pbr.ShowProgress(sampler, start, running)
	}

	pbr.WritePNG(*out, renderer.Rgb())
	pbr.WritePNG(*heat, renderer.Heat())
}
