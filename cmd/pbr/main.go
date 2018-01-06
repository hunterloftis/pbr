package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr"
)

func main() {
	o := options()

	if o.Render == "" && !o.Info {
		fmt.Fprintln(os.Stderr, "error: render file required (or run with -info)")
		os.Exit(1)
	}

	scene := pbr.NewScene(o.Sky, o.Ground)

	if len(o.Env) > 0 {
		hdr, _ := os.Open(o.Env) // TODO: handle err
		defer hdr.Close()
		scene.SetPano(hdr, 150) // TODO: read radiosity info or allow it as an option
	}

	obj, err := os.Open(o.Scene)
	if err != nil {
		fmt.Println("Unable to open scene", o.Scene)
		os.Exit(1)
	}
	defer obj.Close()
	scene.ImportObj(obj)

	// whitePlastic := pbr.Plastic(0.25, 0.25, 0.25, 0.7)
	// bluePlastic := pbr.Plastic(0, 0, 0, 0.9)
	gold := pbr.Metal(1.022, 0.782, 0.344, 0.9)
	greenGlass := pbr.Plastic(0.9, 0.9, 0.9, 0.2)
	scene.Add(pbr.UnitCube(greenGlass, pbr.Trans(0, -5, 0), pbr.Scale(750, 10, 750)))
	// scene.Add(pbr.UnitSphere(greenGlass, pbr.Trans(65, 50, 85), pbr.Scale(100, 100, 100)))
	scene.Add(pbr.UnitSphere(gold, pbr.Trans(-75, 50, -125), pbr.Scale(100, 100, 100)))

	scene.Prepare()
	min, max, center, surfaces := scene.Info()
	fmt.Println("surfaces:", surfaces)
	fmt.Println("center of mass:", center)
	fmt.Println("minX:", min.X, "maxX:", max.X)
	fmt.Println("minY:", min.Y, "maxY:", max.Y)
	fmt.Println("minZ:", min.Z, "maxZ:", max.Z)
	fmt.Println()

	if o.Info {
		os.Exit(0)
	}

	size := o.Width * o.Height
	cutoff := uint(float64(size) * o.Complete)
	camera := pbr.NewCamera(o.Width, o.Height, pbr.CameraConfig{
		Lens:     o.Lens / 1000.0,
		Position: o.From,
		Target:   o.To,
		Focus:    o.Focus,
		FStop:    o.FStop,
	})
	renderer := pbr.NewRenderer(camera, scene, pbr.RenderConfig{
		Bounces: o.Bounce,
		Adapt:   o.Adapt,
	})

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	if o.Profile {
		f, _ := os.Create("cpu.pprof")
		pprof.StartCPUProfile(f)
	}

	ticker := time.NewTicker(time.Second * 60)
	start := time.Now()
	for samples := range renderer.Start(time.Second / 4) {
		select {
		case <-interrupt:
			renderer.Stop()
		case <-ticker.C:
			pbr.WritePNG(o.Render, renderer.Rgb(o.Expose))
			if len(o.Heat) > 0 {
				pbr.WritePNG(o.Heat, renderer.Heat())
			}
			if len(o.Noise) > 0 {
				pbr.WritePNG(o.Noise, renderer.Noise())
			}
		default:
			if samples >= cutoff {
				renderer.Stop()
			}
		}
		pbr.ShowProgress(renderer, start)
	}
	fmt.Println()

	if o.Profile {
		pprof.StopCPUProfile()
	}

	pbr.WritePNG(o.Render, renderer.Rgb(o.Expose))
	if len(o.Heat) > 0 {
		pbr.WritePNG(o.Heat, renderer.Heat())
	}
	if len(o.Noise) > 0 {
		pbr.WritePNG(o.Noise, renderer.Noise())
	}
}
