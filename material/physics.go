package material

import (
	"math"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Beer's Law.
// http://www.epolin.com/converting-absorbance-transmittance
func beers(dist float64, absorb rgb.Energy) rgb.Energy {
	red := math.Exp(-absorb.X * dist)
	green := math.Exp(-absorb.Y * dist)
	blue := math.Exp(-absorb.Z * dist)
	return rgb.Energy{red, green, blue}
}

// Schlick's approximation of Fresnel
func fresnelSchlick(in, normal geom.Direction, f0 float64) float64 {
	return f0 + (1-f0)*math.Pow(1-normal.Dot(in), 5)
}

// GGX Normal Distribution Function
// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func ggx(in, out, normal geom.Direction, roughness float64) float64 {
	wm := in.Half(out)
	a := roughness * roughness
	a2 := a * a
	cosTheta := normal.Dot(wm)
	exp := (a2-1)*cosTheta*cosTheta + 1
	return a2 / (math.Pi * exp * exp)
}

// Smith geometric shadowing for a GGX distribution
// http://graphicrants.blogspot.com/2013/08/specular-brdf-reference.html
func smithGGX(out, normal geom.Direction, roughness float64) float64 {
	a := roughness * roughness
	a2 := a * a
	nv := normal.Dot(out)
	return (2 * nv) / (nv + math.Sqrt(a2+(1-a2)*nv*nv))
}
