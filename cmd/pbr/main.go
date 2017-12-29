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

	// For debugging until we're actually parsing scene files
	// scene.Add(pbr.UnitCube(pbr.Plastic(1, 0, 0, 1), pbr.Rot(pbr.Vector3{0, 1, 0}), pbr.Scale(0.5, 0.5, 0.5)))
	mesh := pbr.Mesh{
		Tris: []pbr.Triangle{
			pbr.NewTriangle(pbr.Vector3{0, 0, 0}, pbr.Vector3{0, 1, 0}, pbr.Vector3{1, 0, 0}),
		},
		Pos: pbr.Identity(),
		Mat: pbr.Plastic(1, 0, 0, 1),
	}
	scene.Add(&mesh)

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
