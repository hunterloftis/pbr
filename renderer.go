package pbr

import (
	"image"
	"math"
	"runtime"
	"time"
)

// Renderer renders the samples in a Sampler to an Image.
type Renderer struct {
	Width  int
	Height int
	RenderConfig

	camera *Camera
	scene  *Scene

	// state
	active bool
	pixels []float64 // stored in a flat array, chunked by Stride
	count  uint
}

// RenderConfig configures rendering settings.
type RenderConfig struct {
	Bounces int
	Adapt   bool
}

// NewRenderer creates a renderer referencing a Sampler.
func NewRenderer(c *Camera, s *Scene, config ...RenderConfig) *Renderer {
	conf := RenderConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	return &Renderer{
		Width:        c.Width,
		Height:       c.Height,
		RenderConfig: conf,

		camera: c,
		scene:  s,
	}
}

func (r *Renderer) Start(tick time.Duration) <-chan uint {
	r.active = true
	n := runtime.NumCPU()
	samplers := make([]*Sampler, n)
	ticker := time.NewTicker(tick)
	ch := make(chan uint)
	results := make(chan result, n)
	stop := make(chan struct{})
	for i := 0; i < n; i++ {
		samplers[i] = NewSampler(r.camera, r.scene, SamplerConfig{})
		samplers[i].Sample(results, stop)
	}
	go func() {
		for {
			select {
			case res := <-results:
				r.integrate(res)
				r.count++
			case <-ticker.C:
				ch <- r.count
			default:
				if !r.active {
					close(stop)
					close(results)
					close(ch)
					return
				}
			}
		}
	}()
	return ch
}

func (r *Renderer) Stop() {
	r.active = false
}

func (r *Renderer) Active() bool {
	return r.active
}

func (r *Renderer) Count() uint {
	return r.count
}

func (r *Renderer) Size() uint {
	return uint(r.camera.Width * r.camera.Height)
}

func (r *Renderer) integrate(res result) {

}

// Rgb averages each sample into an rgb value.
func (r *Renderer) Rgb(expose float64) image.Image {
	im := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	length := uint(len(r.pixels))
	for i := uint(0); i < length; i += Stride {
		i2 := i / Stride * 4
		count := r.pixels[i+Count]
		im.Pix[i2] = r.color(r.pixels[i+Red] / count * expose)
		im.Pix[i2+1] = r.color(r.pixels[i+Green] / count * expose)
		im.Pix[i2+2] = r.color(r.pixels[i+Blue] / count * expose)
		im.Pix[i2+3] = 255
	}
	return im
}

// Heat returns a heatmap of the sample count for each pixel.
func (r *Renderer) Heat() image.Image {
	im := image.NewRGBA(image.Rect(0, 0, r.Width, r.Height))
	max := 0.0
	length := uint(len(r.pixels))
	for i := Count; i < length; i += Stride {
		max = math.Max(max, r.pixels[i])
	}
	for i := uint(0); i < length; i += Stride {
		i2 := i / Stride * 4
		im.Pix[i2] = r.color(r.pixels[i+Count] / max * 255)
		im.Pix[i2+1] = r.color(r.pixels[i+Count] / max * 255)
		im.Pix[i2+2] = r.color(r.pixels[i+Count] / max * 255)
		im.Pix[i2+3] = 255
	}
	return im
}

func (r *Renderer) color(n float64) uint8 {
	return uint8(gamma(math.Min(n, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}
