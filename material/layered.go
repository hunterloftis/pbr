package material

import (
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Layered struct {
	specular  BSDF
	transmit  BSDF
	diffuse   BSDF
	f0        float64
	transmits bool
}

func LayeredBSDF(s *Sample) *Layered {
	sp := s.Specularity
	f0 := rgb.Energy{sp, sp, sp}.Lerp(s.Color, s.Metalness)
	return &Layered{
		f0:        f0.Mean(),
		transmits: s.Transmission > 0,
		specular: Microfacet{
			Specularity: f0,
			Roughness:   s.Roughness,
		},
		transmit: Microfacet{},
		diffuse: Lambert{
			Color:     s.Color,
			Metalness: s.Metalness,
			Roughness: s.Roughness,
		},
	}
}

func (l *Layered) Sample(wo geom.Direction, rnd *rand.Rand) (geom.Direction, float64) {
	// f := fresnelSchlick(wo.Dot(geom.Up), l.f0)
	// if rnd.Float64() < f {
	// 	return l.specular.Sample(wo, rnd)
	// }
	// // if l.transmits {
	// // 	return l.transmit.Sample(wo, rnd)
	// // }
	return l.diffuse.Sample(wo, rnd)
}

func (l *Layered) Eval(wi, wo geom.Direction) rgb.Energy {
	light := rgb.White
	spec := light.Times(l.specular.Eval(wi, wo))
	light = light.Minus(spec)
	// light = light.Times(l.transmit.Eval(wi, wo))
	diff := light.Times(l.diffuse.Eval(wi, wo))
	return spec.Plus(diff)
}
