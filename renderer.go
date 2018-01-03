package pbr

import (
	"image"
	"math"
	"math/rand"
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
	rnd    *rand.Rand

	// state
	active       bool
	count        uint
	cursor       uint
	meanVariance float64
	image        Image
}

// RenderConfig configures rendering settings.
type RenderConfig struct {
	Uniform  bool
	Bounces  int
	Direct   uint // TODO
	Indirect uint // TODO
}

type Sample struct {
	index  uint
	energy Energy
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
		camera:       c,
		scene:        s,
		meanVariance: math.MaxFloat64,
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		image:        NewImage(uint(c.Width), uint(c.Height)),
	}
}

func (r *Renderer) Start(tick time.Duration) <-chan uint {
	r.active = true
	n := runtime.NumCPU()
	ticker := time.NewTicker(tick)
	update := make(chan uint)
	result := make(chan Sample, n*2)
	pixel := make(chan uint, n*1024) // TODO: calc the actual limit this can be given # of workers, image size, etc
	planned := uint(0)
	for i := 0; i < n; i++ {
		r.sample(pixel, result)
		planned += r.next(pixel)
	}
	go func() {
		for {
			select {
			case res := <-result:
				r.image.Integrate(res.index, res.energy, !r.Uniform)
				r.count++
				for planned < r.count+uint(n) {
					planned += r.next(pixel)
				}
			case <-ticker.C:
				update <- r.count
			default:
				if !r.active {
					close(pixel)
					close(update)
					return
				}
			}
		}
	}()
	return update
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

func (r *Renderer) Rgb(expose float64) image.Image {
	return r.image.Rgb(expose)
}

func (r *Renderer) Heat() image.Image {
	return r.image.Heat(Count)
}

func (r *Renderer) Noise() image.Image {
	return r.image.Heat(Noise)
}

func (r *Renderer) sample(in <-chan uint, out chan<- Sample) {
	size := uint(r.Width * r.Height)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for {
			if p, ok := <-in; ok {
				i := p % size
				x, y := r.pixelAt(i)
				sample := r.trace(x, y, rnd)
				out <- Sample{i, sample}
			} else {
				return
			}
		}
	}()
}

func (r *Renderer) pixelAt(i uint) (x, y float64) {
	return float64(i % uint(r.Width)), float64(i / uint(r.Width))
}

func (r *Renderer) trace(x, y float64, rnd *rand.Rand) Energy {
	ray := r.camera.ray(x, y, rnd)
	signal := Energy{1, 1, 1}
	energy := Energy{0, 0, 0}

	for bounce := 0; bounce < r.Bounces; bounce++ {
		hit, surface, dist := r.scene.Intersect(ray)
		if !hit {
			energy = energy.Merged(r.scene.Env(ray), signal)
			break
		}
		point := ray.Moved(dist)
		normal, mat := surface.At(point)
		energy = energy.Merged(mat.Emit(normal, ray.Dir), signal)
		signal = signal.RandomGain(rnd) // "Russian Roulette"
		if signal == (Energy{}) {
			break
		}
		if next, dir, str := mat.Bsdf(normal, ray.Dir, dist, rnd); next {
			signal = signal.Strength(str)
			ray = Ray3{point, dir}
		} else {
			break
		}
	}
	return energy
}

// TODO: make this more sophisticated (like using max-mean variance vs just max)
func (r *Renderer) next(pixels chan<- uint) uint {
	count := uint(1)
	if !r.Uniform {
		noise := r.image.Noise(r.cursor * Stride) // TODO: shouldn't have to calc with Stride
		ratio := (noise + 1) / (r.image.MaxVariance() + 1)
		scale := uint(math.Max(math.Min(ratio, 100), 0)) // TODO: remove the magic number
		count += scale
	}
	for i := uint(0); i < count; i++ {
		pixels <- r.cursor
	}
	r.cursor = (r.cursor + 1) % uint(r.Width*r.Height)
	if r.cursor == 0 {
		r.image.UpdateVariance()
	}
	return count
}
