package trace

import (
	"math"
	"math/rand"
)

const adapt = 0.25

// Sampler traces samples from light paths in a scene
type Sampler struct {
	Width   int
	Height  int
	pixels  []float64 // r, g, b, count
	cam     *Camera
	scene   *Scene
	bounces int
	count   int
}

// NewSampler constructs a new Sampler instance
func NewSampler(cam *Camera, scene *Scene, bounces int) *Sampler {
	return &Sampler{
		Width:   cam.Width,
		Height:  cam.Height,
		pixels:  make([]float64, cam.Width*cam.Height*4),
		cam:     cam,
		scene:   scene,
		bounces: bounces,
	}
}

// Sample traces light paths for the full image
func (s *Sampler) Sample() {
	total := float64(s.Width * s.Height)
	for p := 0; p < len(s.pixels); p += 4 {
		val := s.value(p)
		x, y := s.offsetPixel(p)
		average := math.Floor(float64(s.count) / total)
		limit := int(average) + 1
		for j := 0; j < limit; j++ {
			sample := s.trace(x, y)
			variance := sample.Minus(val).Length() / (val.Length() + 1e-6)
			rgb := sample.Array()
			val = sample
			s.pixels[p] += rgb[0]
			s.pixels[p+1] += rgb[1]
			s.pixels[p+2] += rgb[2]
			s.pixels[p+3]++
			s.count++
			if variance < adapt {
				break
			}
		}
	}
}

func (s *Sampler) value(i int) Vector3 {
	sample := Vector3{s.pixels[i], s.pixels[i+1], s.pixels[i+2]}
	return sample.Scale(1 / s.pixels[i+3])
}

func (s *Sampler) trace(x, y int) Vector3 {
	ray := s.cam.Ray(x, y)
	signal := Vector3{1, 1, 1}
	energy := Vector3{0, 0, 0}

	for bounce := 0; bounce < s.bounces; bounce++ {
		intersected, hit := s.scene.Intersect(ray)
		if !intersected {
			energy = energy.Add(s.scene.Env(ray).Mult(signal))
			break
		}
		light := hit.Mat.Emit(hit.Normal, ray.Dir)
		energy = energy.Add(light.Mult(signal))
		if rand.Float64() > signal.Max() {
			break
		}
		signal = signal.Scale(1 / signal.Max())
		next, dir, strength := hit.Mat.Bsdf(hit.Normal, ray.Dir, hit.Dist)
		if !next {
			break
		}
		ray = Ray3{hit.Point, dir}
		signal = signal.Mult(strength)
	}

	return energy
}

func (s *Sampler) offsetPixel(i int) (x, y int) {
	pos := i / 4
	return pos % s.Width, pos / s.Width
}

// Values gets the average sampled rgb at each pixel
func (s *Sampler) Values() []float64 {
	rgb := make([]float64, s.Width*s.Height*3)
	for p := 0; p < len(s.pixels); p += 4 {
		val := s.value(p).Array()
		i := p / 4 * 3
		rgb[i] = val[0]
		rgb[i+1] = val[1]
		rgb[i+2] = val[2]
	}
	return rgb
}

// Counts returns the sample count at each pixel as rgb
func (s *Sampler) Counts() []float64 {
	rgb := make([]float64, s.Width*s.Height*3)
	var max float64
	for p := 0; p < len(s.pixels); p += 4 {
		max = math.Max(max, s.pixels[p+3])
	}
	for p := 0; p < len(s.pixels); p += 4 {
		val := (s.pixels[p+3] / max) * 255
		i := p / 4 * 3
		rgb[i] = val
		rgb[i+1] = val
		rgb[i+2] = val
	}
	return rgb
}
