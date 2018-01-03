package pbr

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"
)

// WritePNG saves an image to a png file and prints the filename.
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

// ShowProgress shows the current sampling progress (billion samples, samples/pixel, samples/ms, shutting down?).
// https://stackoverflow.com/a/15442704/1911432
func ShowProgress(r *Renderer, start time.Time) {
	var pp, pms uint
	var note string
	var mil float64
	if !r.Active() {
		note = " (wrapping up...)"
	}
	samples := r.Count()
	pixels := r.Size()
	passed := time.Now().Sub(start)
	secs := passed.Seconds()
	ms := uint(passed.Nanoseconds() / 1e6)
	if pixels > 0 && ms > 0 {
		mil = float64(samples) / 1e6
		pp = samples / pixels
		pms = samples / ms
	}
	fmt.Printf("\r%ds - %.3f million samples - %v samples/pixel - %v samples/ms%v", int(secs), mil, pp, pms, note)
}
