package surface

import (
	"math"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/surface/material"
)

// Sphere describes a 3d sphere
// TODO: make all of these private, this is accessed through interfaces anyway
type Sphere struct {
	Pos *geom.Matrix4
	Mat *material.Material
	box *Box
}

// UnitSphere returns a pointer to a new 1x1x1 Sphere Surface with a given material and optional transforms.
func UnitSphere(m ...*material.Material) *Sphere {
	s := &Sphere{
		Pos: geom.Identity(),
		Mat: material.Default,
	}
	if len(m) > 0 {
		s.Mat = m[0]
	}
	return s.transform(geom.Identity())
}

// TODO: https://tavianator.com/exact-bounding-boxes-for-spheres-ellipsoids/
func (s *Sphere) transform(m *geom.Matrix4) *Sphere {
	s.Pos = s.Pos.Mult(m)
	min := s.Pos.MultPoint(geom.Vector3{-1, -1, -1})
	max := s.Pos.MultPoint(geom.Vector3{1, 1, 1})
	s.box = NewBox(min, max)
	return s
}

func (s *Sphere) Material() *material.Material {
	return s.Mat
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
func (s *Sphere) At(point geom.Vector3) (geom.Direction, *material.Material) {
	i := s.Pos.Inverse()
	p := i.MultPoint(point)
	return s.Pos.MultDir(p.Unit()), s.Mat
}
