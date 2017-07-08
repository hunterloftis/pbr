package pbr

import (
	"math"
)

// Cube describes a unit cube scaled, rotated, and translated by Pos.
type Cube struct {
	Pos  *Matrix4
	Mat  *Material
	Grid *Material
}

// UnitCube returns a pointer to a new 1x1x1 Cube Surface with position pos and material mat.
func UnitCube(pos *Matrix4, mat *Material, grid *Material) *Cube {
	return &Cube{
		Pos:  pos,
		Mat:  mat,
		Grid: grid,
	}
}

// Intersect tests for an intersection between a Ray3 and this Cube
// It returns whether there was an intersection (bool) and the intersection distance along the ray (float64)
// Both the Ray3 and the distance are in world space.
// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection
// https://tavianator.com/fast-branchless-raybounding-box-intersections/
func (c *Cube) Intersect(ray Ray3) (bool, float64) {
	inv := c.Pos.Inverse() // global to local transform
	r := inv.MultRay(ray)  // translate ray into local space
	or := [3]float64{r.Origin.X, r.Origin.Y, r.Origin.Z}
	dir := [3]float64{r.Dir.X, r.Dir.Y, r.Dir.Z}
	t0 := 0.0
	t1 := math.Inf(1)
	for i := 0; i < 3; i++ {
		tNear := (-0.5 - or[i]) / dir[i]
		tFar := (0.5 - or[i]) / dir[i]
		if tNear > tFar {
			tNear, tFar = tFar, tNear
		}
		if tNear > t0 {
			t0 = tNear
		}
		if tFar < t1 {
			t1 = tFar
		}
		if t0 > t1 {
			return false, 0
		}
	}
	if t0 > 0 {
		if dist := c.Pos.MultDir(r.Dir.Scaled(t0)).Len(); dist >= Bias {
			return true, dist
		}
	}
	if t1 > 0 {
		if dist := c.Pos.MultDir(r.Dir.Scaled(t1)).Len(); dist >= Bias {
			return true, dist
		}
	}
	return false, 0
}

// At returns the normal Vector3 at this point on the Surface
func (c *Cube) At(p Vector3) (Vector3, *Material) {
	i := c.Pos.Inverse() // global to local transform
	p1 := i.MultPoint(p) // translate point into local space
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
	// translate normal from local to global space
	mat := c.Mat
	if c.Grid != nil {
		x, z := p.X*20, p.Z*20
		if dx := math.Abs(x - math.Floor(x)); dx < 0.03 {
			mat = c.Grid
		} else if dz := math.Abs(z - math.Floor(z)); dz < 0.03 {
			mat = c.Grid
		}
	}
	return c.Pos.MultNormal(normal), mat
}
