package trace

// Vector3 holds x, y, z
type Vector3 struct {
	X, Y, Z float64
}

// Scale multiplies by a scalar
func (v *Vector3) Scale(n float64) *Vector3 {
	return &Vector3{v.X * n, v.Y * n, v.Z * n}
}
