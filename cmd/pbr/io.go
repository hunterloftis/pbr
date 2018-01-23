package main

import (
	"fmt"
	"math"
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
// TODO: clean up, this is messy
// TODO: hide this in --silent mode
func printProgress(r *pbr.Render, start time.Time, rays uint64, filename string, perc float64) {
	var pp, pms uint
	var mil, mRays, rayPerSample float64
	samples := r.Count()
	pixels := r.Size()
	passed := time.Since(start)
	secs := passed.Seconds()
	ms := uint(passed.Nanoseconds() / 1e6)
	if pixels > 0 && ms > 0 {
		mil = float64(samples) / 1e6
		pp = samples / pixels
		pms = samples / ms
		mRays = float64(rays) / 1e6
		rayPerSample = float64(rays) / float64(samples)
	}
	progress := ""
	count := int(math.Min(math.Ceil(perc*8), 8))
	for i := 0; i < count; i++ {
		progress += "="
	}
	progress += ">"
	for i := count; i < 8; i++ {
		progress += " "
	}
	const str = "\r%.1fs, samples (%.1f mil, %v/ms), rays (%.1f mil, %.1f/sample) [ %v s/px %v %v ]  "
	fmt.Printf(str, secs, mil, pms, mRays, rayPerSample, pp, progress, filename)
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
