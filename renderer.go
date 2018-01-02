package pbr

import (
	"image"
	"math"
	"math/rand"
	"runtime"
	"time"
)

// Pixel elements are stored in specific offsets.
// These constants allow easy access, eg `someFloat64Array[i + Blue]`
const (
	Red uint = iota
	Green
	Blue
	Count
	Noise
	Stride
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
	pixels       []float64 // stored in a flat array, chunked by Stride
	count        uint
	cursor       uint
	meanVariance float64
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
		pixels:       make([]float64, uint(c.Width*c.Height)*Stride),
		meanVariance: math.MaxFloat64,
		rnd:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *Renderer) Start(tick time.Duration) <-chan uint {
	r.active = true
	n := runtime.NumCPU()
	ticker := time.NewTicker(tick)
	updates := make(chan uint)
	results := make(chan Sample, n*2)
	pixels := make(chan uint)
	for i := 0; i < n; i++ {
		r.sample(pixels, results)
		r.request(pixels)
	}
	go func() {
		for {
			select {
			case res := <-results:
				r.integrate(res)
				r.request(pixels)
			case <-ticker.C:
				updates <- r.count
			default:
				if !r.active {
					close(pixels)
					close(updates)
					return
				}
			}
		}
	}()
	return updates
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

func (r *Renderer) request(pixels chan<- uint) {
	size := uint(r.Width * r.Height)
	pixels <- r.cursor
	if r.Uniform {
		r.cursor = (r.cursor + 1) % size
		return
	}
	p := r.cursor * Stride
	ratio := math.Min((r.pixels[p+Noise]+1)/(r.meanVariance+1), 5)
	rand := r.rnd.Float64() * ratio
	if rand < 0.9 {
		r.cursor++
		if r.cursor%size == 0 {
			r.cursor = 0
			r.meanVariance = 0
			for i := uint(0); i < size; i++ {
				r.meanVariance += r.pixels[p+Noise] / float64(size)
			}
		}
	}
}

func (r *Renderer) integrate(res Sample) {
	p := res.index * Stride
	rgb := [3]float64{res.energy.X, res.energy.Y, res.energy.Z}
	r.pixels[p+Red] += rgb[0]
	r.pixels[p+Green] += rgb[1]
	r.pixels[p+Blue] += rgb[2]
	r.pixels[p+Count]++
	r.count++
	r.computeNoise(res)
}

func (r *Renderer) computeNoise(res Sample) {
	if r.Uniform {
		return
	}
	p := res.index * Stride
	mean := r.average(res.index)
	variance := res.energy.Variance(mean)
	count := r.pixels[p+Count]
	oldNoise := r.pixels[p+Noise] * (count - 1) / count
	newNoise := variance / count
	r.pixels[p+Noise] = oldNoise + newNoise
}

func (r *Renderer) average(pixel uint) Energy {
	i := pixel * Stride
	c := float64(r.pixels[i+Count])
	red := r.pixels[i+Red] / c
	green := r.pixels[i+Green] / c
	blue := r.pixels[i+Blue] / c
	return Energy{red, green, blue}
}

func (r *Renderer) color(n float64) uint8 {
	return uint8(gamma(math.Min(n, 255), 2.2))
}

func gamma(n, g float64) float64 {
	return math.Pow(n/255, (1/g)) * 255
}
