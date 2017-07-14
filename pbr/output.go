package pbr

import (
	"fmt"
	"image"
	"image/png"
	"os"
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
func ShowProgress(samples, pixels, ns int, stopped bool) {
	var pp, pms int
	var note string
	var bil float64
	var ms = ns / 1e6
	if stopped {
		note = " (wrapping up...)"
	}
	if pixels > 0 && ms > 0 {
		bil = float64(samples) / 1e9
		pp = samples / pixels
		pms = samples / ms
	}
	fmt.Printf("\rsamples: %.3f billion - %v/pixel - %v/ms%v", bil, pp, pms, note) // https://stackoverflow.com/a/15442704/1911432
}
