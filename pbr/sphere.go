package pbr

import "math"

// Sphere describes a 3d sphere
type Sphere struct {
	Center Vector3
	Radius float64
	Mat    Material
}

// Intersect tests whether the sphere intersects a given ray
func (s *Sphere) Intersect(r Ray3) (hit bool, dist float64) {
	op := s.Center.Minus(r.Origin)
	b := op.Dot(r.Dir)
	det := b*b - op.Dot(op) + s.Radius*s.Radius
	if det < 0 {
		return false, 0
	}
	root := math.Sqrt(det)
	t1 := b - root
	if t1 > BIAS {
		return true, t1
	}
	t2 := b + root
	if t2 > BIAS {
		return true, t2
	}
	return false, 0
}

// NormalAt returns the surface normal given a point on the surface
func (s *Sphere) NormalAt(v Vector3) Vector3 {
	return v.Minus(s.Center).Normalize()
}

// MaterialAt returns the material at a given point on the surface
func (s *Sphere) MaterialAt(v Vector3) Material {
	return s.Mat
}
