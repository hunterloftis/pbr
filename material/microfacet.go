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

func (m Microfacet) Sample(out, normal geom.Direction, rnd *rand.Rand) geom.Direction {
	// TODO: better sampling
	return normal.RandHemi(rnd)
}

func (m Microfacet) PDF(in, normal geom.Direction) float64 {
	// TODO: PDF that matches a better sampling distribution
	return 1 / (2 * math.Pi)
}

// https://computergraphics.stackexchange.com/questions/130/trying-to-implement-microfacet-brdf-but-my-result-images-are-wrong
// https://schuttejoe.github.io/post/ggximportancesamplingpart2/
func (m Microfacet) Eval(in, out, normal geom.Direction) rgb.Energy {
	F := schlick2(in, normal, m.F0.Mean())  // The Fresnel function
	D := ggx(in, out, normal, m.Roughness)  // The NDF (Normal Distribution Function)
	G := smithGGX(out, normal, m.Roughness) // The Geometric Shadowing function
	r := (F * D * G) / (4 * normal.Dot(in) * normal.Dot(out))
	return m.F0.Scaled(r)
}
