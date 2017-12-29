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
	scene := pbr.NewScene(*o.Sky, *o.Ground)
	camera := pbr.NewCamera(o.Width, o.Height, pbr.CameraConfig{
		Lens:     o.Lens / 1000.0,
		Position: o.From,
		Target:   o.To,
		Focus:    o.Focus,
		FStop:    o.FStop,
	})
	adapt := 5
	if o.Uniform {
		adapt = 0
	}
	// TODO: should Renderer and Sampler be separate?
	sampler := pbr.NewSampler(camera, scene, pbr.SamplerConfig{
		Bounces: o.Bounce,
		Adapt:   adapt, // TODO: make this boolean
	})
	renderer := pbr.NewRenderer(sampler, pbr.RenderConfig{
		Exposure: o.Expose,
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

	render(sampler, renderer, o.Exit)
	pbr.WritePNG(o.Render, renderer.Rgb()) // TODO: should o.Expose be passed in here instead of as a global option?
	if len(o.Heat) > 0 {
		pbr.WritePNG(o.Heat, renderer.Heat())
	}
}

func render(sampler *pbr.Sampler, renderer *pbr.Renderer, quality float64) {
	start := time.Now()
	running := true
	interrupt := make(chan os.Signal, 2)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // TODO: abstract this?
	go func() {
		<-interrupt
		running = false
	}()

	for running && sampler.PerPixel() < quality {
		pbr.ShowProgress(sampler, start, running)
		sampler.Sample()
	}
	pbr.ShowProgress(sampler, start, running)
}
