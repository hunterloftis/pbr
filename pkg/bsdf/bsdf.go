package bsdf

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
)

// Schlick's approximation of Fresnel
// https://en.wikipedia.org/wiki/Schlick%27s_approximation
// cosTheta is the cosine between the (view or light) ray and the surface normal (or microfacet half-vector)
func fresnelSchlick(cosTheta, f0 float64) float64 {
	x := math.Pow(1-cosTheta, 5)
	return math.Max(0, math.Min(1, f0+(1-f0)*x))
}

// GGX Normal Distribution Function
// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func ggx(in, out, normal geom.Dir, roughness float64) float64 {
	wm := in.Half(out)
	a := roughness
	a2 := a * a
	cosTheta := normal.Dot(wm)
	exp := (a2-1)*cosTheta*cosTheta + 1
	return a2 / (math.Pi * exp * exp)
}

// Smith geometric shadowing for a GGX distribution
// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func smithGGX(out, normal geom.Dir, roughness float64) float64 {
	a := roughness * roughness
	a2 := a * a
	nv := normal.Dot(out)
	return (2 * nv) / (nv + math.Sqrt(a2+(1-a2)*nv*nv))
}
