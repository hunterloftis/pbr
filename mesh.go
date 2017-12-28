package pbr

import (
	"fmt"
	"math"
)

// Mesh describes a triangle mesh
type Mesh struct {
	Pos  *Matrix4
	Mat  *Material
	Tris []Triangle
}

// Intersect returns whether the ray intersects and where
// TODO: implement next
// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (m *Mesh) Intersect(ray Ray3) (bool, float64, int) {
	const EPS float64 = 0.000001
	nearest := math.MaxFloat64
	id := -1
	for _, t := range m.Tris {
		edge1 := t.Points[1].Minus(t.Points[0])
		edge2 := t.Points[2].Minus(t.Points[0])
		h := ray.Dir.Cross(Direction(edge2))
		a := edge1.Dot(Vector3(h))
		if a > -EPS && a < EPS {
			continue
		}
		f := 1.0 / a
		s := ray.Origin.Minus(t.Points[0])
		u := f * s.Dot(Vector3(h))
		if u < 0 || u > 1 {
			continue
		}
		q := s.Cross(edge1)
		v := f * Vector3(ray.Dir).Dot(q)
		if v < 0 || (u+v) > 1 {
			continue
		}
		dist := f * edge2.Dot(q)
		if dist < EPS {
			continue
		}
		if dist < nearest {
			nearest = dist
		}
	}
	if id == -1 {
		return false, 0, 0
	}
	fmt.Println("hit at", nearest)
	panic("woah")
	return true, nearest, id
}

// At returns the material at a point on the mesh
// TODO: implement after Intersect
func (m *Mesh) At(v Vector3, id int) (normal Direction, material *Material) {
	return Vector3{0, 1, 0}.Unit(), m.Mat
}
