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

func (l Lambert) Sample(out geom.Direction, rnd *rand.Rand) (geom.Direction, float64) {
	normal := geom.Up
	return normal.RandHemiCos(rnd), 0
}

func (l Lambert) PDF(in, out geom.Direction) (float64, float64) {
	normal := geom.Up
	return in.Dot(normal) * math.Pi, 0
}

func (l Lambert) Eval(in, out geom.Direction) rgb.Energy {
	normal := geom.Up
	return rgb.Energy{l.R, l.G, l.B}.Scaled(math.Pi * in.Dot(normal))
}
