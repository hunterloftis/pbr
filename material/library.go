package material

import "github.com/hunterloftis/pbr/rgb"

var Default = UniformMaterial(Sample{
	Color:        rgb.Energy{1, 1, 1},
	Metalness:    0,
	Roughness:    0.5,
	Specularity:  0.04,
	Emission:     0,
	Transmission: 0,
})

var Copper = UniformMaterial(Sample{
	Color:     rgb.Energy{0.95, 0.64, 0.54},
	Metalness: 1,
	Roughness: 0.2,
})

var RedLambert = UniformMaterial(Sample{
	Color:       rgb.Energy{1, 0.5, 0.5},
	Roughness:   1,
	Specularity: 0,
})

func Halogen(brightness float64) *Uniform {
	c, _ := rgb.Energy{4781, 4518, 4200}.Compressed(1)
	return UniformMaterial(Sample{
		Color:    c,
		Emission: brightness,
	})
}
