package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr"
)

// pbr [options] scene.obj render.png

func main() {
	o := options()
	fmt.Println(o)

	scene := pbr.EmptyScene()
	camera := pbr.NewCamera(o.Width, o.Height, pbr.CameraConfig{
		Lens:     o.Lens / 1000.0,
		Position: &o.From, // TODO: why by address?
		Target:   &o.To,   // TODO: why by address?
		Focus:    o.Focus,
		FStop:    o.FStop,
	})
	adapt := 5
	if o.Uniform {
		adapt = 0
	}
	sampler := pbr.NewSampler(camera, scene, pbr.SamplerConfig{
		Bounces: o.Bounce,
		Adapt:   adapt, // TODO: make this boolean
	})
	renderer := pbr.NewRenderer(sampler, pbr.RenderConfig{
		Exposure: o.Expose,
	})
	scene.SetSky(o.Sky, pbr.Vector3{})
	if len(o.Env) > 0 {
		hdr, _ := os.Open(o.Env)
		defer hdr.Close()
		scene.SetPano(hdr, 100) // TODO: read radiosity info or allow it as an option
	}

	// For debugging until we're actually parsing collada files
	scene.Add(pbr.UnitCube(pbr.Plastic(1, 0, 0, 1), pbr.Rot(pbr.Vector3{0, 1, 0}), pbr.Scale(0.5, 0.5, 0.5)))

	start := time.Now()
	running := true
	interrupt := make(chan os.Signal, 2)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // TODO: abstract this?
	go func() {
		<-interrupt
		running = false
		pbr.ShowProgress(sampler, start, running)
	}()

	pbr.ShowProgress(sampler, start, running)
	for running && sampler.PerPixel() < o.Exit {
		sampler.Sample()
		pbr.ShowProgress(sampler, start, running)
	}

	pbr.WritePNG(o.Render, renderer.Rgb())
	if len(o.Heat) > 0 {
		pbr.WritePNG(o.Heat, renderer.Heat())
	}
}
