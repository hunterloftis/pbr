package geom

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var Origin = Vec{0, 0, 0}

// Vec holds x, y, z values.
type Vec struct {
	X, Y, Z float64
}

func ArrayToVec(a [3]float64) Vec {
	return Vec{a[0], a[1], a[2]}
}

// Scaled multiplies by a scalar
func (a Vec) Scaled(n float64) Vec {
	return Vec{a.X * n, a.Y * n, a.Z * n}
}

// By multiplies by a Vector3
func (a Vec) By(b Vec) Vec {
	return Vec{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Plus adds Vector3s together
func (a Vec) Plus(b Vec) Vec {
	return Vec{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Ave returns the average of X, Y, and Z
func (a Vec) Ave() float64 {
	return (a.X + a.Y + a.Z) / 3
}

// Max returns the highest of X, Y, and Z
func (a Vec) Greatest() float64 {
	return math.Max(a.X, math.Max(a.Y, a.Z))
}

// Dot returns the dot product of two vectors
func (a Vec) Dot(b Vec) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Cross returns the cross product of vectors a and b
func (a Vec) Cross(b Vec) Vec {
	return Vec{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}

// Minus subtracts another vector from this one
func (a Vec) Minus(b Vec) Vec {
	return Vec{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Len finds the length of the vector
func (a Vec) Len() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// Lerp linearly interpolates between two vectors
func (a Vec) Lerp(b Vec, n float64) Vec {
	m := 1 - n
	return Vec{a.X*m + b.X*n, a.Y*m + b.Y*n, a.Z*m + b.Z*n}
}

// Equals compares two vectors
func (a Vec) Equals(b Vec) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}

// Abs converts X, Y, and Z to absolute values
func (a Vec) Abs() Vec {
	return Vec{math.Abs(a.X), math.Abs(a.Y), math.Abs(a.Z)}
}

// String returns a string representation of this vector
func (a *Vec) String() string {
	if a == nil {
		return ""
	}
	x := strconv.FormatFloat(a.X, 'f', -1, 64)
	y := strconv.FormatFloat(a.Y, 'f', -1, 64)
	z := strconv.FormatFloat(a.Z, 'f', -1, 64)
	return strings.Join([]string{x, y, z}, ",")
}

func (a Vec) Min(b Vec) Vec {
	x := math.Min(a.X, b.X)
	y := math.Min(a.Y, b.Y)
	z := math.Min(a.Z, b.Z)
	return Vec{x, y, z}
}

func (a Vec) Max(b Vec) Vec {
	x := math.Max(a.X, b.X)
	y := math.Max(a.Y, b.Y)
	z := math.Max(a.Z, b.Z)
	return Vec{x, y, z}
}

func (a Vec) Axis(n int) float64 {
	switch n {
	case 0:
		return a.X
	case 1:
		return a.Y
	default:
		return a.Z
	}
}

func (a Vec) GreaterEqual(b Vec) bool {
	return a.X >= b.X && a.Y >= b.Y && a.Z >= b.Z
}

func (a Vec) LessEqual(b Vec) bool {
	return a.X <= b.X && a.Y <= b.Y && a.Z <= b.Z
}

func (a Vec) Array() [3]float64 {
	return [3]float64{a.X, a.Y, a.Z}
}

func (a Vec) Projected(d Dir) Vec {
	b := Vec(d)
	return b.Scaled(a.Dot(b))
}

// Unit normalizes a Vector3 into a Direction.
func (a Vec) Unit() (Dir, bool) {
	d := a.Len()
	return Dir{a.X / d, a.Y / d, a.Z / d}, d > 0
}

// Set sets the vector from a string value
func (a *Vec) Set(b Vec) {
	a.X = b.X
	a.Y = b.Y
	a.Z = b.Z
}

// UnmarshalText unmarshals a byte slice into a Vector3 value
func (a *Vec) UnmarshalText(b []byte) error {
	v, err := ParseVec(string(b))
	if err != nil {
		return err
	}
	a.Set(v)
	return nil
}

func ParseVec(s string) (v Vec, err error) {
	xyz := strings.Split(s, ",")
	if len(xyz) != 3 {
		return v, fmt.Errorf("pbr: 3 values required for Vector3, received %v", len(xyz))
	}
	v.X, err = strconv.ParseFloat(xyz[0], 64)
	if err != nil {
		return v, err
	}
	v.Y, err = strconv.ParseFloat(xyz[1], 64)
	if err != nil {
		return v, err
	}
	v.Z, err = strconv.ParseFloat(xyz[2], 64)
	return v, err
}
