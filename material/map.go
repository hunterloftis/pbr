package material

import (
	"image"

	"github.com/hunterloftis/pbr/rgb"
)

type Map struct {
	color       image.Image
	metalness   image.Image
	roughness   image.Image
	specularity image.Image
	emission    image.Image
	defaults    Sample
}

func MappedMaterial(defaults Sample) *Map {
	return &Map{defaults: defaults}
}

func (m *Map) SetColor(i image.Image) {
	m.color = i
}

func (m *Map) At(u, v float64) *Sample {
	s := m.defaults
	if m.color != nil {
		w := m.color.Bounds().Max.X
		h := m.color.Bounds().Max.Y
		u2 := u * float64(w)
		v2 := -v * float64(h)
		x := int(u2) % w
		if x < 0 {
			x += w
		}
		y := int(v2) % h
		if y < 0 {
			y += h
		}
		r, g, b, _ := m.color.At(x, y).RGBA()
		s.Color = rgb.Energy{float64(r), float64(g), float64(b)}
	}
	return &s
}

func (m *Map) Emits() bool {
	return m.defaults.Emission > 0 || m.emission != nil
}
