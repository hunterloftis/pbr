package pbr

import (
	"math"
	"math/rand"
)

// Direction is a unit vector that specifies a direction in 3D space.
type Direction Vector3

// Unit normalizes a Vector3 into a Direction.
func (a Vector3) Unit() Direction {
	d := a.Len()
	return Direction{a.X / d, a.Y / d, a.Z / d}
}

// Inv inverts a Direction.
func (a Direction) Inv() Direction {
	return Direction{-a.X, -a.Y, -a.Z}
}

// Enters returns whether this Vector is entering the plane represented by a normal Vector.
func (a Direction) Enters(normal Direction) bool {
	return normal.Cos(a) < 0
}

// Cos returns the dot product of two unit vectors, which is also the cosine of the angle between them.
func (a Direction) Cos(b Direction) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Refracted refracts a vector through the plane represented by a normal, based on the ratio of refraction indices.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func (a Direction) Refracted(normal Direction, indexA, indexB float64) (bool, Direction) {
	ratio := indexA / indexB
	cos := normal.Cos(a)
	k := 1 - ratio*ratio*(1-cos*cos)
	if k < 0 {
		return false, a
	}
	offset := normal.Scaled(ratio*cos + math.Sqrt(k))
	return true, a.Scaled(ratio).Minus(offset).Unit()
}

// Reflected reflects the vector about a normal.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func (a Direction) Reflected(normal Direction) Direction {
	cos := normal.Cos(a)
	return Vector3(a).Minus(normal.Scaled(2 * cos)).Unit()
}

// Scaled multiplies a Direction by a scalar to produce a Vector3.
func (a Direction) Scaled(n float64) Vector3 {
	return Vector3(a).Scaled(n)
}

// Cross returns the cross product of unit vectors a and b.
func (a Direction) Cross(b Direction) Direction {
	return Direction{a.Y*b.Z - a.Z*b.Y, a.Z*b.X - a.X*b.Z, a.X*b.Y - a.Y*b.X}
}

// Cone returns a random vector within a cone about Direction a.
// size is 0-1, where 0 is the original vector and 1 is anything within the original hemisphere.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/util.go#L24
func (a Direction) Cone(size float64, rnd *rand.Rand) Direction {
	u := rnd.Float64()
	v := rnd.Float64()
	theta := size * 0.5 * math.Pi * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a2 := v * 2 * math.Pi
	q := RandDirection(rnd)
	s := a.Cross(q)
	t := a.Cross(s)
	d := Vector3{}
	d = d.Plus(s.Scaled(m1 * math.Cos(a2)))
	d = d.Plus(t.Scaled(m1 * math.Sin(a2)))
	d = d.Plus(a.Scaled(m2))
	return d.Unit()
}

// RandDirection returns a random unit vector (a point on the edge of a unit sphere).
func RandDirection(rnd *rand.Rand) Direction {
	return AngleDirection(rnd.Float64()*math.Pi*2, math.Asin(rnd.Float64()*2-1))
}

// AngleDirection creates a unit vector based on theta and phi.
func AngleDirection(theta, phi float64) Direction {
	return Direction{math.Cos(theta) * math.Cos(phi), math.Sin(phi), math.Sin(theta) * math.Cos(phi)}
}

// RandHemiCos returns a random unit vector within the hemisphere of the normal direction a.
// It distributes these random vectors with a cosine weight.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/ray.go#L28
// NOTE: Added .Unit() because this doesn't always return a unit vector otherwise
func (a Direction) RandHemiCos(rnd *rand.Rand) Direction {
	u := rnd.Float64()
	v := rnd.Float64()
	r := math.Sqrt(u)
	theta := 2 * math.Pi * v
	s := a.Cross(RandDirection(rnd))
	t := a.Cross(s)
	d := Vector3{}
	d = d.Plus(s.Scaled(r * math.Cos(theta)))
	d = d.Plus(t.Scaled(r * math.Sin(theta)))
	d = d.Plus(a.Scaled(math.Sqrt(1 - u)))
	return d.Unit()
}
