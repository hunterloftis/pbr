package pbr

import (
	"math"
	"math/rand"
	"runtime"
	"time"
)

// Sampler samples pixels for a Scene by tracing Rays from a Camera.
type Sampler struct {
	Width  int
	Height int
	SamplerConfig

	samples   []float64 // stored in a flat array, chunked by Stride
	cam       *Camera
	scene     *Scene
	count     int
	meanNoise float64
}

// SamplerConfig configures a Sampler.
type SamplerConfig struct {
	Bounces int
	Samples float64
	Adapt   int
}

type sampleStat struct {
	count int
	noise float64
}

// NewSampler constructs a new Sampler for a given Camera and Scene.
// The Sampler samples Rays from the Camera into the Scene.
// bounces specifies the maximum number of times a Ray can bounce around the scene (eg, 10).
// adapt specifies how adaptive sampling should be to noise (0 = none, 3 = medium, 4 = high).
func NewSampler(cam *Camera, scene *Scene, config ...SamplerConfig) *Sampler {
	conf := SamplerConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	if conf.Bounces == 0 {
		conf.Bounces = 10 // Reasonable default
	}
	if conf.Samples == 0 {
		conf.Samples = math.Inf(1) // Sample forever by default
	}
	if conf.Adapt == 0 { // TODO: 0 should be a valid value
		conf.Adapt = 4
	}
	return &Sampler{
		Width:         cam.Width,
		Height:        cam.Height,
		SamplerConfig: conf,
		samples:       make([]float64, index(cam.Width*cam.Height)*Stride),
		cam:           cam,
		scene:         scene,
	}
}

// Sample samples every pixel in the Camera's frame at least once.
// Depending on the Sampler's `adapt` value, noisy pixels may be sampled several times.
// It returns the total number of samples taken.
// TODO: clean this up a bit
// https://stackoverflow.com/questions/22517614/golang-concurrent-array-access
func (s *Sampler) Sample() {
	length := index(len(s.samples))
	workers := index(runtime.NumCPU())
	ch := make(chan sampleStat, workers)

	for i := index(0); i < workers; i++ {
		go func(i index, adapt, max int, mean float64) {
			var stat sampleStat
			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			for p := i * Stride; p < length; p += Stride * workers {
				samples := adaptive(s.samples[p+Noise], adapt, max, mean)
				stat.noise += s.samplePixel(p, rnd, samples)
				stat.count += samples
			}
			ch <- stat
		}(i, s.Adapt, s.Adapt*3, s.meanNoise+BIAS)
	}

	var sample sampleStat
	for i := index(0); i < workers; i++ {
		stat := <-ch
		sample.count += stat.count
		sample.noise += stat.noise
	}
	s.count += sample.count
	s.meanNoise = sample.noise / float64(sample.count)
}

// Count returns the total sample count.
func (s *Sampler) Count() int {
	return s.count
}

// PerPixel returns the per-pixel sample count.
func (s *Sampler) PerPixel() float64 {
	return float64(s.count) / float64(s.Width*s.Height)
}

// Samples returns an array of float64 pixel values.
func (s *Sampler) Samples() []float64 {
	return s.samples
}

// Pixels returns the number of pixels per sample.
func (s *Sampler) Pixels() int {
	return s.Width * s.Height
}

// adaptive returns the number of samples to take given specific and average noise values.
// TODO: this is a pretty slow function (think it's inlined), speed it up
func adaptive(noise float64, adapt, max int, mean float64) int {
	ratio := noise/mean + BIAS
	return int(math.Min(math.Ceil(math.Pow(ratio, float64(adapt))), float64(max)))
}

// samplePixel samples a single pixel `samples` times.
// The pixel is specified by the index `p`.
func (s *Sampler) samplePixel(p index, rnd *rand.Rand, samples int) float64 {
	x, y := s.pixelAt(p)
	before := value(s.samples, p)
	for i := 0; i < samples; i++ {
		sample := s.trace(x, y, rnd)
		rgb := [3]float64{sample.X, sample.Y, sample.Z}
		s.samples[p+Red] += rgb[0]
		s.samples[p+Green] += rgb[1]
		s.samples[p+Blue] += rgb[2]
		s.samples[p+Count]++
	}
	after := value(s.samples, p)
	scale := (before.Len()+after.Len())/2 + 1e-6
	noise := before.Minus(after).Len() / scale
	s.samples[p+Noise] = noise
	return noise
}

func value(pixels []float64, i index) Vector3 {
	if pixels[i+Count] == 0 {
		return Vector3{}
	}
	sample := Vector3{pixels[i+Red], pixels[i+Green], pixels[i+Blue]}
	return sample.Scaled(1 / pixels[i+Count])
}

func (s *Sampler) trace(x, y float64, rnd *rand.Rand) Energy {
	ray := s.cam.ray(x, y, rnd)
	signal := Energy{1, 1, 1}
	energy := Energy{0, 0, 0}

	for bounce := 0; bounce < s.Bounces; bounce++ {
		hit, surface, dist := s.scene.Intersect(ray)
		if !hit {
			energy = energy.Merged(s.scene.Env(ray), signal)
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

func (s *Sampler) pixelAt(i index) (x, y float64) {
	pos := int(i / Stride)
	return float64(pos % s.Width), float64(pos / s.Width)
}
