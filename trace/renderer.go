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

func (r *Renderer) Write(file string) (n int, err error) {
	f, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	m := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	m.Set(5, 5, color.RGBA{255, 0, 0, 255})
	return 0, png.Encode(f, m)
}
