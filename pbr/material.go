package pbr

import (
	"math"
	"math/rand"
)

// Material describes the properties of a physically-based material
type Material struct {
	Color      Vector3 // Diffuse color for opaque surfaces, transmission coefficients for transparent surfaces
	Fresnel    Vector3 // Fresnel coefficients, used for fresnel reflectivity
	Light      Vector3 // Light emittance, used if this Material is a light source
	Refract    float64 // Index of refraction
	Opacity    float64 // 0 = transparent, 1 = opaque, (0-1) = tinted thin surface
	Gloss      float64 // Microsurface roughness (how "polished" is this Material)
	Metal      float64 // The metallic range of electric (1) or dielectric (0), controls energy absorption
	absorbance Vector3 // cache for computed absorbance
}

// Light constructs a new light
// r, g, b (0-Inf) specifies the light color
func Light(r, g, b float64) *Material {
	return &Material{Light: Vector3{r, g, b}}
}

// DayLight constructs a new light with a DayLight color temperature.
func DayLight() *Material {
	return Light(2550, 2550, 2510)
}

// Plastic constructs a new plastic material
// r, g, b (0-1) controls the color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Plastic(r, g, b float64, gloss float64) *Material {
	return &Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.04, 0.04, 0.04},
		Opacity: 1,
		Gloss:   gloss,
	}
}

// Lambert constructs a new plastic material
// r, g, b (0-1) controls the color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Lambert(r, g, b float64, gloss float64) *Material {
	return &Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.02, 0.02, 0.02},
		Opacity: 1,
		Gloss:   gloss,
	}
}

// Metal constructs a new metal material
// r, g, b (0-1) controls the fresnel color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Metal(r, g, b float64, gloss float64) *Material {
	return &Material{
		Fresnel: Vector3{r, g, b},
		Opacity: 1,
		Gloss:   gloss,
		Metal:   1,
	}
}

// Glass constructs a new glass material
// r, g, b (0-1) controls the transmission of the glass (beer-lambert)
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Glass(r, g, b, gloss float64) *Material {
	return &Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.042, 0.042, 0.042},
		Refract: 1.514,
		Opacity: 0,
		Gloss:   gloss,
		absorbance: Vector3{
			X: 2 - math.Log10(r*100),
			Y: 2 - math.Log10(g*100),
			Z: 2 - math.Log10(b*100),
		},
	}
}

// Bsdf is an attempt at a new bsdf
func (m *Material) Bsdf(norm, inc Vector3, dist float64, rnd *rand.Rand) (bool, Vector3, Vector3) {
	if inc.Enters(norm) {
		reflect := schlick(norm, inc, m.Fresnel.Ave(), 0, 0) // TODO: cache m.Fresnel.Ave() as m.fresnel
		switch {
		// reflect
		case rnd.Float64() < reflect:
			return m.reflect(norm, inc, rnd)
		// transmit (in)
		case rnd.Float64() >= m.Opacity:
			return m.transmit(norm, inc, rnd)
		// absorb
		case rnd.Float64() < m.Metal:
			return m.absorb(norm, inc)
		// diffuse
		default:
			return m.diffuse(norm, inc, rnd)
		}
	}
	// transmit (out)
	return m.exit(norm, inc, dist, rnd)
}

// Emit returns the amount of light emitted from the Material at a given angle.
func (m *Material) Emit(normal Vector3, dir Vector3) (bool, Vector3) {
	if m.Light.Max() == 0 { // TODO: cache m.Light.Max() as m.light
		return false, Vector3{}
	}
	cos := math.Max(normal.Dot(dir.Scaled(-1)), 0) // instead of scaling -1, can I invert the normal?
	return true, m.Light.Scaled(cos)
}

func (m *Material) reflect(norm, inc Vector3, rnd *rand.Rand) (bool, Vector3, Vector3) {
	if refl := inc.Reflected(norm).Cone(1-m.Gloss, rnd); !refl.Enters(norm) {
		return true, refl, Vector3{1, 1, 1}.Lerp(m.Fresnel, m.Metal)
	}
	return m.diffuse(norm, inc, rnd)
}

func (m *Material) transmit(norm, inc Vector3, rnd *rand.Rand) (bool, Vector3, Vector3) {
	if entered, refr := inc.Refracted(norm, 1, m.Refract); entered {
		if spread := refr.Cone(1-m.Gloss, rnd); spread.Enters(norm) {
			return true, spread, Vector3{1, 1, 1}
		}
		return true, refr, Vector3{1, 1, 1}
	}
	return m.diffuse(norm, inc, rnd)
}

func (m *Material) exit(norm, inc Vector3, dist float64, rnd *rand.Rand) (bool, Vector3, Vector3) {
	if m.Opacity == 1 {
		return false, inc, Vector3{}
	}
	if rnd.Float64() >= schlick(norm, inc, 0, m.Refract, 1.0) {
		if exited, refr := inc.Refracted(norm.Scaled(-1), m.Refract, 1); exited {
			if spread := refr.Cone(1-m.Gloss, rnd); !spread.Enters(norm) {
				return true, spread, beers(dist, m.absorbance)
			}
			return true, refr, beers(dist, m.absorbance)
		}
	}
	return true, inc.Reflected(norm.Scaled(-1)), beers(dist, m.absorbance)
}

func (m *Material) diffuse(norm, inc Vector3, rnd *rand.Rand) (bool, Vector3, Vector3) {
	return true, norm.RandHemiCos(rnd), m.Color.Scaled(1 / math.Pi)
}

func (m *Material) absorb(norm, inc Vector3) (bool, Vector3, Vector3) {
	return false, inc, Vector3{}
}

// Schlick's approximation.
// Returns a number between 0-1 indicating the percentage of light reflected vs refracted.
// 0 = no reflection, all refraction; 1 = 100% reflection, no refraction.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
// http://blog.selfshadow.com/publications/s2015-shading-course/hoffman/s2015_pbs_physics_math_slides.pdf
// http://graphics.stanford.edu/courses/cs348b-10/lectures/reflection_i/reflection_i.pdf
func schlick(incident, normal Vector3, r0, n1, n2 float64) float64 {
	cosX := -normal.Dot(incident)
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
func beers(dist float64, absorb Vector3) Vector3 {
	red := math.Exp(-absorb.X * dist)
	green := math.Exp(-absorb.Y * dist)
	blue := math.Exp(-absorb.Z * dist)
	return Vector3{red, green, blue}
}
