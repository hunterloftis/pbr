package surface

import (
	"math"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/material"
)

// Sphere describes a 3d sphere
// TODO: make all of these private, this is accessed through interfaces anyway
type Sphere struct {
	Pos *geom.Matrix4
	Mat material.Description
	box *Box
}

// UnitSphere returns a pointer to a new 1x1x1 Sphere Surface with a given material and optional transforms.
func UnitSphere(m ...material.Description) *Sphere {
	s := &Sphere{
		Pos: geom.Identity(),
		Mat: material.Default,
	}
	if len(m) > 0 {
		s.Mat = m[0]
	}
	return s.transform(geom.Identity())
}

// TODO: unify with cube.transform AABB calc
func (s *Sphere) transform(t *geom.Matrix4) *Sphere {
	s.Pos = s.Pos.Mult(t)
	min := s.Pos.MultPoint(geom.Vector3{})
	max := s.Pos.MultPoint(geom.Vector3{})
	for x := -0.5; x <= 0.5; x += 1 {
		for y := -0.5; y <= 0.5; y += 1 {
			for z := -0.5; z <= 0.5; z += 1 {
				pt := s.Pos.MultPoint(geom.Vector3{x, y, z})
				min = min.Min(pt)
				max = max.Max(pt)
			}
		}
	}
	s.box = NewBox(min, max)
	return s
}

func (s *Sphere) Move(x, y, z float64) *Sphere {
	return s.transform(geom.Trans(x, y, z))
}

func (s *Sphere) Scale(x, y, z float64) *Sphere {
	return s.transform(geom.Scale(x, y, z))
}

func (s *Sphere) Rotate(x, y, z float64) *Sphere {
	return s.transform(geom.Rot(geom.Vector3{x, y, z}))
}

func (s *Sphere) Center() geom.Vector3 {
	return s.Pos.MultPoint(geom.Vector3{})
}

func (s *Sphere) Box() *Box {
	return s.box
}

// Intersect tests whether the sphere intersects a given ray.
// http://tfpsly.free.fr/english/index.html?url=http://tfpsly.free.fr/english/3d/Raytracing.html
// TODO: http://kylehalladay.com/blog/tutorial/math/2013/12/24/Ray-Sphere-Intersection.html
func (s *Sphere) Intersect(ray *geom.Ray3) Hit {
	if ok, _, _ := s.box.Check(ray); !ok {
		return Miss
	}
	i := s.Pos.Inverse()
	r := i.MultRay(ray)
	op := geom.Vector3{}.Minus(r.Origin)
	b := op.Dot(geom.Vector3(r.Dir))
	det := b*b - op.Dot(op) + 0.5*0.5
	if det < 0 {
		return Miss
	}
	root := math.Sqrt(det)
	t1 := b - root
	if t1 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t1)).Len()
		if dist > bias {
			return NewHit(s, dist)
		}
	}
	t2 := b + root
	if t2 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t2)).Len()
		if dist > bias {
			return NewHit(s, dist)
		}
	}
	return Miss
}

// At returns the surface normal given a point on the surface.
func (s *Sphere) At(pt geom.Vector3) (normal geom.Direction, material *material.Sample) {
	i := s.Pos.Inverse()
	p := i.MultPoint(pt)
	return s.Pos.MultDir(p.Unit()), s.Mat.At(0, 0)
}

func (s *Sphere) Emits() bool {
	return s.Mat.Emits()
}
