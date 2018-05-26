package geom

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var UnitVector3 = Vector3{1, 1, 1}

// Vector3 holds x, y, z values.
// TODO: rename to Vec3, V3, Vector, V, Vec? Used very frequently.
type Vector3 struct {
	X, Y, Z float64
}

func ArrayToVector3(a [3]float64) Vector3 {
	return Vector3{a[0], a[1], a[2]}
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
func (a Vector3) Greatest() float64 {
	return math.Max(a.X, math.Max(a.Y, a.Z))
}

// Dot returns the dot product of two vectors
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

// String returns a string representation of this vector
func (a *Vector3) String() string {
	if a == nil {
		return ""
	}
	x := strconv.FormatFloat(a.X, 'f', -1, 64)
	y := strconv.FormatFloat(a.Y, 'f', -1, 64)
	z := strconv.FormatFloat(a.Z, 'f', -1, 64)
	return strings.Join([]string{x, y, z}, ",")
}

func (a Vector3) Min(b Vector3) Vector3 {
	x := math.Min(a.X, b.X)
	y := math.Min(a.Y, b.Y)
	z := math.Min(a.Z, b.Z)
	return Vector3{x, y, z}
}

func (a Vector3) Max(b Vector3) Vector3 {
	x := math.Max(a.X, b.X)
	y := math.Max(a.Y, b.Y)
	z := math.Max(a.Z, b.Z)
	return Vector3{x, y, z}
}

func (a Vector3) Axis(n int) float64 {
	switch n {
	case 0:
		return a.X
	case 1:
		return a.Y
	default:
		return a.Z
	}
}

func (a Vector3) GreaterEqual(b Vector3) bool {
	return a.X >= b.X && a.Y >= b.Y && a.Z >= b.Z
}

func (a Vector3) LessEqual(b Vector3) bool {
	return a.X <= b.X && a.Y <= b.Y && a.Z <= b.Z
}

func (a Vector3) Array() [3]float64 {
	return [3]float64{a.X, a.Y, a.Z}
}

func (a Vector3) Projected(d Direction) Vector3 {
	b := Vector3(d)
	return b.Scaled(a.Dot(b))
}

// Unit normalizes a Vector3 into a Direction.
func (a Vector3) Unit() (Direction, bool) {
	d := a.Len()
	return Direction{a.X / d, a.Y / d, a.Z / d}, d > 0
}

// Set sets the vector from a string value
func (a *Vector3) Set(b Vector3) {
	a.X = b.X
	a.Y = b.Y
	a.Z = b.Z
}

// UnmarshalText unmarshals a byte slice into a Vector3 value
func (a *Vector3) UnmarshalText(b []byte) error {
	v, err := ParseVector3(string(b))
	if err != nil {
		return err
	}
	a.Set(v)
	return nil
}

func ParseVector3(s string) (v Vector3, err error) {
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
