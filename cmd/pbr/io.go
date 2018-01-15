package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/hunterloftis/pbr"
	"github.com/hunterloftis/pbr/obj"
)

func loadScene(filename string, sky, ground *pbr.Energy, thin bool) (*pbr.Scene, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open scene %v, %v", filename, err)
	}
	defer f.Close()
	scene := pbr.NewScene(sky, ground)
	scanner := obj.NewScanner(f)
	for scanner.Scan() {
		if m := scanner.Material(); len(m) > 0 {
			matfile := filepath.Join(filepath.Dir(filename), m)
			f, err := os.Open(matfile)
			if err != nil {
				continue
			}
			defer f.Close()
			err = scanner.ReadMaterials(f, thin)
			if err != nil {
				return nil, err
			}
			continue
		}
		scene.Add(scanner.Surface())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return scene, nil
}

// TODO: only show stuff in --verbose mode
func showSceneInfo(bounds *pbr.Box, center pbr.Vector3, surfaces int) {
	fmt.Printf("surfaces: %v\n", surfaces)
	fmt.Printf("center of mass: (%.2f, %.2f, %.2f)\n", center.X, center.Y, center.Z)
	fmt.Printf("X range: [%.2f : %.2f]\n", bounds.Min.X, bounds.Max.X)
	fmt.Printf("Y range: [%.2f : %.2f]\n", bounds.Min.Y, bounds.Max.Y)
	fmt.Printf("Z range: [%.2f : %.2f]\n", bounds.Min.Z, bounds.Max.Z)
}

func showRenderInfo(o *Options, c *pbr.Camera) {
	fmt.Printf("Camera lens: %vmm, dist: %.2f, polar: %.2f\n", o.Lens, o.Dist, o.Polar)
	fmt.Printf("Camera position: (%.2f, %.2f, %.2f)\n", c.Position.X, c.Position.Y, c.Position.Z)
}

func loadEnvironment(s *pbr.Scene, filename string, exposure float64) error {
	if len(filename) == 0 {
		return nil
	}
	hdr, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer hdr.Close()
	s.SetPano(hdr, exposure)
	return nil
}

func write(renderer *pbr.Renderer, out, heat, noise string, expose float64) error {
	if err := writePNG(out, renderer.Rgb(expose)); err != nil {
		return err
	}
	if len(heat) > 0 {
		if err := writePNG(heat, renderer.Heat()); err != nil {
			return err
		}
	}
	if len(noise) > 0 {
		if err := writePNG(noise, renderer.Noise()); err != nil {
			return err
		}
	}
	return nil
}

func writePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, i)
	return err
}

// https://stackoverflow.com/a/15442704/1911432
// TODO: hide this in --silent mode
func showProgress(r *pbr.Renderer, start time.Time, filename string) {
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
