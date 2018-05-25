package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Lambert struct {
	Color     rgb.Energy
	Metalness float64
}

func (l Lambert) Sample(out geom.Direction, rnd *rand.Rand) geom.Direction {
	normal := geom.Up
	return normal.RandHemiCos(rnd)
}

func (l Lambert) PDF(in, out geom.Direction) float64 {
	normal := geom.Up
	return in.Dot(normal) * math.Pi
}

func (l Lambert) Eval(in, out geom.Direction) rgb.Energy {
	return l.Color.Lerp(rgb.Black, l.Metalness)
}
