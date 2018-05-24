package material

type Uniform struct {
	Sample
}

func UniformMaterial(s Sample) *Uniform {
	return &Uniform{s}
}

func (un *Uniform) At(u, v float64) *Sample {
	return &un.Sample
}

func (un *Uniform) Emits() bool {
	return un.Sample.Emission > 0
}
