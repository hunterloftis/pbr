package material

import "github.com/hunterloftis/pbr/rgb"

var Default = UniformMaterial(Sample{
	Color:        rgb.Energy{1, 1, 1},
	Metalness:    0,
	Roughness:    0.1,
	Specularity:  0.02,
	Emission:     0,
	Transmission: 0,
})

var Copper = UniformMaterial(Sample{
	Color:     rgb.Energy{0.95, 0.64, 0.54},
	Metalness: 1,
	Roughness: 0.2,
})

var Mirror = UniformMaterial(Sample{
	Color:     rgb.Energy{1, 1, 1},
	Metalness: 1,
	Roughness: 0.01,
})

var WhiteLambert = UniformMaterial(Sample{
	Color:       rgb.Energy{1, 1, 1},
	Roughness:   1,
	Specularity: 0,
})

var RedLambert = UniformMaterial(Sample{
	Color:       rgb.Energy{1, 0.5, 0.5},
	Roughness:   1,
	Specularity: 0,
})

var RedPlastic = UniformMaterial(Sample{
	Color:       rgb.Energy{1, 0, 0},
	Roughness:   0.1,
	Specularity: 0.02,
})

var TealPlastic = UniformMaterial(Sample{
	Color:       rgb.Energy{0, 1, 1},
	Roughness:   0.1,
	Specularity: 0.02,
})

var ShinyPlastic = UniformMaterial(Sample{
	Color:       rgb.Energy{1, 1, 1},
	Roughness:   0.01,
	Specularity: 1,
})

func Halogen(brightness float64) *Uniform {
	c, _ := rgb.Energy{4781, 4518, 4200}.Compressed(1)
	return UniformMaterial(Sample{
		Color:    c,
		Emission: brightness,
	})
}

func Gold(roughness float64) *Uniform {
	return UniformMaterial(Sample{
		Color:     rgb.Energy{1.022, 0.782, 0.344},
		Metalness: 1,
		Roughness: roughness,
	})
}
