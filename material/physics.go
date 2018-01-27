package material

import (
	"math"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Schlick's approximation.
// Returns a number between 0-1 indicating the percentage of light reflected vs refracted.
// 0 = no reflection, all refraction; 1 = 100% reflection, no refraction.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
// http://blog.selfshadow.com/publications/s2015-shading-course/hoffman/s2015_pbs_physics_math_slides.pdf
// http://graphics.stanford.edu/courses/cs348b-10/lectures/reflection_i/reflection_i.pdf
func schlick(incident, normal geom.Direction, r0, n1, n2 float64) float64 {
	cosX := -normal.Cos(incident)
	if r0 == 0 {
		r0 = (n1 - n2) / (n1 + n2)
		r0 *= r0
		if n1 > n2 {
			n := n1 / n2
			sinT2 := n * n * (1.0 - cosX*cosX)
			if sinT2 > 1.0 {
				return 1.0 // Total Internal Reflection
			}
			cosX = math.Sqrt(1.0 - sinT2)
		}
	}
	x := 1.0 - cosX
	return r0 + (1.0-r0)*x*x*x*x*x
}

// Beer's Law.
// http://www.epolin.com/converting-absorbance-transmittance
func beers(dist float64, absorb rgb.Energy) rgb.Energy {
	red := math.Exp(-absorb.X * dist)
	green := math.Exp(-absorb.Y * dist)
	blue := math.Exp(-absorb.Z * dist)
	return rgb.Energy{red, green, blue}
}
