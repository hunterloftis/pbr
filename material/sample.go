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
	entering := wo.Dot(geom.Up) > 0
	if entering {
		spec := rgb.Energy{s.Specularity, s.Specularity, s.Specularity}
		F0 := spec.Lerp(s.Color, s.Metalness)
		reflect := fresnelSchlick(wo.Dot(geom.Up), F0.Max())
		switch {
		case rnd.Float64() < reflect:
			return Microfacet{F0: F0, Roughness: s.Roughness}
		case s.Transmission > 0:
			// TODO: handle Microfacet transmission
			return Microfacet{F0: F0, Roughness: s.Roughness}
		default:
			return Lambert{Color: s.Color, Metalness: s.Metalness}
		}
	}
	// TODO: return/handle exits
	return Lambert{Color: s.Color, Metalness: s.Metalness}

	// return Testing{
	// 	sample:   *s,
	// 	u:        rnd.Float64(),
	// 	specular: rnd.Float64() < 0.1,
	// }

	// if s.Specularity == 0 && s.Metalness == 0 {
	// 	return Lambert{Color: s.Color}
	// }
	// spec := rgb.Energy{s.Specularity, s.Specularity, s.Specularity}
	// F0 := spec.Lerp(s.Color, s.Metalness)
	// F := fresnelSchlick(wo, geom.Up, F0.Mean())
	// c := s.Color.Lerp(rgb.Black, s.Metalness)
	// r := rnd.Float64()
	// // if F > 1 {
	// // 	// TODO: aha here's the issue.
	// // 	// This should NEVER be > 1

	// // 	fmt.Println(spec, F0.Mean(), F, c, r, wo.Dot(geom.Up))
	// // 	// panic("wtf")
	// // }
	// // fmt.Println("ok")
	// // if r < F && (s.Metalness == 1 || rnd.Float64() < 0.00001) {
	// if r < F && s.Metalness == 1 {
	// 	// if s.Metalness != 1 {
	// 	// 	fmt.Println(r, "<", F, "where F0.Mean() =", F0.Mean())
	// 	// 	fmt.Println(spec, F0.Mean(), F, c, r, s.Emission)
	// 	// 	panic("wtf")
	// 	// }
	// 	return Microfacet{F0: F0, Roughness: s.Roughness, nonmetal: s.Metalness == 0}
	// }
	// return Lambert{Color: c}
}
