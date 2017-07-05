package pbr

import (
	"math"
)

// Cube describes a unit cube scaled, rotated, and translated by Transform
type Cube struct {
	Pos Matrix4
	Mat Material
}

// Intersect tests for an intersection
// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection
// https://tavianator.com/fast-branchless-raybounding-box-intersections/
func (c *Cube) Intersect(ray Ray3) (bool, float64) {
	i := (&c.Pos).Inverse() // global to local transform
	r := i.MultRay(ray)     // translate ray into local space
	tx1 := (-0.5 - r.Origin.X) / r.Dir.X
	tx2 := (0.5 - r.Origin.X) / r.Dir.X
	tmin, tmax := math.Min(tx1, tx2), math.Max(tx1, tx2)
	ty1 := (-0.5 - r.Origin.Y) / r.Dir.Y
	ty2 := (0.5 - r.Origin.Y) / r.Dir.Y
	tmin, tmax = math.Max(tmin, math.Min(ty1, ty2)), math.Min(tmax, math.Max(ty1, ty2))
	tz1 := (-0.5 - r.Origin.Z) / r.Dir.Z
	tz2 := (0.5 - r.Origin.Z) / r.Dir.Z
	tmin, tmax = math.Max(tmin, math.Min(tz1, tz2)), math.Min(tmax, math.Max(tz1, tz2))
	if hit := tmin > 0 && tmax >= tmin; !hit {
		return false, 0
	}
	dist := c.Pos.MultDir(r.Dir.Scale(tmin)).Length() // translate distance from local to global space
	return dist >= BIAS, dist
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
