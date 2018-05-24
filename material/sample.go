package material

import (
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Description interface {
	At(u, v float64) *Sample
	Emits() bool
}

type BSDF interface {
	Sample(out geom.Direction, rnd *rand.Rand) (in geom.Direction)
	PDF(in, out geom.Direction) float64
	Eval(in, out geom.Direction) rgb.Energy
}

type Sample struct {
	Color        rgb.Energy
	Metalness    float64
	Roughness    float64
	Specularity  float64
	Emission     float64
	Transmission float64
}

func (s *Sample) Light() rgb.Energy {
	return s.Color.Scaled(s.Emission)
}

func (s *Sample) BSDF() BSDF {
	if s.Metalness > 0 {
		return Microfacet{F0: s.Color, Roughness: s.Roughness}
	}
	return Lambert{R: s.Color.X, G: s.Color.Y, B: s.Color.Z}
}
