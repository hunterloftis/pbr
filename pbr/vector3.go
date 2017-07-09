package pbr

import (
	"math"
)

// Vector3 holds x, y, z
type Vector3 struct {
	X, Y, Z float64
}

// Scaled multiplies by a scalar
func (a Vector3) Scaled(n float64) Vector3 {
	return Vector3{a.X * n, a.Y * n, a.Z * n}
}

// By multiplies by a Vector3
func (a Vector3) By(b Vector3) Vector3 {
	return Vector3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Plus adds Vector3s together
func (a Vector3) Plus(b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Ave returns the average of X, Y, and Z
func (a Vector3) Ave() float64 {
	return (a.X + a.Y + a.Z) / 3
}

// Max returns the highest of X, Y, and Z
func (a Vector3) Max() float64 {
	return math.Max(a.X, math.Max(a.Y, a.Z))
}

// Dot returns the dot product of two vectors
// (which is also the cosine of the angle between them)
func (a Vector3) Dot(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross returns the cross product of vectors a and b
func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}

// Minus subtracts another vector from this one
func (a Vector3) Minus(b Vector3) Vector3 {
	return Vector3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Len finds the length of the vector
func (a Vector3) Len() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// Lerp linearly interpolates between two vectors
func (a Vector3) Lerp(b Vector3, n float64) Vector3 {
	m := 1 - n
	return Vector3{a.X*m + b.X*n, a.Y*m + b.Y*n, a.Z*m + b.Z*n}
}

// Equals compares two vectors
func (a Vector3) Equals(b Vector3) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}

// Abs converts X, Y, and Z to absolute values
func (a Vector3) Abs() Vector3 {
	return Vector3{math.Abs(a.X), math.Abs(a.Y), math.Abs(a.Z)}
}
