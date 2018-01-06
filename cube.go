package pbr

import (
	"math"
)

// Cube describes the orientation and material of a unit cube
type Cube struct {
	Pos      *Matrix4
	Mat      *Material
	GridMat  *Material
	GridSize float64
	box      *Box
}

// UnitCube returns a pointer to a new 1x1x1 Cube Surface with material and optional transforms.
func UnitCube(m *Material, transforms ...*Matrix4) *Cube {
	pos := Identity()
	for _, t := range transforms { // TODO: factor this so all surfaces can share it
		pos = pos.Mult(t)
	}
	c := &Cube{
		Pos: pos,
		Mat: m,
	}
	min := c.Pos.MultPoint(Vector3{-1, -1, -1})
	max := c.Pos.MultPoint(Vector3{1, 1, 1})
	c.box = NewBox(min, max)
	return c
}

// SetGrid adds a second material to the cube which is applied as a grid across its surface
func (c *Cube) SetGrid(mat *Material, size float64) *Cube {
	c.GridMat = mat
	c.GridSize = size
	return c
}

// TODO: unify with Box.Check?
func (c *Cube) Intersect(ray *Ray3) Hit {
	ok, _ := c.box.Check(ray)
	if !ok {
		return Miss
	}
	inv := c.Pos.Inverse() // global to local transform
	r := inv.MultRay(ray)  // translate ray into local space
	dir := Vector3(r.Dir).Array()
	or := r.Origin.Array()
	min := Vector3{-0.5, -0.5, -0.5}.Array()
	max := Vector3{0.5, 0.5, 0.5}.Array()
	tmin := 0.0
	tmax := math.Inf(1)
	for a := 0; a < 3; a++ {
		invD := 1 / dir[a]
		t0 := (min[a] - or[a]) * invD
		t1 := (max[a] - or[a]) * invD
		if invD < 0 {
			t0, t1 = t1, t0
		}
		if t0 > tmin {
			tmin = t0
		}
		if t1 < tmax {
			tmax = t1
		}
		if tmax < tmin {
			return Miss
		}
	}
	// TODO: lots of calculations here going from point to dist. optimize.
	pointLocal := r.Moved(tmin)
	point := c.Pos.MultPoint(pointLocal)
	if dist := point.Minus(ray.Origin).Len(); dist >= BIAS {
		return NewHit(c, dist)
	}
	return Miss
}

func (c *Cube) Center() Vector3 {
	return c.Pos.MultPoint(Vector3{})
}

// At returns the normal Vector3 at this point on the Surface
func (c *Cube) At(p Vector3) (normal Direction, mat *Material) {
	i := c.Pos.Inverse() // global to local transform
	p1 := i.MultPoint(p) // translate point into local space
	abs := p1.Abs()
	switch {
	case abs.X > abs.Y && abs.X > abs.Z:
		normal = Direction{math.Copysign(1, p1.X), 0, 0}
	case abs.Y > abs.Z:
		normal = Direction{0, math.Copysign(1, p1.Y), 0}
	default:
		normal = Direction{0, 0, math.Copysign(1, p1.Z)}
	}
	// translate normal from local to global space
	mat = c.Mat
	if c.GridMat != nil && c.GridSize > 0 {
		x, z := p.X/c.GridSize, p.Z/c.GridSize
		if dx := math.Abs(x - math.Floor(x)); dx < 0.08 { // TODO: this should not be a magic number
			mat = c.GridMat
		} else if dz := math.Abs(z - math.Floor(z)); dz < 0.08 {
			mat = c.GridMat
		}
	}
	return c.Pos.MultDir(normal), mat
}

func (c *Cube) Box() *Box {
	return c.box
}
