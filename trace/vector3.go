package trace

import "math"

// Vector3 holds x, y, z
type Vector3 struct {
	X, Y, Z float64
}

// Scale multiplies by a scalar
func (a Vector3) Scale(n float64) Vector3 {
	return Vector3{a.X * n, a.Y * n, a.Z * n}
}

// Dot returns the dot product of two vectors
// (which is also the cosine of the angle between them)
func (a Vector3) Dot(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Minus subtracts another vector from this one
func (a Vector3) Minus(b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Normalize makes the vector of length 1
func (a Vector3) Normalize() Vector3 {
	d := a.Length()
	return Vector3{a.X / d, a.Y / d, a.Z / d}
}

// Length finds the length of the vector
func (a Vector3) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}
