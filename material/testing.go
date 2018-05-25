package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Cook-Torrance microfacet model
type Testing struct {
	sample   Sample
	u        float64
	specular bool
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
func (t Testing) Sample(wo geom.Direction, rnd *rand.Rand) geom.Direction {
	if !t.specular {
		normal := geom.Up
		return normal.RandHemiCos(rnd)
	}
	r0 := rnd.Float64()
	r1 := rnd.Float64()
	a := t.sample.Roughness
	a2 := a * a
	theta := math.Acos(math.Sqrt((1 - r0) / ((a2-1)*r0 + 1)))
	phi := 2 * math.Pi * r1
	wm := geom.SphericalDirection(theta, phi)
	wi := wo.Reflect2(wm)
	return wi
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
// https://en.wikipedia.org/wiki/List_of_common_coordinate_transformations#From_Cartesian_coordinates_2
func (t Testing) PDF(wi, wo geom.Direction) float64 {
	if !t.specular {
		normal := geom.Up
		return wi.Dot(normal) * math.Pi
	}
	wg := geom.Up
	wm := wo.Half(wi)
	a := t.sample.Roughness
	a2 := a * a
	cosTheta := wg.Dot(wm)
	exp := (a2-1)*cosTheta*cosTheta + 1
	D := a2 / (math.Pi * exp * exp)
	return (D * wm.Dot(wg)) / (4 * wo.Dot(wm))
}

// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func (t Testing) Eval(wi, wo geom.Direction) rgb.Energy {
	// return t.sample.Color
	wg := geom.Up
	wm := wo.Half(wi)
	if wi.Y <= 0 || wi.Dot(wm) <= 0 {
		return rgb.Energy{100, 0, 0}
	}
	s := t.sample
	c := s.Color.Lerp(rgb.Black, s.Metalness)
	spec := rgb.Energy{s.Specularity, s.Specularity, s.Specularity}
	F0 := spec.Lerp(s.Color, s.Metalness)
	F := fresnelSchlick(wi.Dot(wm), F0.Mean())
	if s.Metalness == 1 { //t.u < F {
		D := ggx(wi, wo, wg, s.Roughness)  // The NDF (Normal Distribution Function)
		G := smithGGX(wo, wg, s.Roughness) // The Geometric Shadowing function
		r := (F * D * G) / (4 * wg.Dot(wi) * wg.Dot(wo))
		return F0.Scaled(r)
	}
	return c
}
