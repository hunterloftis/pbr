package pbr

import (
	"math"
	"math/rand"
)

// Material describes the properties of a physically-based material
type Material struct {
	Color   Vector3 // Diffuse color for opaque surfaces, transmission coefficients for transparent surfaces
	Fresnel Vector3 // Fresnel coefficients, used for fresnel reflectivity
	Light   Vector3 // Light emittance, used if this Material is a light source
	Refract float64 // Index of refraction
	Opacity float64 // 0 = transparent, 1 = opaque, (0-1) = tinted thin surface
	Gloss   float64 // Microsurface roughness (how "polished" is this Material)
	Metal   float64 // The metallic range of electric (1) or dielectric (0), controls energy absorption
}

// Light constructs a new light
// r, g, b (0-Inf) specifies the light color
func Light(r, g, b float64) Material {
	return Material{Light: Vector3{r, g, b}}
}

// DayLight constructs a new light with a DayLight color temperature.
func DayLight() Material {
	return Light(2550, 2550, 2510)
}

// Plastic constructs a new plastic material
// r, g, b (0-1) controls the color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Plastic(r, g, b float64, gloss float64) Material {
	return Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.04, 0.04, 0.04},
		Opacity: 1,
		Gloss:   gloss,
	}
}

// Metal constructs a new metal material
// r, g, b (0-1) controls the fresnel color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Metal(r, g, b float64, gloss float64) Material {
	return Material{
		Fresnel: Vector3{r, g, b},
		Opacity: 1,
		Gloss:   gloss,
		Metal:   1,
	}
}

// Glass constructs a new glass material
// r, g, b (0-1) controls the transmission of the glass (beer-lambert)
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Glass(r, g, b, gloss float64) Material {
	return Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.04, 0.04, 0.04},
		Refract: 1.52,
		Opacity: 0,
		Gloss:   gloss,
	}
}

// Bsdf returns next rays predicted by the Material's
// Bidirectional Scattering Distribution Function
func (m *Material) Bsdf(normal Vector3, incident Vector3, dist float64, rnd *rand.Rand) (next bool, dir Vector3, signal Vector3) {
	if incident.Enters(normal) {
		// reflected
		reflect := m.schlick(normal, incident)
		if rnd.Float64() <= reflect.Ave() {
			tint := Vector3{1, 1, 1}.Lerp(m.Fresnel, m.Metal)
			refl := incident.Reflected(normal).Cone(1-m.Gloss, rnd)
			if refl.Enters(normal) {
				refl = normal.RandHemiCos(rnd) // If cone passes into surface, diffuse instead
			}
			return true, refl, tint
		}
		// transmitted (entering)
		if rnd.Float64() >= m.Opacity {
			refracted, dir := incident.Refracted(normal, 1, m.Refract)
			if refracted {
				spread := dir.Cone(1-m.Gloss, rnd)
				if spread.Enters(normal) {
					return true, spread, Vector3{1, 1, 1}
				}
				return true, dir, Vector3{1, 1, 1}
			}
		}
		// absorbed
		if rnd.Float64() < m.Metal {
			return false, incident, Vector3{0, 0, 0}
		}
		// diffused
		return true, normal.RandHemiCos(rnd), m.Color.Scaled(1 / math.Pi)
	}
	if m.Opacity == 1 {
		return false, incident, Vector3{}
		// Rays shouldn't be exiting from opaque surfaces, to test:
		// panic("Exit from opaque Surface")
	}
	reflect := m.schlick2(normal, incident)
	if rnd.Float64() >= reflect { // refract
		exited, dir := incident.Refracted(normal.Scaled(-1), m.Refract, 1)
		if exited {
			return true, dir, m.beers(dist)
		}
	}
	return true, incident.Reflected(normal.Scaled(-1)), m.beers(dist) // internal reflection
}

// Emit returns the amount of light emitted
// from the Material at a given angle
func (m *Material) Emit(normal Vector3, dir Vector3) Vector3 {
	if m.Light.Max() == 0 {
		return Vector3{}
	}
	cos := math.Max(normal.Dot(dir.Scaled(-1)), 0)
	return m.Light.Scaled(cos)
}

// http://blog.selfshadow.com/publications/s2015-shading-course/hoffman/s2015_pbs_physics_math_slides.pdf
// http://graphics.stanford.edu/courses/cs348b-10/lectures/reflection_i/reflection_i.pdf
func (m *Material) schlick(incident, normal Vector3) Vector3 {
	cos := incident.Scaled(-1).Dot(normal)
	invFresnel := Vector3{1, 1, 1}.Minus(m.Fresnel)
	scaled := invFresnel.Scaled(math.Pow(1-cos, 5))
	return m.Fresnel.Plus(scaled)
}

// Schlick returns a number between 0-1 indicating the percentage of light reflected vs refracted.
// 0 = no reflection, all refraction; 1 = 100% reflection, no refraction.
// TODO: remove m.Fresnel, convert m.Refract to a Vector3, compute IoRs backwards from known Fresnel coefficients
// unify both reflect and refract to use schlick2
// https://www.bramz.net/data/writings/reflection_transmission.pdf
// Unify: http://www.visual-barn.com/f0-converting-substance-fresnel-vray-values/
func (m *Material) schlick2(incident, normal Vector3) float64 {
	n1 := m.Refract
	n2 := 1.0
	r0 := (n1 - n2) / (n1 + n2)
	r0 *= r0
	cosX := -normal.Dot(incident)
	if n1 > n2 {
		n := n1 / n2
		sinT2 := n * n * (1.0 - cosX*cosX)
		if sinT2 > 1.0 {
			return 1.0 // Total Internal Reflection
		}
		cosX = math.Sqrt(1.0 - sinT2)
	}
	x := 1.0 - cosX
	return r0 + (1.0-r0)*x*x*x*x*x
}

// https://www.bramz.net/data/writings/reflection_transmission.pdf
// func (m *Material) fresnel(incident Vector3, normal Vector3) float64 {

// }

// Beer's Law
// http://www.epolin.com/converting-absorbance-transmittance
// TODO: cache the absorbance values to avoid repeating these calculations
func (m *Material) beers(dist float64) Vector3 {
	ar := 2 - math.Log10(m.Color.X*100)
	ag := 2 - math.Log10(m.Color.Y*100)
	ab := 2 - math.Log10(m.Color.Z*100)
	red := math.Exp(-ar * dist)
	green := math.Exp(-ag * dist)
	blue := math.Exp(-ab * dist)
	return Vector3{red, green, blue}
}
