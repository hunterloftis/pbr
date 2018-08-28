package material

import "github.com/hunterloftis/pbr/pkg/rgb"

func Plastic(r, g, b, roughness float64) *Uniform {
	return &Uniform{
		Color:       rgb.Energy{r, g, b},
		Roughness:   roughness,
		Specularity: 0.04,
	}
}
