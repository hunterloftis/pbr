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
	scene := pbr.NewScene(o.Sky, o.Ground)

	obj, err := os.Open(o.Scene)
	if err != nil {
		fmt.Println("Unable to open scene", o.Scene, "error:", err)
		os.Exit(1)
	}
	defer obj.Close()
	scene.ImportObj(obj)

	scene.Prepare()
	min, max, center, surfaces := scene.Info()
	fmt.Println()
	fmt.Printf("surfaces: %v\n", surfaces)
	fmt.Printf("center of mass: (%.2f, %.2f, %.2f)\n", center.X, center.Y, center.Z)
	fmt.Printf("X range: [%.2f : %.2f]\n", min.X, max.X)
	fmt.Printf("Y range: [%.2f : %.2f]\n", min.Y, max.Y)
	fmt.Printf("Z range: [%.2f : %.2f]\n", min.Z, max.Z)
	fmt.Println()

	if o.Info {
		os.Exit(0)
	}

	if o.From == nil {
		twoThirds := pbr.Vector3{max.X * 9, max.Y, max.Z * 6}
		o.From = &twoThirds
	}
	if o.To == nil {
		o.To = &center
	}
	if o.Focus == nil {
		o.Focus = o.To
	}

	if len(o.Env) > 0 {
		hdr, _ := os.Open(o.Env) // TODO: handle err
		defer hdr.Close()
		scene.SetPano(hdr, o.Rad)
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
			pbr.WritePNG(o.Out, renderer.Rgb(o.Expose))
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

	fmt.Println("->", o.Out)
	pbr.WritePNG(o.Out, renderer.Rgb(o.Expose))
	if len(o.Heat) > 0 {
		pbr.WritePNG(o.Heat, renderer.Heat())
	}
	if len(o.Noise) > 0 {
		pbr.WritePNG(o.Noise, renderer.Noise())
	}
}
