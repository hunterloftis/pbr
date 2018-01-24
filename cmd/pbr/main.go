package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/surface"
	"github.com/hunterloftis/pbr/surface/material"
)

func main() {
	if err := run(options()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func run(o *Options) error {
	print := o.Verbose || o.Info
	scene := pbr.NewScene()
	if o.Ambient != nil {
		scene.SetAmbient(*o.Ambient)
	}

	err := scene.ReadObjFile(o.Scene, o.Thin)
	if err != nil {
		return err
	}

	box, surfaces := scene.Info()
	if print {
		printSceneInfo(box, len(surfaces))
	}
	if o.Info {
		return nil
	}

	if len(o.Env) > 0 {
		err = scene.ReadHdrFile(o.Env, o.Rad)
		if err != nil {
			return err
		}
	}

	if o.Floor {
		floor := surface.UnitCube(material.Plastic(0, 0, 0, 1))
		floor.Move(box.Center.X, box.Min.Y-0.5, box.Center.Z)
		floor.Scale(100000, 1, 100000)
		scene.Add(floor)
	}

	camera := pbr.NewCamera(o.Width, o.Height)
	camera.SetLens(o.Lens)
	camera.SetStop(o.FStop)

	if o.Target == nil {
		o.Target = &box.Center
	}
	if o.Focus == nil {
		f := *o.Target
		o.Focus = &f
	}
	if o.Dist == 0 {
		o.Dist = camera.FrameDistance(box)
	}
	pos := o.Target.Plus(geom.AngleDirection(o.Lon, o.Lat).Scaled(o.Dist))
	camera.MoveTo(pos.X, pos.Y, pos.Z)
	camera.LookAt(*o.Target, *o.Focus)

	render := pbr.NewRender(scene, camera)
	render.SetBounces(o.Bounce)
	render.SetAdapt(o.Adapt)
	render.SetDirect(o.Direct)
	render.SetBranch(o.Branch)

	if print {
		printRenderInfo(o, camera)
	}
	err = iterativeRender(render, scene, o)
	if err != nil {
		return err
	}

	err = render.WritePngs(o.Out, o.Heat, o.Noise, o.Expose)
	return err
}

// TODO: move this iterative functionality into Render
func iterativeRender(r *pbr.Render, s *pbr.Scene, o *Options) error {
	perc := 0.0
	cutoff := float64(o.Width*o.Height) * o.Complete
	interrupt := make(chan os.Signal, 2)
	start := time.Now()
	end := start.Add(time.Second * time.Duration(o.Time))
	write := func(n, total int) {
		perc = float64(n) / float64(total)
		if n == total {
			r.WritePngs(o.Out, o.Heat, o.Noise, o.Expose)
		}
	}
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	if o.Profile {
		f, err := createProfile()
		if err != nil {
			return err
		}
		defer stopProfile(f)
	}
	r.Start(write)
Loop:
	for range time.NewTicker(time.Second / 10).C {
		samples := r.Count()
		select {
		case <-interrupt:
			break Loop
		default:
			completeS := float64(samples) >= cutoff
			completeT := o.Time < math.Inf(1) && time.Now().After(end)
			printProgress(r, start, s.Rays(), o.Out, perc)
			if completeS || completeT {
				break Loop
			}
		}
	}
	r.Stop()
	fmt.Println()
	return nil
}
