package trace

// Sampler traces samples from light paths in a scene
type Sampler struct {
	Width  int
	Height int
}

// Samples returns the average sampled value at each pixel
// in a format compatible with image.RGBA.Pix
func (s *Sampler) Samples() []uint8 {
	samples := make([]uint8, s.Width*s.Height*4)
	for i := 0; i < len(samples); i += 4 {
		samples[i] = 255
		samples[i+1] = 0
		samples[i+2] = 0
		samples[i+3] = 255
	}
	return samples
}
