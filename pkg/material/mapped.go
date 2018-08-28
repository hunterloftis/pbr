package material

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Mapped struct {
	Color     image.Image
	Metalness image.Image
	Roughness image.Image
	Normal    image.Image
	Base      *Uniform
}

func NewMapped(base *Uniform) *Mapped {
	m := Mapped{
		Base: base,
	}
	return &m
}

func colToEnergy(c color.Color) rgb.Energy {
	r, g, b, _ := c.RGBA()
	return rgb.Energy{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
}

func colToFloat(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return float64(r+g+b) / (65535 * 3) // TODO: Need to take the length / square root here?
}

func (m *Mapped) At(u, v float64, in, norm geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf render.BSDF) {
	sample := *m.Base
	img := image.Image(nil)
	x := 0
	y := 0
	if m.Color != nil {
		img = m.Color
	} else if m.Roughness != nil {
		img = m.Roughness
	}
	if img != nil {
		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y
		u2 := u * float64(w)
		v2 := -v * float64(h)
		x = int(u2) % w
		if x < 0 {
			x += w
		}
		y = int(v2) % h
		if y < 0 {
			y += h
		}
	}
	if m.Color != nil {
		sample.Color = colToEnergy(m.Color.At(x, y))
	}
	// if m.Metalness != nil {
	// 	sample.Metalness = colToFloat(m.Metalness.At(x, y))
	// }
	if m.Roughness != nil {
		sample.Roughness = colToFloat(m.Roughness.At(x, y))
	}
	_, bsdf = sample.At(u, v, in, norm, rnd)
	// TODO: combine normal map at u,v with norm
	return geom.Up, bsdf
}

func (m *Mapped) Light() rgb.Energy {
	return m.Base.Light()
}

func (m *Mapped) Transmit() rgb.Energy {
	return m.Base.Transmit()
}
