package pbr

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"
)

// WritePNG saves an image to a png file
func WritePNG(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, i)
	if err == nil {
		fmt.Printf("\n-> %v\n", file)
	}
	return err
}

// ShowProgress shows the current sampling progress
// https://stackoverflow.com/a/15442704/1911432
func ShowProgress(s *Sampler, start time.Time, running bool) {
	var pp, pms int
	var note string
	var bil float64
	if !running {
		note = " (wrapping up...)"
	}
	samples := s.Count()
	pixels := s.Pixels()
	ms := int(time.Now().Sub(start).Nanoseconds() / 1e6)
	if s.Pixels() > 0 && ms > 0 {
		bil = float64(samples) / 1e9
		pp = samples / pixels
		pms = samples / ms
	}
	fmt.Printf("\rsamples: %.3f billion - %v/pixel - %v/ms%v", bil, pp, pms, note)
}
