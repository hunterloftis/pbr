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
func (m Microfacet) Sample(wo geom.Direction, rnd *rand.Rand) (geom.Direction, float64) {
	r0 := rnd.Float64()
	r1 := rnd.Float64()
	a := m.Roughness * m.Roughness
	a2 := a * a
	theta := math.Acos(math.Sqrt((1 - r0) / ((a2-1)*r0 + 1)))
	phi := 2 * math.Pi * r1
	x := math.Sin(theta) * math.Cos(phi)
	y := math.Cos(theta)
	z := math.Sin(theta) * math.Sin(phi)
	wm := geom.Vector3{x, y, z}.Unit()
	wi := wo.Reflect2(wm)
	return wi, theta
}

// https://schuttejoe.github.io/post/ggximportancesamplingpart1/
// https://agraphicsguy.wordpress.com/2015/11/01/sampling-microfacet-brdf/
// https://en.wikipedia.org/wiki/List_of_common_coordinate_transformations#From_Cartesian_coordinates_2
func (m Microfacet) PDF(wi, wo geom.Direction) (float64, float64) {
	wm := wo.Half(wi)
	a := m.Roughness * m.Roughness
	a2 := a * a
	theta := math.Atan2(math.Sqrt(wm.X*wm.X+wm.Z*wm.Z), wm.Y)
	cosTheta := math.Cos(theta)
	num := a2 * cosTheta * math.Sin(theta)
	exp := (a2-1)*cosTheta*cosTheta + 1
	den := math.Pi * (exp * exp)
	return num / den, theta
}

// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func (m Microfacet) Eval(wi, wo geom.Direction) rgb.Energy {
	normal := geom.Up
	wm := wo.Half(wi)
	if wi.Y <= 0 || wi.Dot(wm) <= 0 {
		return rgb.Energy{0, 0, 0}
	}
	F := schlick2(wi, normal, m.F0.Mean()) // The Fresnel function
	D := ggx(wi, wo, normal, m.Roughness)  // The NDF (Normal Distribution Function)
	G := smithGGX(wo, normal, m.Roughness) // The Geometric Shadowing function
	r := (F * D * G) / (4 * normal.Dot(wi) * normal.Dot(wo))
	return m.F0.Scaled(r)
}

// https://github.com/schuttejoe/ShootyEngine/blob/6a301e9f7d2a46db3d1f9bc846f3637ce876a06f/Source/Applications/PathTracer/Source/PathTracerShading.cpp#L440-L445
// func (m Microfacet) Eval(wi, wo geom.Direction, rnd *rand.Rand) rgb.Energy {
// 	if wi.Y <= 0 {
// 		return rgb.Energy{0, 0, 0}
// 	}
// 	wm := wi.Half(wo)
// 	a2 := m.Roughness * m.Roughness
// 	// F := schlick3(m.F0, wi.Dot(wm))
// 	F := schlick2(wi, wm, m.F0.Mean())
// 	if rnd.Float64() < F {
// 		return rgb.Energy{0, 0, 0}
// 	}
// 	G1 := smithGGXMasking(wo, wm, a2)
// 	G2 := smithGGXMaskingShading(wi, wo, wm, a2)
// 	weight := G2 / G1
// 	return m.F0.Scaled(weight)
// }

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
