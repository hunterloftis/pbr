package pbr

import (
	"math"
)

// Sphere describes a 3d sphere
type Sphere struct {
	Pos *Matrix4
	Mat *Material
	box *Box
}

// UnitSphere returns a pointer to a new 1x1x1 Sphere Surface with a given material and optional transforms.
func UnitSphere(m *Material, transforms ...*Matrix4) *Sphere {
	pos := Identity()
	for _, t := range transforms {
		pos = pos.Mult(t)
	}
	s := &Sphere{
		Pos: pos,
		Mat: m,
	}
	min := s.Pos.MultPoint(Vector3{-1, -1, -1})
	max := s.Pos.MultPoint(Vector3{1, 1, 1})
	s.box = NewBox(min, max)
	return s
}

func (s *Sphere) Center() Vector3 {
	return s.Pos.MultPoint(Vector3{})
}

func (s *Sphere) Box() *Box {
	return s.box
}

// Intersect tests whether the sphere intersects a given ray.
// http://tfpsly.free.fr/english/index.html?url=http://tfpsly.free.fr/english/3d/Raytracing.html
func (s *Sphere) Intersect(ray *Ray3) Hit {
	if ok, _ := s.box.Check(ray); !ok {
		return Miss
	}
	i := s.Pos.Inverse()
	r := i.MultRay(ray)
	op := Vector3{}.Minus(r.Origin)
	b := op.Dot(Vector3(r.Dir))
	det := b*b - op.Dot(op) + 0.5*0.5
	if det < 0 {
		return Miss
	}
	root := math.Sqrt(det)
	t1 := b - root
	if t1 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t1)).Len()
		if dist > BIAS {
			return NewHit(s, dist)
		}
	}
	t2 := b + root
	if t2 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t2)).Len()
		if dist > BIAS {
			return NewHit(s, dist)
		}
	}
	return Miss
}

// At returns the surface normal given a point on the surface.
func (s *Sphere) At(point Vector3, dir Direction) (Direction, *Material) {
	i := s.Pos.Inverse()
	p := i.MultPoint(point)
	return s.Pos.MultDir(p.Unit()), s.Mat
}
