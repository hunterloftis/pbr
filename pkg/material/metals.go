package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

// https://i.stack.imgur.com/Q73nz.png

func Gold(roughness, metalness float64) *Uniform {
	return &Uniform{
		Color:     rgb.Energy{1, 0.86, 0.57},
		Metalness: metalness,
		Roughness: roughness,
	}
}

func Mirror(roughness float64) *Uniform {
	return &Uniform{
		Color:     rgb.Energy{0.8, 0.8, 0.8},
		Metalness: 1,
		Roughness: roughness,
	}
}

func Copper(roughness, metalness float64) *Uniform {
	return &Uniform{
		Color:     rgb.Energy{0.98, 0.82, 0.76},
		Metalness: metalness,
		Roughness: roughness,
	}
}
