package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Microfacet struct {
	F0        rgb.Energy
	Roughness float64
}

func (m Microfacet) Sample(out geom.Direction, rnd *rand.Rand) geom.Direction {
	normal := geom.Up
	return normal.RandHemi(rnd)
}

func (m Microfacet) PDF(in, out geom.Direction) float64 {
	return 1 / (2 * math.Pi)
}

func (m Microfacet) Eval(in, out geom.Direction) rgb.Energy {
	normal := geom.Up
	F := schlick2(in, normal, m.F0.Mean())  // The Fresnel function
	D := ggx(in, out, normal, m.Roughness)  // The NDF (Normal Distribution Function)
	G := smithGGX(out, normal, m.Roughness) // The Geometric Shadowing function
	r := (F * D * G) / (4 * normal.Dot(in) * normal.Dot(out))
	return m.F0.Scaled(r)
}

// https://github.com/jeremypaton/wombleman/blob/master/src/bsdfs/microfacet.cpp
// func (m Microfacet) Eval(in, out, normal geom.Direction) rgb.Energy {
// 	if out.Dot(normal) <= 0 {
// 		return rgb.Energy{0, 0, 0}
// 	}
// 	a := m.Roughness // squared?
// 	wh := in.Half(out)
// 	ci := in.Z
// 	co := out.Z
// 	ho := wh.Dot(wo)
// 	term1 := m.F0.Scaled(1 / math.Pi)
// 	denom := 4 * ci * co
// 	num := beckmann(wh, a) * schlick2(ho)

// }
