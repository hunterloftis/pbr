package trace

import "math"

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

// Bsdf returns next rays predicted by the material's bidirectional scattering distribution function
func (m *Material) Bsdf(normal Vector3, dir Vector3, dist float64) {

}

// Emit returns the amount of light emitted
func (m *Material) Emit(normal Vector3, dir Vector3) Vector3 {
	cos := math.Max(normal.Dot(dir.Scale(-1)), 0)
	return m.Light.Scale(cos)
}
