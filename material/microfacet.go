package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Just copper for now (0.98, 0.82, 0.76)
type Microfacet struct {
	F0        rgb.Energy
	Roughness float64
}

func (m Microfacet) Sample(out, normal geom.Direction, rnd *rand.Rand) geom.Direction {
	return normal.RandHemiCos(rnd)
}

func (m Microfacet) Probability(in, normal geom.Direction) float64 {
	return 1 / (in.Cos(normal) * math.Pi)
}

// https://computergraphics.stackexchange.com/questions/130/trying-to-implement-microfacet-brdf-but-my-result-images-are-wrong
func (m Microfacet) Radiance(in, out, normal geom.Direction) rgb.Energy {
	F := schlick2(in, normal, m.F0.Average()) // The Fresnel function
	D := ggx(in, out, normal, m.Roughness)    // The NDF (Normal Distribution Function)
	G := smithGGX(out, normal, m.Roughness)   // The Geometric Shadowing function
	r := (F * D * G) / (4 * normal.Cos(in) * normal.Cos(out))
	return m.F0.Amplified(r)
}
