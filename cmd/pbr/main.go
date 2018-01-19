package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/surface"
)

// TODO: make structs with constructors private?
// http://www.golangpatterns.info/object-oriented/constructors
func main() {
	if err := run(options()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func run(o *Options) error {
	scene := pbr.NewScene()
	if o.Ambient != nil {
		scene.SetAmbient(*o.Ambient)
	}

	err := scene.ReadObjFile(o.Scene, o.Thin)
	if err != nil {
		return err
	}

	box, surfaces := scene.Info()
	printSceneInfo(box, len(surfaces))
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
		floor := surface.UnitCube()
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
		o.Focus = o.Target
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

	printRenderInfo(o, camera)
	err = iterativeRender(render, o)
	if err != nil {
		return err
	}

	err = render.WritePngs(o.Out, o.Heat, o.Noise, o.Expose)
	return err
}

// TODO: this is a bit messy. Better way?
func iterativeRender(r *pbr.Render, o *Options) error {
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
	fmt.Println()
	savePoint := uint(size)
	start := time.Now()
	ticker := time.NewTicker(time.Second / 4)
	r.Start()
Loop:
	for range ticker.C {
		samples := r.Count()
		select {
		case <-interrupt:
			ticker.Stop()
			break Loop
		default:
			if float64(samples) >= cutoff {
				ticker.Stop()
				break Loop
			} else if samples >= savePoint {
				r.WritePngs(o.Out, o.Heat, o.Noise, o.Expose)
				savePoint *= 2
			}
			printProgress(r, start, o.Out)
		}
	}
	r.Stop()
	printProgress(r, start, o.Out)
	fmt.Println()
	return nil
}
