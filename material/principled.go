package material

import (
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Principled struct{}

func (p Principled) Sample(wo geom.Direction, rnd *rand.Rand) geom.Direction {
	F := fresnelSchlick(wi, wg, m.F0.Mean()) // The Fresnel function
}

func (p Principled) PDF(wi, wo geom.Direction) float64 {

}

func (p Principled) Eval(wi, wo geom.Direction) rgb.Energy {

}
