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

// Refract refracts a vector based on the ratio of coefficients of refraction
func (a Vector3) Refract(b Vector3, refractA, refractB float64) (bool, Vector3) {
	ratio := refractA / refractB
	cos := b.Dot(a)
	k := 1 - ratio*ratio*(1-cos*cos)
	if k < 0 {
		return false, a
	}
	offset := b.Scale(ratio*cos + math.Sqrt(k))
	return true, a.Scale(ratio).Minus(offset).Normalize()
}

// Ave returns the average of X, Y, and Z
func (a Vector3) Ave() float64 {
	return (a.X + a.Y + a.Z) / 3
}

// Max returns the highest of X, Y, and Z
func (a Vector3) Max() float64 {
	return math.Max(a.X, math.Max(a.Y, a.Z))
}

// Cone returns a random vector within a Cone of the original vector
// size is 0-1, where 0 is the original vector and 1 is anything within the original hemisphere
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/util.go#L24
func (a Vector3) Cone(size float64, rnd *rand.Rand) Vector3 {
	u := rnd.Float64()
	v := rnd.Float64()
	theta := size * 0.5 * math.Pi * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a2 := v * 2 * math.Pi
	q := NewVectorSphere(rnd)
	s := a.Cross(q)
	t := a.Cross(s)
	d := Vector3{}
	d = d.Add(s.Scale(m1 * math.Cos(a2)))
	d = d.Add(t.Scale(m1 * math.Sin(a2)))
	d = d.Add(a.Scale(m2))
	return d.Normalize()
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
func (a Vector3) RandHemiCos(rnd *rand.Rand) Vector3 {
	u := rnd.Float64()
	v := rnd.Float64()
	r := math.Sqrt(u)
	theta := 2 * math.Pi * v
	s := a.Cross(NewVectorSphere(rnd)).Normalize()
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

// NewVectorSphere returns a random unit vector (some point on the edge of a unit sphere)
func NewVectorSphere(rnd *rand.Rand) Vector3 {
	return NewVectorAngles(rnd.Float64()*math.Pi*2, math.Asin(rnd.Float64()*2-1))
}

// NewVectorAngles creates a vector based on theta and phi
func NewVectorAngles(theta, phi float64) Vector3 {
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

// Lerp linearly interpolates between two vectors
func (a Vector3) Lerp(b Vector3, n float64) Vector3 {
	m := 1 - n
	return Vector3{a.X*m + b.X*n, a.Y*m + b.Y*n, a.Z*m + b.Z*n}
}

// Reflect reflects the vector about a normal (b)
func (a Vector3) Reflect(b Vector3) Vector3 {
	cos := b.Dot(a)
	return a.Minus(b.Scale(2 * cos)).Normalize()
}
