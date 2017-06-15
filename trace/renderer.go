package trace

import (
	"image"
	"image/png"
	"os"
)

// Renderer renders the results of a trace to a file
type Renderer struct {
	Width  int
	Height int
}

// NewRenderer sizes a Renderer to match a Camera
func NewRenderer(cam *Camera) *Renderer {
	return &Renderer{Width: cam.Width, Height: cam.Height}
}

func (r *Renderer) Write(samples []uint8, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	m := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	for i := 0; i < len(samples); i += 4 {
		m.Pix[i] = samples[i]
		m.Pix[i+1] = samples[i+1]
		m.Pix[i+2] = samples[i+2]
		m.Pix[i+3] = 255
	}
	return png.Encode(f, m)
}
