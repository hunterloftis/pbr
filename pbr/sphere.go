package pbr

import "math"

// Sphere describes a 3d sphere
type Sphere struct {
	Pos *Matrix4
	Mat *Material
}

// UnitSphere returns a pointer to a new 1x1x1 Sphere Surface with position pos and material mat.
func UnitSphere(m *Material, transforms ...*Matrix4) *Sphere {
	pos := Identity()
	for _, t := range transforms {
		pos = pos.Mult(t)
	}
	return &Sphere{
		Pos: pos,
		Mat: m,
	}
}

// Intersect tests whether the sphere intersects a given ray
// http://tfpsly.free.fr/english/index.html?url=http://tfpsly.free.fr/english/3d/Raytracing.html
func (s *Sphere) Intersect(ray Ray3) (hit bool, dist float64) {
	i := s.Pos.Inverse()
	r := i.MultRay(ray)
	op := Vector3{}.Minus(r.Origin)
	b := op.Dot(Vector3(r.Dir))
	det := b*b - op.Dot(op) + 0.5*0.5
	if det < 0 {
		return false, 0
	}
	root := math.Sqrt(det)
	t1 := b - root
	if t1 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t1)).Len()
		if dist > Bias {
			return true, dist
		}
	}
	t2 := b + root
	if t2 > 0 {
		dist := s.Pos.MultDist(r.Dir.Scaled(t2)).Len()
		if dist > Bias {
			return true, dist
		}
	}
	return false, 0
}

// At returns the surface normal given a point on the surface
func (s *Sphere) At(point Vector3) (Direction, *Material) {
	i := s.Pos.Inverse()
	p := i.MultPoint(point)
	return s.Pos.MultDir(p.Unit()), s.Mat
}
