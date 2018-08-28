package env

import (
	"math"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

type Gradient struct {
	Up, Down rgb.Energy
	Bias     float64
}

func NewGradient(down, up rgb.Energy, bias float64) *Gradient {
	return &Gradient{
		Down: down,
		Up:   up,
		Bias: bias,
	}
}

func (g *Gradient) At(dir geom.Dir) rgb.Energy {
	cos := dir.Dot(geom.Up)
	vertical := (1 + cos) / 2
	return g.Down.Lerp(g.Up, math.Pow(vertical, g.Bias))
}
