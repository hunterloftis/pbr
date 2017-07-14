package pbr

import (
	"math"
	"math/rand"
	"time"
)

// Sampler samples pixels for a Scene by tracing Rays from a Camera.
type Sampler struct {
	Width  int
	Height int
	SamplerConfig

	pixels []float64 // stored in a flat array of Elements
	cam    *Camera
	scene  *Scene
	noise  float64
	rnd    *rand.Rand
	count  int
}

// SamplerConfig configures a Sampler
type SamplerConfig struct {
	Bounces int
	Samples float64
	Adapt   int
}

// NewSampler constructs a new Sampler instance.
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
		conf.Adapt = 3
	}
	return &Sampler{
		Width:         cam.Width,
		Height:        cam.Height,
		SamplerConfig: conf,
		pixels:        make([]float64, cam.Width*cam.Height*Elements),
		cam:           cam,
		scene:         scene,
		rnd:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SampleFrame samples every pixel in the Camera's frame at least once.
// Depending on the Sampler's `adapt` value, noisy pixels may be sampled several times.
// It returns the total number of samples taken.
func (s *Sampler) SampleFrame() (total int) {
	noise := 0.0
	mean := s.noise + Bias
	max := s.Adapt * 3
	length := len(s.pixels)
	for p := 0; p < length; p += Elements {
		samples := s.Adaptive(s.pixels[p+Noise], mean, max)
		noise += s.sample(p, s.rnd, samples)
		total += samples
	}
	s.noise = noise / float64(s.Width*s.Height)
	s.count += total
	return
}

// PerPixel returns the per pixel sample count
func (s *Sampler) PerPixel() float64 {
	return float64(s.count) / float64(s.Width*s.Height)
}

// Adaptive returns the number of samples to take given specific and average noise values.
func (s *Sampler) Adaptive(noise float64, mean float64, max int) int {
	ratio := noise/mean + Bias
	return int(math.Min(math.Ceil(math.Pow(ratio, float64(s.Adapt))), float64(max)))
}

// sample samples a single pixel `samples` times.
// The pixel is specified by the index `p`.
func (s *Sampler) sample(p int, rnd *rand.Rand, samples int) float64 {
	x, y := s.pixelAt(p)
	before := value(s.pixels, p)
	for i := 0; i < samples; i++ {
		sample := s.trace(x, y, rnd)
		rgb := [3]float64{sample.X, sample.Y, sample.Z}
		s.pixels[p+Red] += rgb[0]
		s.pixels[p+Green] += rgb[1]
		s.pixels[p+Blue] += rgb[2]
		s.pixels[p+Count]++
	}
	after := value(s.pixels, p)
	scale := (before.Len()+after.Len())/2 + 1e-6
	noise := before.Minus(after).Len() / scale
	s.pixels[p+Noise] = noise
	return noise
}

// Pixels returns an array of float64 pixel values.
func (s *Sampler) Pixels() []float64 {
	return s.pixels
}

func value(pixels []float64, i int) Vector3 {
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

func (s *Sampler) pixelAt(i int) (x, y float64) {
	pos := i / Elements
	return float64(pos % s.Width), float64(pos / s.Width)
}
