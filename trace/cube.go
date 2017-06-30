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
	r := i.Ray(ray)
	_ = r
	// - test AABB intersection (https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection)
	// https://tavianator.com/fast-branchless-raybounding-box-intersections/

	// - compute distance
	return false, 0
}

// NormalAt returns the normal at this point on the surface
func (c *Cube) NormalAt(p Vector3) Vector3 {
	var axis Vector3
	// - translate point into local space
	p1 := c.Pos.Point(p)
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
	return c.Pos.Point(axis)
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
