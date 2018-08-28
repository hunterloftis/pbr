package obj

import (
	"math/rand"

	"github.com/hunterloftis/pbr/pkg/geom"
	"github.com/hunterloftis/pbr/pkg/render"
	"github.com/hunterloftis/pbr/pkg/rgb"
	"github.com/hunterloftis/pbr/pkg/surface"
)

type Material struct {
	Name  string
	Files []string
}

func (m *Material) At(u, v float64, in, norm geom.Dir, rnd *rand.Rand) (geom.Dir, render.BSDF) {
	return norm, surface.Lambert{}
}

func (m *Material) Light() rgb.Energy {
	return rgb.Black
}

func (m *Material) Transmit() rgb.Energy {
	return rgb.Black
}
