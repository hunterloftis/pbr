package bsdf

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Ignore struct{}

func (i Ignore) Sample(wo geom.Dir, rnd *rand.Rand) (geom.Dir, float64, bool) {
	return wo.Inv(), 1, false
}

func (i Ignore) PDF(wi, wo geom.Dir) float64 {
	return 1
}

func (i Ignore) Eval(wi, wo geom.Dir) rgb.Energy {
	return rgb.White
}
