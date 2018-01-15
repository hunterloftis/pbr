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
	if err := run(options()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func run(o *Options) error {
	scene, err := loadScene(o.Scene, o.Sky, o.Ground, o.Thin)
	if err != nil {
		return err
	}

	bounds, center, surfaces := scene.Info()
	showSceneInfo(bounds, center, len(surfaces))
	completeOptions(o, bounds, center, surfaces)

	if o.Info {
		return nil
	}

	err = loadEnvironment(scene, o.Env, o.Rad)
	if err != nil {
		return err
	}

	// TODO: implement a Plane surface type and use that instead of a scaled cube
	if o.Floor {
		p := pbr.Plastic(0.5, 0.5, 0.5, 0.1)
		t := pbr.Trans(center.X, bounds.Min.Y-0.5, center.Z)
		s := pbr.Scale(100000, 1, 100000)
		floor := pbr.UnitCube(p, t, s)
		scene.Add(floor)
	}

	// from, to, focus := cameraOptions(o, bounds, center, surfaces)
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

	showRenderInfo(o, camera)
	scene.Prepare() // TODO: make this unnecessary
	err = render(renderer, o)
	if err != nil {
		return err
	}

	err = write(renderer, o.Out, o.Heat, o.Noise, o.Expose)
	return err
}

func render(r *pbr.Renderer, o *Options) error {
	size := o.Width * o.Height
	cutoff := float64(o.Width*o.Height) * o.Complete
	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	if o.Profile {
		f, err := createProfile()
		if err != nil {
			return err
		}
		defer stopProfile(f)
	}
	savePoint := uint(size)
	start := time.Now()
	fmt.Println()
	for samples := range r.Start(time.Second / 4) {
		select {
		case <-interrupt:
			r.Stop()
			showProgress(r, start, o.Out)
		default:
			if float64(samples) >= cutoff {
				r.Stop()
			} else if samples >= savePoint {
				write(r, o.Out, o.Heat, o.Noise, o.Expose)
				savePoint *= 2
			}
			showProgress(r, start, o.Out)
		}
	}
	fmt.Println()
	return nil
}
