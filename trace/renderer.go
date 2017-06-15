package trace

import (
	"image"
	"image/png"
	"math"
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

func (r *Renderer) Write(samples []float64, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	m := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	for i := 0; i < len(samples); i += 3 {
		i2 := i / 3 * 4
		m.Pix[i2] = uint8(math.Min(samples[i], 255))
		m.Pix[i2+1] = uint8(math.Min(samples[i+1], 255))
		m.Pix[i2+2] = uint8(math.Min(samples[i+2], 255))
		m.Pix[i2+3] = 255
	}
	return png.Encode(f, m)
}
