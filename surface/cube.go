package surface

import (
	"math"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/surface/material"
)

// Cube describes the orientation and material of a unit cube
type Cube struct {
	Pos      *geom.Matrix4
	Mat      *material.Material
	GridMat  *material.Material
	GridSize float64
	box      *Box
}

// UnitCube returns a pointer to a new 1x1x1 Cube Surface with material and optional transforms.
func UnitCube(m ...*material.Material) *Cube {
	c := &Cube{
		Pos: geom.Identity(),
		Mat: material.Default,
	}
	if len(m) > 0 {
		c.Mat = m[0]
	}
	return c.transform(geom.Identity())
}

func (c *Cube) transform(m *geom.Matrix4) *Cube {
	c.Pos = c.Pos.Mult(m)
	min := c.Pos.MultPoint(geom.Vector3{})
	max := c.Pos.MultPoint(geom.Vector3{})
	for x := -1.0; x <= 1; x += 2 {
		for y := -1.0; y <= 1; y += 2 {
			for z := -1.0; z <= 1; z += 2 {
				pt := c.Pos.MultPoint(geom.Vector3{x, y, z})
				min = min.Min(pt)
				max = max.Max(pt)
			}
		}
	}
	c.box = NewBox(min, max)
	return c
}

func (c *Cube) Move(x, y, z float64) *Cube {
	return c.transform(geom.Trans(x, y, z))
}

func (c *Cube) Scale(x, y, z float64) *Cube {
	return c.transform(geom.Scale(x, y, z))
}

func (c *Cube) Rotate(x, y, z float64) *Cube {
	return c.transform(geom.Rot(geom.Vector3{x, y, z}))
}

// SetGrid adds a second material to the cube which is applied as a grid across its surface
func (c *Cube) SetGrid(mat *material.Material, size float64) *Cube {
	c.GridMat = mat
	c.GridSize = size
	return c
}

func (c *Cube) Intersect(ray *geom.Ray3) Hit {
	if ok, _, _ := c.box.Check(ray); !ok {
		return Miss
	}
	inv := c.Pos.Inverse() // global to local transform
	r := inv.MultRay(ray)  // translate ray into local space
	tmin := 0.0
	tmax := math.Inf(1)
	for a := 0; a < 3; a++ {
		t0 := (-0.5 - r.OrArray[a]) * r.InvArray[a]
		t1 := (0.5 - r.OrArray[a]) * r.InvArray[a]
		if r.InvArray[a] < 0 {
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
	if tmin > 0 {
		if dist := c.Pos.MultDist(r.Dir.Scaled(tmin)).Len(); dist >= bias {
			return NewHit(c, dist)
		}
	}
	if tmax > 0 {
		if dist := c.Pos.MultDist(r.Dir.Scaled(tmax)).Len(); dist >= bias {
			return NewHit(c, dist)
		}
	}
	return Miss
}

func (c *Cube) Center() geom.Vector3 {
	return c.Pos.MultPoint(geom.Vector3{})
}

func (c *Cube) Material() *material.Material {
	return c.Mat
}

// At returns the normal geom.Vector3 at this point on the Surface
func (c *Cube) At(p geom.Vector3) (normal geom.Direction, mat *material.Material) {
	i := c.Pos.Inverse() // global to local transform
	p1 := i.MultPoint(p) // translate point into local space
	abs := p1.Abs()
	switch {
	case abs.X > abs.Y && abs.X > abs.Z:
		normal = geom.Direction{math.Copysign(1, p1.X), 0, 0}
	case abs.Y > abs.Z:
		normal = geom.Direction{0, math.Copysign(1, p1.Y), 0}
	default:
		normal = geom.Direction{0, 0, math.Copysign(1, p1.Z)}
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
