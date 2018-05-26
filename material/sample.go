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
	Sample(out geom.Direction, rnd *rand.Rand) (in geom.Direction, pdf float64)
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

func (s *Sample) BSDF(rnd *rand.Rand) BSDF {
	if s.Transmission > 0 {
		return Microfacet{}
	}
	if rnd.Float64() <= s.Metalness {
		return Microfacet{
			Specular:   s.Color,
			Roughness:  s.Roughness,
			Multiplier: 1,
		}
	}
	if rnd.Float64() < 0.5 {
		return Microfacet{
			Specular:   rgb.Energy{s.Specularity, s.Specularity, s.Specularity},
			Roughness:  s.Roughness,
			Multiplier: 2,
		}
	}
	return Lambert{
		Color:       s.Color,
		Roughness:   s.Roughness,
		Specularity: s.Specularity,
		Multiplier:  2,
	}
}
