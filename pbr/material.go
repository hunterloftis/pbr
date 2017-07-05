package pbr

import (
	"math"
	"math/rand"
)

// Material describes the properties of a physically based material
type Material struct {
	Color   Vector3
	Fresnel Vector3
	Light   Vector3
	Refract float64
	Opacity float64
	Gloss   float64
	Metal   float64
}

// Light constructs a new light
func Light(r, g, b float64) Material {
	return Material{Light: Vector3{r, g, b}}
}

// Plastic constructs a new plastic material
func Plastic(r, g, b float64, gloss float64) Material {
	return Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.04, 0.04, 0.04},
		Opacity: 1,
		Gloss:   gloss,
	}
}

// Metal constructs a new metal material
func Metal(r, g, b float64, gloss float64) Material {
	return Material{
		Fresnel: Vector3{r, g, b},
		Opacity: 1,
		Gloss:   gloss,
		Metal:   1,
	}
}

// Glass constructs a new glass material
func Glass(r, g, b, opacity float64, gloss float64) Material {
	return Material{
		Color:   Vector3{r, g, b},
		Fresnel: Vector3{0.04, 0.04, 0.04},
		Refract: 1.52,
		Opacity: opacity,
		Gloss:   gloss,
	}
}

// Bsdf returns next rays predicted by the material's bidirectional scattering distribution function
func (m *Material) Bsdf(normal Vector3, incident Vector3, dist float64, rnd *rand.Rand) (next bool, dir Vector3, signal Vector3) {
	if incident.Enters(normal) {
		// reflected
		reflect := m.schlick(normal, incident)
		if rnd.Float64() <= reflect.Ave() {
			tint := Vector3{1, 1, 1}.Lerp(m.Fresnel, m.Metal)
			return true, incident.Reflect(normal).Cone(1-m.Gloss, rnd), tint
		}
		// transmitted (entering)
		if rnd.Float64() > m.Opacity {
			refracted, dir := incident.Refract(normal, 1, m.Refract)
			return refracted, dir.Cone(1-m.Gloss, rnd), Vector3{1, 1, 1}
		}
		// absorbed
		if rnd.Float64() < m.Metal {
			return false, incident, Vector3{0, 0, 0}
		}
		// diffused
		return true, normal.RandHemiCos(rnd), m.Color.Scale(1 / math.Pi)
	}
	exited, dir := incident.Refract(normal.Scale(-1), m.Refract, 1)
	volume := math.Min(m.Opacity*dist*dist, 1)
	tint := Vector3{1, 1, 1}.Lerp(m.Color, volume)
	return exited, dir, tint
}

// http://blog.selfshadow.com/publications/s2015-shading-course/hoffman/s2015_pbs_physics_math_slides.pdf
// http://graphics.stanford.edu/courses/cs348b-10/lectures/reflection_i/reflection_i.pdf
func (m *Material) schlick(incident Vector3, normal Vector3) Vector3 {
	cos := incident.Scale(-1).Dot(normal)
	invFresnel := Vector3{1, 1, 1}.Minus(m.Fresnel)
	scaled := invFresnel.Scale(math.Pow(1-cos, 5))
	return m.Fresnel.Add(scaled)
}

// Emit returns the amount of light emitted
func (m *Material) Emit(normal Vector3, dir Vector3) Vector3 {
	if m.Light.Max() == 0 {
		return Vector3{}
	}
	cos := math.Max(normal.Dot(dir.Scale(-1)), 0)
	return m.Light.Scale(cos)
}
