package geom

import (
	"math"
	"math/rand"
)

// Dir is a unit vector that specifies a direction in 3D space.
type Dir Vec

// Up is the positive Direction on the vertical (Y) axis.
var Up = Dir{0, 1, 0}

func SphericalDirection(theta, phi float64) (Dir, bool) {
	x := math.Sin(theta) * math.Cos(phi)
	y := math.Cos(theta)
	z := math.Sin(theta) * math.Sin(phi)
	return Vec{x, y, z}.Unit()
}

// Inv inverts a Direction.
func (a Dir) Inv() Dir {
	return Dir{-a.X, -a.Y, -a.Z}
}

// Enters returns whether this Vector is entering the plane represented by a normal Vector.
func (a Dir) Enters(normal Dir) bool {
	return normal.Dot(a) < 0
}

// Dot returns the dot product of two unit vectors, which is also the cosine of the angle between them.
func (a Dir) Dot(b Dir) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Dir) Half(b Dir) Dir {
	dir, _ := Vec(a).Plus(Vec(b)).Unit()
	return dir
}

// Refracted refracts a vector through the plane represented by a normal, based on the ratio of refraction indices.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func (a Dir) Refracted(normal Dir, indexA, indexB float64) (bool, Dir) {
	ratio := indexA / indexB
	cos := normal.Dot(a)
	k := 1 - ratio*ratio*(1-cos*cos)
	if k < 0 {
		return false, a
	}
	offset := normal.Scaled(ratio*cos + math.Sqrt(k))
	dir, _ := a.Scaled(ratio).Minus(offset).Unit()
	return true, dir
}

// Reflected reflects the vector about a normal.
// https://www.bramz.net/data/writings/reflection_transmission.pdf
func (a Dir) Reflected(normal Dir) Dir {
	cos := normal.Dot(a)
	dir, _ := Vec(a).Minus(normal.Scaled(2 * cos)).Unit()
	return dir
}

// To ensure that both face outward
func (a Dir) Reflect2(normal Dir) Dir {
	dir, _ := normal.Scaled(2).Scaled(a.Dot(normal)).Minus(Vec(a)).Unit()
	return dir
}

func (a Dir) Equals(b Dir) bool {
	return Vec(a).Equals(Vec(b))
}

// Scaled multiplies a Direction by a scalar to produce a Vector3.
func (a Dir) Scaled(n float64) Vec {
	return Vec(a).Scaled(n)
}

// Cross returns the cross product of unit vectors a and b.
func (a Dir) Cross(b Dir) (Dir, bool) {
	return Vec(a).Cross(Vec(b)).Unit()
}

// Cone returns a random vector within a cone about Direction a.
// size is 0-1, where 0 is the original vector and 1 is anything within the original hemisphere.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/util.go#L24
func (a Dir) Cone(size float64, rnd *rand.Rand) (Dir, bool) {
	u := rnd.Float64()
	v := rnd.Float64()
	theta := size * 0.5 * math.Pi * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a2 := v * 2 * math.Pi
	q := RandDirection(rnd)
	s, _ := a.Cross(q)
	t, _ := a.Cross(s)
	d := Vec{}
	d = d.Plus(s.Scaled(m1 * math.Cos(a2)))
	d = d.Plus(t.Scaled(m1 * math.Sin(a2)))
	d = d.Plus(a.Scaled(m2))
	return d.Unit()
}

// RandDirection returns a random unit vector (a point on the edge of a unit sphere).
func RandDirection(rnd *rand.Rand) Dir {
	return AngleDirection(rnd.Float64()*math.Pi*2, math.Asin(rnd.Float64()*2-1))
}

// AngleDirection creates a unit vector based on theta and phi.
// http://mathworld.wolfram.com/SphericalCoordinates.html
func AngleDirection(theta, phi float64) Dir {
	return Dir{math.Cos(theta) * math.Cos(phi), math.Sin(phi), math.Sin(theta) * math.Cos(phi)}
}

// RandHemiCos returns a random unit vector within the hemisphere of the normal direction a.
// It distributes these random vectors with a cosine weight.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/ray.go#L28
// NOTE: Added .Unit() because this doesn't always return a unit vector otherwise
func (a Dir) RandHemiCos(rnd *rand.Rand) (Dir, bool) {
	u := rnd.Float64()
	v := rnd.Float64()
	r := math.Sqrt(u)
	theta := 2 * math.Pi * v
	s, _ := a.Cross(RandDirection(rnd))
	t, _ := a.Cross(s)
	d := Vec{}
	d = d.Plus(s.Scaled(r * math.Cos(theta)))
	d = d.Plus(t.Scaled(r * math.Sin(theta)))
	d = d.Plus(a.Scaled(math.Sqrt(1 - u)))
	return d.Unit()
}

// https://stackoverflow.com/questions/5531827/random-point-on-a-given-sphere
// http://www.leadinglesson.com/dot-product-is-positive-for-vectors-in-the-same-general-direction
func (a Dir) RandHemi(rnd *rand.Rand) Dir {
	u := rnd.Float64()
	v := rnd.Float64()
	theta := 2 * math.Pi * u
	phi := math.Acos(2*v - 1)
	x := math.Sin(phi) * math.Cos(theta)
	y := math.Sin(phi) * math.Sin(theta)
	z := math.Cos(phi)
	dir, _ := Vec{x, y, z}.Unit()
	if a.Dot(dir) < 0 {
		return dir.Inv()
	}
	return dir
}

func ParseDirection(s string) (d Dir, err error) {
	v, err := ParseVec(s)
	if err != nil {
		return d, err
	}
	dir, _ := v.Unit()
	return dir, nil
}
