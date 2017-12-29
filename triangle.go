package pbr

// Triangle describes a triangle
type Triangle struct {
	Points [3]Vector3
	Normal Direction
	edge1  Vector3
	edge2  Vector3
}

// NewTriangle creates a new triangle
func NewTriangle(a, b, c Vector3) Triangle {
	edge1 := b.Minus(a)
	edge2 := c.Minus(a)
	return Triangle{
		Points: [3]Vector3{a, b, c},
		Normal: edge2.Cross(edge1).Unit(),
		edge1:  edge1,
		edge2:  edge2,
	}
}

// Intersect determines whether or not, and where, a Ray intersects this Triangle
// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (t *Triangle) Intersect(ray Ray3) (bool, float64) {
	const EPS float64 = 0.000001
	h := ray.Dir.Cross(Direction(t.edge2))
	a := t.edge1.Dot(Vector3(h))
	if a > -EPS && a < EPS {
		return false, 0
	}
	f := 1.0 / a
	s := ray.Origin.Minus(t.Points[0])
	u := f * s.Dot(Vector3(h))
	if u < 0 || u > 1 {
		return false, 0
	}
	q := s.Cross(t.edge1)
	v := f * Vector3(ray.Dir).Dot(q)
	if v < 0 || (u+v) > 1 {
		return false, 0
	}
	dist := f * t.edge2.Dot(q)
	if dist < EPS {
		return false, 0
	}
	return true, dist
}
