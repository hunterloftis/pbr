package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr"
)

func main() {
	o := options()
	size := o.Width * o.Height
	cutoff := uint(float64(size) * o.Complete)
	scene := pbr.NewScene(*o.Sky, *o.Ground)
	camera := pbr.NewCamera(o.Width, o.Height, pbr.CameraConfig{
		Lens:     o.Lens / 1000.0,
		Position: o.From,
		Target:   o.To,
		Focus:    o.Focus,
		FStop:    o.FStop,
	})
	renderer := pbr.NewRenderer(camera, scene, pbr.RenderConfig{
		Bounces: o.Bounce,
		Uniform: o.Uniform,
	})

	if len(o.Env) > 0 {
		hdr, _ := os.Open(o.Env) // TODO: handle err
		defer hdr.Close()
		scene.SetPano(hdr, 100) // TODO: read radiosity info or allow it as an option
	}

	obj, err := os.Open(o.Scene)
	if err != nil {
		fmt.Println("Unable to open scene", o.Scene)
		os.Exit(1)
	}
	defer obj.Close()
	scene.ImportObj(obj)

	whitePlastic := pbr.Plastic(1, 1, 1, 0.8)
	bluePlastic := pbr.Plastic(0, 0, 1, 1)
	scene.Add(pbr.UnitCube(whitePlastic, pbr.Trans(0, 11, -600), pbr.Scale(10000, 1, 10000)).SetGrid(bluePlastic, 8.0))

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	start := time.Now()
	for samples := range renderer.Start(time.Second / 4) {
		select {
		case <-interrupt:
			renderer.Stop()
		default:
			if samples >= cutoff {
				renderer.Stop()
			}
		}
		pbr.ShowProgress(renderer, start)
	}

	pbr.WritePNG(o.Render, renderer.Rgb(o.Expose))
	if len(o.Heat) > 0 {
		pbr.WritePNG(o.Heat, renderer.Heat())
	}
}
