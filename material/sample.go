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

func (s *Sample) BSDF(wo geom.Direction, rnd *rand.Rand) BSDF {
	spec := rgb.Energy{s.Specularity, s.Specularity, s.Specularity}
	F0 := spec.Lerp(s.Color, s.Metalness)
	F := fresnelSchlick(wo, geom.Up, F0.Mean())
	c := s.Color.Lerp(rgb.Black, s.Metalness)
	if rnd.Float64() < F {
		return Microfacet{F0: F0, Roughness: s.Roughness}
	}
	return Lambert{Color: c}
}
