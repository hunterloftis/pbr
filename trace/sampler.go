package trace

import (
	"math"
)

// Sampler traces samples from light paths in a scene
type Sampler struct {
	Width   int
	Height  int
	samples []float64 // r, g, b, count
	cam     *Camera
	scene   *Scene
	bounces int
}

// NewSampler constructs a new Sampler instance
func NewSampler(cam *Camera, scene *Scene, bounces int) *Sampler {
	return &Sampler{
		Width:   cam.Width,
		Height:  cam.Height,
		samples: make([]float64, cam.Width*cam.Height*4),
		cam:     cam,
		scene:   scene,
		bounces: bounces,
	}
}

// Sample traces light paths for the full image
func (s *Sampler) Sample() {
	for i := 0; i < len(s.samples); i += 4 {
		x, y := s.offsetPixel(i)
		rgb := s.trace(x, y)
		s.samples[i] += rgb[0]
		s.samples[i+1] += rgb[1]
		s.samples[i+2] += rgb[2]
		s.samples[i+3]++
	}
}

func (s *Sampler) trace(x, y int) [3]float64 {
	ray := s.cam.Ray(x, y)
	signal := Vector3{1, 1, 1}
	energy := Vector3{0, 0, 0}

	for bounce := 0; bounce < 1; bounce++ { // bounce < s.bounces
		hit := s.scene.Intersect(ray)
		if hit {
			energy = energy.Add(Vector3{255, 255, 255}.Mult(signal))
		} else {
			energy = energy.Add(s.scene.Env(ray).Mult(signal))
		}
	}

	return energy.Array()
}

func (s *Sampler) offsetPixel(i int) (x, y int) {
	pos := i / 4
	return pos % s.Width, pos / s.Width
}

// Values gets the average sampled value at each pixel
// in a format compatible with image.RGBA.Pix
func (s *Sampler) Values() []uint8 {
	rgba := make([]uint8, s.Width*s.Height*4)
	for i := 0; i < len(s.samples); i += 4 {
		count := s.samples[i+3]
		rgba[i] = average(s.samples[i], count)
		rgba[i+1] = average(s.samples[i+1], count)
		rgba[i+2] = average(s.samples[i+2], count)
		rgba[i+3] = 255
	}
	return rgba
}

func average(total, count float64) uint8 {
	return uint8(math.Floor(total / count))
}
