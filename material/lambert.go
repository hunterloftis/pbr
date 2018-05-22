package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Lambert struct {
	R, G, B float64
}

func (l Lambert) Sample(out, normal geom.Direction, rnd *rand.Rand) geom.Direction {
	return normal.RandHemiCos(rnd)
}

func (l Lambert) PDF(in, out, normal geom.Direction) float64 {
	return in.Dot(normal) * math.Pi
}

func (l Lambert) Eval(in, out, normal geom.Direction) rgb.Energy {
	return rgb.Energy{l.R, l.G, l.B}.Scaled(math.Pi * in.Dot(normal))
}
