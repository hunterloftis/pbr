package pbr

import (
	"math"
	"math/rand"
	"time"
)

// Sampler traces samples from light paths in a scene
type Sampler struct {
	Width   int
	Height  int
	pixels  []float64 // stored in blocks of `PROPS`
	cam     *Camera
	scene   *Scene
	bounces int
	count   int
	noise   float64
	adapt   int
}

// NewSampler constructs a new Sampler instance
func NewSampler(cam *Camera, scene *Scene, bounces int, adapt int) *Sampler {
	return &Sampler{
		Width:   cam.Width,
		Height:  cam.Height,
		pixels:  make([]float64, cam.Width*cam.Height*PROPS),
		cam:     cam,
		scene:   scene,
		bounces: bounces,
		adapt:   adapt,
	}
}

// SampleFrame samples a frame
func (s *Sampler) SampleFrame() (total int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	noise := 0.0
	mean := s.noise + 1e-6
	limit := float64(s.adapt * 3)
	for p := 0; p < len(s.pixels); p += PROPS {
		ratio := s.pixels[p+4] / mean
		adaptation := math.Floor(math.Pow(ratio, float64(s.adapt)))
		samples := 1 + int(math.Min(adaptation, limit))
		noise += s.Sample(p, rnd, samples)
		total += samples
	}
	s.noise = noise / float64(s.Width*s.Height)
	return
}

// Sample samples a pixel
func (s *Sampler) Sample(p int, rnd *rand.Rand, samples int) float64 {
	x, y := s.pixelAt(p)
	before := value(s.pixels, p)
	for i := 0; i < samples; i++ {
		sample := s.trace(x, y, rnd)
		rgb := sample.Array()
		s.pixels[p] += rgb[0]
		s.pixels[p+1] += rgb[1]
		s.pixels[p+2] += rgb[2]
		s.pixels[p+3]++
	}
	after := value(s.pixels, p)
	scale := (before.Len()+after.Len())/2 + 1e-6
	noise := before.Minus(after).Len() / scale
	s.pixels[p+4] = noise
	return noise
}

// Pixels returns an array of float64 pixel values
func (s *Sampler) Pixels() []float64 {
	return s.pixels
}

func value(pixels []float64, i int) Vector3 {
	if pixels[i+3] == 0 {
		return Vector3{}
	}
	sample := Vector3{pixels[i], pixels[i+1], pixels[i+2]}
	return sample.Scaled(1 / pixels[i+3])
}

func (s *Sampler) trace(x, y float64, rnd *rand.Rand) Vector3 {
	ray := s.cam.Ray(x, y, rnd)
	signal := Vector3{1, 1, 1}
	energy := Vector3{0, 0, 0}

	for bounce := 0; bounce < s.bounces; bounce++ {
		hit := s.scene.Intersect(ray)
		if math.IsInf(hit.Dist, 1) {
			energy = energy.Plus(s.scene.Env(ray).By(signal))
			break
		}
		if e := hit.Mat.Emit(hit.Normal, ray.Dir); e.Max() > 0 {
			energy = energy.Plus(e.By(signal))
		}
		if rnd.Float64() > signal.Max() {
			break
		}
		if next, dir, strength := hit.Mat.Bsdf(hit.Normal, ray.Dir, hit.Dist, rnd); next {
			signal = signal.Scaled(1 / signal.Max())
			ray = Ray3{hit.Point, dir}
			signal = signal.By(strength)
		} else {
			break
		}
	}
	return energy
}

func (s *Sampler) pixelAt(i int) (x, y float64) {
	pos := i / PROPS
	return float64(pos % s.Width), float64(pos / s.Width)
}
