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
	scene, err := loadScene(o.Scene, o.Sky, o.Ground)
	if err != nil {
		return err
	}

	bounds, center, surfaces := scene.Info()
	showSceneInfo(bounds, center, len(surfaces))
	completeOptions(o, bounds, center, surfaces)
	showRenderInfo(o)

	if o.Info {
		return nil
	}

	err = loadEnvironment(scene, o.Env, o.Rad)
	if err != nil {
		return err
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
