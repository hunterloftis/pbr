package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Cook-Torrance microfacet model
type Microfacet struct {
	F0        rgb.Energy
	Roughness float64
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
func (m Microfacet) Sample(wo geom.Direction, rnd *rand.Rand) geom.Direction {
	r0 := rnd.Float64()
	r1 := rnd.Float64()
	a := m.Roughness
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
func (m Microfacet) PDF(wi, wo geom.Direction) float64 {
	wg := geom.Up
	wm := wo.Half(wi)
	a := m.Roughness
	a2 := a * a
	cosTheta := wg.Dot(wm)
	exp := (a2-1)*cosTheta*cosTheta + 1
	D := a2 / (math.Pi * exp * exp)
	return (D * wm.Dot(wg)) / (4 * wo.Dot(wm))
}

// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func (m Microfacet) Eval(wi, wo geom.Direction) rgb.Energy {
	wg := geom.Up
	wm := wo.Half(wi)
	if wi.Y <= 0 || wi.Dot(wm) <= 0 {
		return rgb.Energy{0, 0, 0}
	}
	F := fresnelSchlick(wi, wg, m.F0.Mean()) // The Fresnel function
	D := ggx(wi, wo, wg, m.Roughness)        // The NDF (Normal Distribution Function)
	G := smithGGX(wo, wg, m.Roughness)       // The Geometric Shadowing function
	r := (F * D * G) / (4 * wg.Dot(wi) * wg.Dot(wo))
	return m.F0.Scaled(r)
}
