package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/surface"
)

// TODO: only show stuff in --verbose mode
func printSceneInfo(box *surface.Box, surfaces int) {
	fmt.Printf("surfaces: %v\n", surfaces)
	fmt.Printf("center of mass: (%.2f, %.2f, %.2f)\n", box.Center.X, box.Center.Y, box.Center.Z)
	fmt.Printf("X range: [%.2f : %.2f]\n", box.Min.X, box.Max.X)
	fmt.Printf("Y range: [%.2f : %.2f]\n", box.Min.Y, box.Max.Y)
	fmt.Printf("Z range: [%.2f : %.2f]\n", box.Min.Z, box.Max.Z)
}

func printRenderInfo(o *Options, c *pbr.Camera) {
	pos, target, focus := c.Orientation()
	fmt.Printf("Camera lens: %vmm, dist: %.2f, lat: %.2f, lon: %.2f\n", o.Lens, o.Dist, o.Lat, o.Lon)
	fmt.Printf("Camera position: %v, target: %v, focus: %v\n", pos, target, focus)
}

// https://stackoverflow.com/a/15442704/1911432
// TODO: hide this in --silent mode
func printProgress(r *pbr.Render, start time.Time, filename string) {
	var pp, pms uint
	var note string
	var mil, raysMs float64
	if !r.Active() {
		note = "(wrapping up...)"
	}
	samples := r.Count()
	pixels := r.Size()
	passed := time.Since(start)
	secs := passed.Seconds()
	ms := uint(passed.Nanoseconds() / 1e6)
	if pixels > 0 && ms > 0 {
		mil = float64(samples) / 1e6
		pp = samples / pixels
		pms = samples / ms
		raysMs = float64(r.Rays()) / float64(ms)
	}
	fmt.Printf("\r%ds - %.3f million samples - %v samples/pixel - %v samples/ms - %.1f rays/ms => %v %v", int(secs), mil, pp, pms, raysMs, filename, note)
}

func createProfile() (*os.File, error) {
	f, err := os.Create("cpu.pprof")
	if err != nil {
		return nil, err
	}
	pprof.StartCPUProfile(f)
	return f, nil
}

func stopProfile(f *os.File) {
	pprof.StopCPUProfile()
	f.Close()
}
