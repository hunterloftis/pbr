package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Lambert struct {
	r, g, b float64
}

// TODO: remove both cosine weights if this doesn't work

func (l Lambert) Sample(in, normal geom.Direction, rnd *rand.Rand) geom.Direction {
	return normal.RandHemiCos(rnd)
}

func (l Lambert) PDF(out, normal geom.Direction) float64 {
	return out.Cos(normal) * math.Pi
}

func (l Lambert) Radiance(in, out, normal geom.Direction) rgb.Energy {
	return rgb.Energy{l.r, l.g, l.b}.Amplified(math.Pi * out.Cos(normal))
}
