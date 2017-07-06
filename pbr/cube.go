package pbr

import (
	"math"
)

// Cube describes a unit cube scaled, rotated, and translated by Pos.
type Cube struct {
	Pos Matrix4
	Mat Material
}

// UnitCube returns a pointer to a new 1x1x1 Cube Surface with position pos and material mat.
func UnitCube(pos Matrix4, mat Material) *Cube {
	return &Cube{
		Pos: pos,
		Mat: mat,
	}
}

// Intersect tests for an intersection
// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection
// https://tavianator.com/fast-branchless-raybounding-box-intersections/
func (c *Cube) Intersect(ray Ray3) (bool, float64) {
	i := (&c.Pos).Inverse() // global to local transform
	r := i.MultRay(ray)     // translate ray into local space
	x1 := (-0.5 - r.Origin.X) / r.Dir.X
	x2 := (0.5 - r.Origin.X) / r.Dir.X
	y1 := (-0.5 - r.Origin.Y) / r.Dir.Y
	y2 := (0.5 - r.Origin.Y) / r.Dir.Y
	z1 := (-0.5 - r.Origin.Z) / r.Dir.Z
	z2 := (0.5 - r.Origin.Z) / r.Dir.Z
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	if z1 > z2 {
		z1, z2 = z2, z1
	}
	min := math.Max(math.Max(x1, y1), z1)
	max := math.Min(math.Min(x2, y2), z2)
	if hit := min > 0 && max >= min; !hit {
		return false, 0
	}
	dist := c.Pos.MultDir(r.Dir.Scaled(min)).Len() // translate distance from local to global space
	return dist >= Bias, dist
}

// NormalAt returns the normal at this point on the surface
func (c *Cube) NormalAt(p Vector3) Vector3 {
	i := (&c.Pos).Inverse() // global to local transform
	p1 := i.MultPoint(p)    // translate point into local space
	abs := p1.Abs()
	var normal Vector3
	switch {
	case abs.X > abs.Y && abs.X > abs.Z:
		normal = Vector3{math.Copysign(1, p1.X), 0, 0}
	case abs.Y > abs.Z:
		normal = Vector3{0, math.Copysign(1, p1.Y), 0}
	default:
		normal = Vector3{0, 0, math.Copysign(1, p1.Z)}
	}
	return c.Pos.MultNormal(normal) // translate normal from local to global space
}

// MaterialAt returns the material at this point on the surface
func (c *Cube) MaterialAt(v Vector3) Material {
	return c.Mat
}
