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

	sampler *Sampler
}

// RenderConfig configures rendering settings
type RenderConfig struct {
	Exposure float64
}

// NewRenderer creates a renderer for a sampler
func NewRenderer(s *Sampler, config ...RenderConfig) *Renderer {
	conf := RenderConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	if conf.Exposure == 0 {
		conf.Exposure = 1
	}
	return &Renderer{
		Width:        s.Width,
		Height:       s.Height,
		RenderConfig: conf,

		sampler: s,
	}
}

// Rgb averages each sample into an rgb value
func (r *Renderer) Rgb() image.Image {
	pixels := r.sampler.Samples()
	im := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	for i := 0; i < len(pixels); i += Stride {
		i2 := i / Stride * 4
		count := pixels[i+Count]
		im.Pix[i2] = r.color(pixels[i+Red] / count)
		im.Pix[i2+1] = r.color(pixels[i+Green] / count)
		im.Pix[i2+2] = r.color(pixels[i+Blue] / count)
		im.Pix[i2+3] = 255
	}
	return im
}

// Heat returns a heatmap of the sample count for each pixel
func (r *Renderer) Heat() image.Image {
	pixels := r.sampler.Samples()
	im := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	max := 0.0
	for i := Count; i < len(pixels); i += Stride {
		max = math.Max(max, pixels[i])
	}
	for i := 0; i < len(pixels); i += Stride {
		i2 := i / Stride * 4
		im.Pix[i2] = r.color(pixels[i+Count] / max * 255)
		im.Pix[i2+1] = r.color(pixels[i+Count] / max * 255)
		im.Pix[i2+2] = r.color(pixels[i+Count] / max * 255)
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
