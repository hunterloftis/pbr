package trace

import "math"

// Sampler traces samples from light paths in a scene
type Sampler struct {
	Width   int
	Height  int
	samples []uint64 // r, g, b, count
}

// NewSampler constructs a new Sampler instance
func NewSampler(width, height int) *Sampler {
	return &Sampler{
		Width:   width,
		Height:  height,
		samples: make([]uint64, width*height*4),
	}
}

// Trace traces light paths for the full image
func (s *Sampler) Trace() {
	for i := 0; i < len(s.samples); i += 4 {
		s.samples[i] += 255
		s.samples[i+1] += 0
		s.samples[i+2] += 0
		s.samples[i+3]++
	}
}

// Samples gets the average sampled value at each pixel
// in a format compatible with image.RGBA.Pix
func (s *Sampler) Samples() []uint8 {
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

func average(total, count uint64) uint8 {
	return uint8(math.Floor(float64(total) / float64(count)))
}
