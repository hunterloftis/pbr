package pbr

import (
	"math"
	"math/rand"
)

// Material describes the properties of a physically-based material
// Zero-value is a black, opaque, diffuse, non-metallic surface
type Material struct {
	d          MaterialDesc
	absorbance Energy  // Initd absorbance
	refract    float64 // Initd index of refraction
	fresnel    float64 // Initd average Fresnel value
}

// TODO: "Thin" is a hack to get rid of.
type MaterialDesc struct {
	Color    Energy  // Diffuse color for opaque surfaces, transmission coefficients for transparent surfaces
	Fresnel  Energy  // Fresnel coefficients, used for fresnel reflectivity and computing the refractive index
	Light    Energy  // Light emittance, used if this Material is a light source
	Transmit float64 // 0 = opaque, 1 = transparent, (0-1) = tinted thin surface
	Gloss    float64 // Microsurface roughness (Material "polish")
	Metal    float64 // The metallic range of electric (1) or dielectric (0), controls energy absorption
	Thin     bool    // Should transparent surfaces be passed through (instead of entered and refracted)
}

// Light constructs a new light
// r, g, b (0-Inf) specifies the light color
func Light(r, g, b float64) *Material {
	return NewMaterial(MaterialDesc{
		Light: Energy{r, g, b},
	})
}

// Plastic constructs a new plastic material
// r, g, b (0-1) controls the color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Plastic(r, g, b float64, gloss float64) *Material {
	return NewMaterial(MaterialDesc{
		Color:   Energy{r, g, b},
		Fresnel: Energy{0.04, 0.04, 0.04},
		Gloss:   gloss,
	})
}

// Lambert constructs a new plastic material
// r, g, b (0-1) controls the color
func Lambert(r, g, b float64) *Material {
	return NewMaterial(MaterialDesc{
		Color:   Energy{r, g, b},
		Fresnel: Energy{0.02, 0.02, 0.02},
	})
}

// Metal constructs a new metal material
// r, g, b (0-1) controls the fresnel color
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Metal(r, g, b float64, gloss float64) *Material {
	return NewMaterial(MaterialDesc{
		Fresnel: Energy{r, g, b},
		Gloss:   gloss,
		Metal:   1,
	})
}

// Glass constructs a new glass material
// r, g, b (0-1) controls the transmission of the glass (beer-lambert)
// gloss (0-1) controls the microfacet roughness (how polished the surface looks)
func Glass(r, g, b, gloss float64) *Material {
	return NewMaterial(MaterialDesc{
		Color:    Energy{r, g, b},
		Fresnel:  Energy{0.042, 0.042, 0.042},
		Transmit: 1,
		Gloss:    gloss,
	})
}

func NewMaterial(d MaterialDesc) *Material {
	m := Material{d: d}
	m.fresnel = math.Max(Vector3(d.Fresnel).Ave(), 0.02)
	if d.Transmit > 0 && d.Thin {
		// TODO: this is a hack. Behavior should be baked into the material model instead of special casing.
		// Maybe make a Material interface instead of a struct?
		m.refract = 1
		m.absorbance = Energy{0, 0, 0}
	} else {
		m.absorbance = Energy{
			X: 2 - math.Log10(d.Color.X*100),
			Y: 2 - math.Log10(d.Color.Y*100),
			Z: 2 - math.Log10(d.Color.Z*100),
		}
		m.refract = (1 + math.Sqrt(m.fresnel)) / (1 - math.Sqrt(m.fresnel))
	}
	return &m
}

func (m *Material) Description() *MaterialDesc {
	return &m.d
}

// Bsdf is an attempt at a new bsdf
func (m *Material) Bsdf(norm, inc Direction, dist float64, rnd *rand.Rand) (bool, Direction, Energy) {
	if inc.Enters(norm) {
		reflect := schlick(norm, inc, m.fresnel, 0, 0)
		switch {
		// reflect
		case rnd.Float64() < reflect:
			return m.reflect(norm, inc, rnd)
		// transmit (in)
		case rnd.Float64() < m.d.Transmit:
			return m.transmit(norm, inc, rnd)
		// absorb
		case rnd.Float64() < m.d.Metal:
			return m.absorb(inc)
		// diffuse
		default:
			return m.diffuse(norm, inc, rnd)
		}
	}
	// transmit (out)
	return m.exit(norm, inc, dist, rnd)
}

// Emit returns the amount of light emitted from the Material at a given angle.
func (m *Material) Emit(normal, dir Direction) Energy {
	cos := math.Max(normal.Cos(dir.Inv()), 0)
	return m.d.Light.Amplified(cos)
}

func (m *Material) reflect(norm, inc Direction, rnd *rand.Rand) (bool, Direction, Energy) {
	// TODO: if reflection enters the normal, invert the reflection about the normal
	if refl := inc.Reflected(norm).Cone(1-m.d.Gloss, rnd); !refl.Enters(norm) {
		return true, refl, Energy(Vector3{1, 1, 1}.Lerp(Vector3(m.d.Fresnel), m.d.Metal))
	}
	return m.diffuse(norm, inc, rnd)
}

func (m *Material) transmit(norm, inc Direction, rnd *rand.Rand) (bool, Direction, Energy) {
	if entered, refr := inc.Refracted(norm, 1, m.refract); entered {
		if spread := refr.Cone(1-m.d.Gloss, rnd); spread.Enters(norm) {
			return true, spread, Energy{1, 1, 1}
		}
		return true, refr, Energy{1, 1, 1}
	}
	return m.diffuse(norm, inc, rnd)
}

func (m *Material) exit(norm, inc Direction, dist float64, rnd *rand.Rand) (bool, Direction, Energy) {
	if m.d.Transmit == 0 {
		// shallow bounce within margin of error
		// isn't really an intersection, so just keep the ray moving
		return true, inc, Energy{1, 1, 1}
	}
	if rnd.Float64() >= schlick(norm, inc, 0, m.refract, 1.0) {
		if exited, refr := inc.Refracted(norm.Inv(), m.refract, 1); exited {
			if spread := refr.Cone(1-m.d.Gloss, rnd); !spread.Enters(norm) {
				return true, spread, beers(dist, m.absorbance)
			}
			return true, refr, beers(dist, m.absorbance)
		}
	}
	return true, inc.Reflected(norm.Inv()), beers(dist, m.absorbance)
}

func (m *Material) diffuse(norm, inc Direction, rnd *rand.Rand) (bool, Direction, Energy) {
	return true, norm.RandHemiCos(rnd), m.d.Color.Amplified(1 / math.Pi)
}

func (m *Material) absorb(inc Direction) (bool, Direction, Energy) {
	return false, inc, Energy{}
}

// Schlick's approximation.
// Returns a number between 0-1 indicating the percentage of light reflected vs refracted.
// 0 = no reflection, all refraction; 1 = 100% reflection, no refraction.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
// http://blog.selfshadow.com/publications/s2015-shading-course/hoffman/s2015_pbs_physics_math_slides.pdf
// http://graphics.stanford.edu/courses/cs348b-10/lectures/reflection_i/reflection_i.pdf
func schlick(incident, normal Direction, r0, n1, n2 float64) float64 {
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
func beers(dist float64, absorb Energy) Energy {
	red := math.Exp(-absorb.X * dist)
	green := math.Exp(-absorb.Y * dist)
	blue := math.Exp(-absorb.Z * dist)
	return Energy{red, green, blue}
}
