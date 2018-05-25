package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

// TODO: Oren-Nayar for roughness
type Lambert struct {
	Color       rgb.Energy
	Roughness   float64
	Specularity float64
}

func (l Lambert) Sample(out geom.Direction, rnd *rand.Rand) geom.Direction {
	normal := geom.Up
	return normal.RandHemiCos(rnd)
}

func (l Lambert) PDF(in, out geom.Direction) float64 {
	normal := geom.Up
	return in.Dot(normal) * math.Pi
}

func (l Lambert) Eval(wi, wo geom.Direction) rgb.Energy {
	wm := wo.Half(wi)
	F := fresnelSchlick(wi.Dot(wm), l.Specularity) // TODO: half-vector or normal (geom.Up)?
	return l.Color.Plus(rgb.Energy{F, F, F}).Limit(1)
}
