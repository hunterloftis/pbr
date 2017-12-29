package pbr

import (
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
	nearest := math.MaxFloat64
	id := -1
	for i, t := range m.Tris {
		hit, dist := t.Intersect(ray)
		if hit && dist < nearest {
			id = i
			nearest = dist
		}
	}
	if id == -1 {
		return false, 0, 0
	}
	return true, nearest, id
}

// At returns the material at a point on the mesh
// TODO: implement after Intersect
func (m *Mesh) At(v Vector3, id int) (normal Direction, material *Material) {
	t := m.Tris[id]
	return t.Normal, m.Mat
}
