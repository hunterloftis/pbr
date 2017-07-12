package pbr

import (
	"image"
	"math"
)

// Renderer renders the results of a trace to a file
type Renderer struct {
	Width  int
	Height int
	RenderConfig
	pixels []float64
}

// RenderConfig configures rendering settings
type RenderConfig struct {
	Exposure float64
}

// CamRenderer sizes a Renderer to match a Camera
func CamRenderer(c *Camera, config ...RenderConfig) *Renderer {
	conf := config[0]
	if conf.Exposure == 0 {
		conf.Exposure = 1
	}
	return &Renderer{
		Width:        c.Width,
		Height:       c.Height,
		RenderConfig: conf,
		pixels:       make([]float64, 0),
	}
}

// Merge merges pixel arrays
func (r *Renderer) Merge(pixels []float64) {
	if len(r.pixels) < len(pixels) {
		r.pixels = make([]float64, len(pixels))
	}
	for i, val := range pixels {
		r.pixels[i] += val
	}
}

// Rgb averages each sample into an rgb value
func (r *Renderer) Rgb() image.Image {
	im := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	for i := 0; i < len(r.pixels); i += Elements {
		i2 := i / Elements * 4
		count := r.pixels[i+Count]
		im.Pix[i2] = r.color(r.pixels[i+Red] / count)
		im.Pix[i2+1] = r.color(r.pixels[i+Green] / count)
		im.Pix[i2+2] = r.color(r.pixels[i+Blue] / count)
		im.Pix[i2+3] = 255
	}
	return im
}

// Heat returns a heatmap of the sample count for each pixel
func (r *Renderer) Heat() image.Image {
	im := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	max := 0.0
	for i := Count; i < len(r.pixels); i += Elements {
		max = math.Max(max, r.pixels[i])
	}
	for i := 0; i < len(r.pixels); i += Elements {
		i2 := i / Elements * 4
		im.Pix[i2] = r.color(r.pixels[i+Count] / max * 255)
		im.Pix[i2+1] = r.color(r.pixels[i+Count] / max * 255)
		im.Pix[i2+2] = r.color(r.pixels[i+Count] / max * 255)
		im.Pix[i2+3] = 255
	}
	return im
}

func (r *Renderer) color(n float64) uint8 {
	return uint8(gamma(math.Min(n*r.Exposure, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}
