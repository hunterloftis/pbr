package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Cook-Torrance microfacet model
type Testing struct {
	sample Sample
	u      float64
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
func (t Testing) Sample(wo geom.Direction, rnd *rand.Rand) geom.Direction {
	normal := geom.Up
	return normal.RandHemiCos(rnd)
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
// https://en.wikipedia.org/wiki/List_of_common_coordinate_transformations#From_Cartesian_coordinates_2
func (t Testing) PDF(wi, wo geom.Direction) float64 {
	normal := geom.Up
	return wi.Dot(normal) * math.Pi
}

// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func (t Testing) Eval(wi, wo geom.Direction) rgb.Energy {
	wg := geom.Up
	wm := wo.Half(wi)
	if wi.Y <= 0 || wi.Dot(wm) <= 0 {
		return rgb.Energy{0, 0, 0}
	}
	s := t.sample
	spec := rgb.Energy{s.Specularity, s.Specularity, s.Specularity}
	F0 := spec.Lerp(s.Color, s.Metalness)
	F := fresnelSchlick(wo, geom.Up, F0.Mean())
	c := s.Color.Lerp(rgb.Black, s.Metalness)
	if t.u < F {
		D := ggx(wi, wo, wg, s.Roughness)  // The NDF (Normal Distribution Function)
		G := smithGGX(wo, wg, s.Roughness) // The Geometric Shadowing function
		r := (D * G) / (4 * wg.Dot(wi) * wg.Dot(wo))
		return F0.Scaled(r)
	}
	return c
}
