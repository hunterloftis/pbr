package material

import "github.com/hunterloftis/pbr2/pkg/rgb"

func Light(r, g, b float64) *Uniform {
	c, e := rgb.Energy{r, g, b}.Compressed(1)
	return &Uniform{
		Color:    c,
		Emission: e,
	}
}

func Halogen(brightness float64) *Uniform {
	c, _ := rgb.Energy{4781, 4518, 4200}.Compressed(1)
	return &Uniform{
		Color:    c,
		Emission: brightness,
	}
}

func Daylight(brightness float64) *Uniform {
	c, _ := rgb.Energy{255, 255, 251}.Compressed(1)
	return &Uniform{
		Color:    c,
		Emission: brightness,
	}
}
