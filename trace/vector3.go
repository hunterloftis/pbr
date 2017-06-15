package trace

import (
	"math"
	"math/rand"
)

// Vector3 holds x, y, z
type Vector3 struct {
	X, Y, Z float64
}

// Scale multiplies by a scalar
func (a Vector3) Scale(n float64) Vector3 {
	return Vector3{a.X * n, a.Y * n, a.Z * n}
}

// Mult multiplies by a Vector3
func (a Vector3) Mult(b Vector3) Vector3 {
	return Vector3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Add adds Vector3s together
func (a Vector3) Add(b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Array converts this Vector3 to a fixed Array of length 3
func (a Vector3) Array() [3]float64 {
	return [3]float64{a.X, a.Y, a.Z}
}

// Enters returns whether this Vector is entering the plane represented by a normal Vector
func (a Vector3) Enters(b Vector3) bool {
	return b.Dot(a) < 0
}

// RandHemiCos returns a random unit vector sharing a hemisphere with this Vector with a cosine weighted distribution
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/ray.go#L28
func (a Vector3) RandHemiCos() Vector3 {
	u := rand.Float64()
	v := rand.Float64()
	r := math.Sqrt(u)
	theta := 2 * math.Pi * v
	s := a.Cross(VectorRandUnit()).Normalize()
	t := a.Cross(s)
	d := Vector3{}
	d = d.Add(s.Scale(r * math.Cos(theta)))
	d = d.Add(t.Scale(r * math.Sin(theta)))
	d = d.Add(a.Scale(math.Sqrt(1 - u)))
	return d
}

// Dot returns the dot product of two vectors
// (which is also the cosine of the angle between them)
func (a Vector3) Dot(b Vector3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// VectorRandUnit returns a random unit vector (some point on the edge of a unit sphere)
func VectorRandUnit() Vector3 {
	return VectorFromAngles(rand.Float64()*math.Pi*2, math.Asin(rand.Float64()*2-1))
}

// VectorFromAngles creates a vector based on theta and phi
func VectorFromAngles(theta, phi float64) Vector3 {
	return Vector3{math.Cos(theta) * math.Cos(phi), math.Sin(phi), math.Sin(theta) * math.Cos(phi)}
}

// Cross returns the cross product of vectors a and b
func (a Vector3) Cross(b Vector3) Vector3 {
	return Vector3{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
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
