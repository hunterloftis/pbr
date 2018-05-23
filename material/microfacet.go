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

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
func (m Microfacet) Sample(wo geom.Direction, rnd *rand.Rand) geom.Direction {
	r0 := rnd.Float64()
	r1 := rnd.Float64()
	a := m.Roughness
	a2 := a * a
	theta := math.Acos(math.Sqrt((1 - r0) / ((a2-1)*r0 + 1)))
	phi := 2 * math.Pi * r1
	x := math.Sin(theta) * math.Cos(phi)
	y := math.Cos(theta)
	z := math.Sin(theta) * math.Sin(phi)
	wm := geom.Vector3{x, y, z}.Unit()
	wi := wo.Reflect2(wm)
	return wi
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
// https://en.wikipedia.org/wiki/List_of_common_coordinate_transformations#From_Cartesian_coordinates_2
func (m Microfacet) PDF(wi, wo geom.Direction) float64 {
	wm := wo.Half(wi)
	a := m.Roughness
	a2 := a * a
	theta := math.Atan2(math.Sqrt(wm.X*wm.X+wm.Z*wm.Z), wm.Y)
	cosTheta := math.Cos(theta)
	num := a2 * cosTheta * math.Sin(theta)
	exp := (a2-1)*cosTheta*cosTheta + 1
	den := math.Pi * (exp * exp)
	return num / den
}

// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func (m Microfacet) Eval(wi, wo geom.Direction) rgb.Energy {
	normal := geom.Up
	wm := wo.Half(wi)
	if wi.Y <= 0 || wi.Dot(wm) <= 0 {
		return rgb.Energy{0, 0, 0}
	}
	F := fresnelSchlick(wi, normal, m.F0.Mean()) // The Fresnel function
	D := ggx(wi, wo, normal, m.Roughness)        // The NDF (Normal Distribution Function)
	G := smithGGX(wo, normal, m.Roughness)       // The Geometric Shadowing function
	r := (F * D * G) / (4 * normal.Dot(wi) * normal.Dot(wo))
	return m.F0.Scaled(r)
}
