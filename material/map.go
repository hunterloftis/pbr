package material // TODO: Move up a level into the surface package (surface.Material)

import (
	"image"
	"math"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// Map describes the properties of a physically-based material
// Zero-value is a black, opaque, diffuse, non-metallic surface
// TODO: represent everything from MaterialDesc and stop nesting a desc inside
type Map struct {
	d            MaterialDesc
	absorbance   rgb.Energy // Initd absorbance
	refract      float64    // Initd index of refraction
	fresnel      float64    // Initd average Fresnel value (TODO: rename to "specular")
	transmission rgb.Energy // Initd "alpha" value
}

// TODO: http://saarela.github.io/ShapeToolbox/gs-material-texture.html
// http://igorsklyar.com/system/images/development_descriptions/189/disney_1.jpeg?1432292046
type MaterialDesc struct {
	Name     string
	Color    rgb.Energy // Diffuse color for opaque surfaces, transmission coefficients for transparent surfaces
	Fresnel  rgb.Energy // Fresnel coefficients, used for fresnel reflectivity and computing the refractive index
	Light    rgb.Energy // Light emittance, used if this Material is a light source
	Transmit float64    // 0 = opaque, 1 = transparent, (0-1) = tinted thin surface
	Rough    float64    // Microfacet roughness (Material "polish")
	Metal    float64    // The metallic range of electric (1) or dielectric (0), controls energy absorption
	Thin     bool       // The material is a thin, double-sided surface
	Coat     float64    // Glossy clear-coat
	Texture  image.Image
}

func New(d MaterialDesc) *Map {
	m := Map{d: d}
	m.fresnel = math.Max(geom.Vector3(d.Fresnel).Ave(), 0.02)
	if d.Thin {
		m.transmission = d.Color.Scaled(d.Transmit)
		m.absorbance = rgb.Energy{0, 0, 0} // TODO: This is confusingly named (has nothing to do with m.absorb())
		m.refract = 1
	} else {
		m.transmission = rgb.Energy{d.Transmit, d.Transmit, d.Transmit}
		m.absorbance = rgb.Energy{
			X: 2 - math.Log10(d.Color.X*100),
			Y: 2 - math.Log10(d.Color.Y*100),
			Z: 2 - math.Log10(d.Color.Z*100),
		}
		// specular = ((ior - 1) / (ior +1))^2/0.08 <-- looks like would yield .fresnel * 0.1
		// (https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html)
		m.refract = (1 + math.Sqrt(m.fresnel)) / (1 - math.Sqrt(m.fresnel))
	}
	return &m
}

func (m *Map) Name() string {
	return m.d.Name
}

// TODO: get rid of this after a Light interface exists
func (m *Map) Emit() rgb.Energy {
	return m.d.Light
}

// specular = ((ior - 1) / (ior +1))^2/0.08 <-- looks like would yield .fresnel * 0.1
// (https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html)
// http://images.slideplayer.com/42/11344425/slides/slide_97.jpg
// https://www.opengl.org/discussion_boards/showthread.php/169451-negative-texture-coords-in-wavefront-obj-format
func (m *Map) At(u, v float64) *Sample {
	f := m.d.Fresnel.Mean()
	s := &Sample{
		Color:    m.d.Color,
		Light:    m.d.Light,
		Fresnel:  m.d.Fresnel,
		Transmit: m.d.Transmit,
		Refract:  (1 + math.Sqrt(f)) / (1 - math.Sqrt(f)),
		Rough:    m.d.Rough,
		Coat:     m.d.Coat,
		Metal:    m.d.Metal,
		Thin:     m.d.Thin,
	}
	if m.d.Texture != nil {
		width := m.d.Texture.Bounds().Max.X
		height := m.d.Texture.Bounds().Max.Y
		u2 := u * float64(width)
		v2 := -v * float64(height)
		x := int(u2) % width
		if x < 0 {
			x += width
		}
		y := int(v2) % height
		if y < 0 {
			y += height
		}
		r, g, b, a := m.d.Texture.At(x, y).RGBA()
		s.Color = rgb.Energy{float64(r), float64(g), float64(b)}.Scaled(1 / float64(a))
		// TODO: deal with metals; s.Fresnel = s.Fresnel.Blend(s.Color, s.Metal)?
	}
	return s
}
