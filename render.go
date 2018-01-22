package pbr

import (
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/hunterloftis/pbr/rgb"
)

// Render samples Rays from a Camera into a Scene and records the results into a Framebuffer.
// A Render manages its own concurrency and can be started and stopped at any point.
type Render struct {
	sampler

	adapt float64

	// state
	rnd          *rand.Rand
	active       bool
	cursor       uint
	meanVariance float64
	buffer       *rgb.Framebuffer
	samples      uint
}

// NewRender constructs a new Render from a Camera into a Scene.
func NewRender(s *Scene, c *Camera) *Render {
	return &Render{
		sampler: sampler{
			bounces: 8,
			branch:  32,
			direct:  8,
			camera:  c,
			scene:   s,
		},
		adapt:        8,
		meanVariance: math.MaxFloat64,
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		buffer:       rgb.NewBuffer(uint(c.Width()), uint(c.Height())),
	}
}

// SetBounces sets the maximum number of times a light Ray can bounce within a Scene.
// Set to zero for direct light only (no bounces). 9 is very high quality.
func (r *Render) SetBounces(b int) {
	r.bounces = b
}

// SetAdapt sets the maximum number of extra samples that can be taken per frame to resolve noisy pixels.
func (r *Render) SetAdapt(a float64) {
	r.adapt = a
}

// SetDirect sets the maximum number of direct lights that can be sampled per bounce.
// Higher numbers increase the per-pixel sample accuracy in scenes with specific light sources.
func (r *Render) SetDirect(d int) {
	r.direct = d
}

// SetBranch sets the maximum number of light Ray branches that will be created from primary Rays that hit objects.
// Higher values increase per-pixel sample accuracy at the cost of per-frame render time.
func (r *Render) SetBranch(b int) {
	r.branch = b
}

// Active returns true if the Render has been started or false otherwise.
func (r *Render) Active() bool {
	return r.active
}

// Count returns the total number of samples that have been taken.
func (r *Render) Count() uint {
	return r.samples
}

// Size returns the total number of pixels in this Render.
func (r *Render) Size() uint {
	return uint(r.camera.Width() * r.camera.Height())
}

// Buffer returns a reference to the Framebuffer storing all of the light energy (rgb) data for this Render.
func (r *Render) Buffer() *rgb.Framebuffer {
	return r.buffer
}

// Image processes the current Render state into a 2D RGB Image.
func (r *Render) Image(expose float64) image.Image {
	return r.buffer.Image(expose)
}

// Start begins rendering the Scene.
func (r *Render) Start() {
	r.active = true
	r.scene.prepare()
	n := runtime.NumCPU()
	pixel := make(chan *[sampleSize]uint, n)
	result := make(chan *[sampleSize]sample, n)
	for i := 0; i < n; i++ {
		r.sampler.start(pixel, result)
		pixel <- r.next()
	}
	go func() {
		for {
			for _, res := range <-result {
				r.samples++
				r.buffer.Integrate(res.index, res.energy)
			}
			pixel <- r.next()
			if !r.active {
				close(pixel)
				return
			}
		}
	}()
}

// Stop stops rendering.
func (r *Render) Stop() {
	r.active = false
}

// WritePngs writes up to three images to the filesystem:
// A standard RGB image from the current Framebuffer,
// a heatmap of the sample frequencies,
// and a noisemap of the variance of each pixel.
// Expose multiplies the brightness of the RGB image.
// Empty filenames ("") are skipped.
// Returns any error encountered while writing files.
func (r *Render) WritePngs(out, heat, noise string, expose float64) error {
	if len(out) > 0 {
		if err := writePng(out, r.buffer.Image(expose)); err != nil {
			return err
		}
	}
	if len(heat) > 0 {
		if err := writePng(heat, r.buffer.Heatmap()); err != nil {
			return err
		}
	}
	if len(noise) > 0 {
		if err := writePng(noise, r.buffer.Noisemap()); err != nil {
			return err
		}
	}
	return nil
}

func (r *Render) deficit(i uint) (count int) {
	brightness := r.buffer.Average(i).Average()
	if brightness < 255 && brightness > 0 {
		midtones := (((255 - brightness) / 255) + 3) / 4
		noise := r.buffer.Noise(i)
		varMean, countMean := r.buffer.Variance()
		ratio := (noise + 1) / (varMean + 1)
		targetCount := ratio * countMean * midtones
		correction := targetCount - r.buffer.Count(i)
		adapted := math.Max(0, math.Min(r.adapt, correction))
		count += int(adapted)
	}
	return count
}

func (r *Render) next() *[sampleSize]uint {
	buffer := &[sampleSize]uint{}
	size := uint(r.camera.Width() * r.camera.Height()) // TODO: replace most of these uints with ints for simplicity
	if r.adapt > 0 {
		i := 0
		for i < sampleSize {
			end := i + 1 + r.deficit(r.cursor)
			for i < end && i < sampleSize {
				buffer[i] = r.cursor
				i++
			}
			if end <= sampleSize {
				r.cursor = (r.cursor + 1) % size
				if r.cursor == 0 {
					r.buffer.UpdateVariance()
				}
			}
		}
		return buffer
	}
	for i := 0; i < sampleSize; i++ {
		buffer[i] = r.cursor
		r.cursor = (r.cursor + 1) % size
	}
	return buffer
}

func writePng(file string, i image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, i)
	return err
}
