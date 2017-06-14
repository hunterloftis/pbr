package trace

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// Renderer renders the results of a trace to a file
type Renderer struct {
	Width  int
	Height int
}

func (r *Renderer) Write(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	m := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	m.Set(5, 5, color.RGBA{255, 0, 0, 255})
	return png.Encode(f, m)
}
