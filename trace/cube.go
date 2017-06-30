package trace

import "math"

// Cube describes a unit cube scaled, rotated, and translated by Transform
type Cube struct {
	Pos Matrix4
	Mat Material
}

// Intersect tests for an intersection
func (c *Cube) Intersect(ray Ray3) (hit bool, dist float64) {
	// - translate ray into local space with s.Transform
	i := c.Pos.Inverse()
	_ = i
	r := ray // i.Ray(ray)
	// - test AABB intersection (https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection)
	// https://tavianator.com/fast-branchless-raybounding-box-intersections/
	tx1 := (-0.5 - r.Origin.X) / r.Dir.X
	tx2 := (0.5 - r.Origin.X) / r.Dir.X

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (-0.5 - r.Origin.Y) / r.Dir.Y
	ty2 := (0.5 - r.Origin.Y) / r.Dir.Y

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (-0.5 - r.Origin.Z) / r.Dir.Z
	tz2 := (0.5 - r.Origin.Z) / r.Dir.Z

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	hit = tmax > 0 && tmax > tmin
	dist = tmin
	return
}

// NormalAt returns the normal at this point on the surface
func (c *Cube) NormalAt(p Vector3) Vector3 {
	var axis Vector3
	// - translate point into local space
	p1 := p // p1 := c.Pos.Point(p)
	// - test x, y, and z to see which one is largest/smallest
	x := math.Abs(p1.X)
	y := math.Abs(p1.Y)
	z := math.Abs(p1.Z)
	if x > y && x > z {
		axis = Vector3{sign(x), 0, 0}
	} else if y > z {
		axis = Vector3{0, sign(y), 0}
	} else {
		axis = Vector3{0, 0, sign(z)}
	}
	// - translate that axis normal back into world space
	// TODO: one of these needs to be inverted
	return axis // c.Pos.Point(axis)
}

func sign(n float64) float64 {
	if n > 0 {
		return 1
	}
	return -1
}

// MaterialAt returns the material at this point on the surface
func (c *Cube) MaterialAt(v Vector3) Material {
	return c.Mat
}
