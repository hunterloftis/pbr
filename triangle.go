package pbr

// Triangle describes a triangle
// TODO: store per-vertex Normal data so .obj file curved surfaces can be read in and rendered smoothly / without edges
type Triangle struct {
	Points  [3]Vector3
	Normals [3]Direction
	Mat     *Material
	edge1   Vector3
	edge2   Vector3
}

// NewTriangle creates a new triangle
func NewTriangle(a, b, c Vector3, m *Material) Triangle {
	edge1 := b.Minus(a)
	edge2 := c.Minus(a)
	n := edge1.Cross(edge2).Unit()
	return Triangle{
		Points:  [3]Vector3{a, b, c},
		Normals: [3]Direction{n, n, n},
		Mat:     m,
		edge1:   edge1,
		edge2:   edge2,
	}
}

// Intersect determines whether or not, and where, a Ray intersects this Triangle
// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (t *Triangle) Intersect(ray Ray3) (bool, float64) {
	h := ray.Dir.Cross(Direction(t.edge2))
	a := t.edge1.Dot(Vector3(h))
	if a > -BIAS && a < BIAS {
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
	if dist < BIAS {
		return false, 0
	}
	return true, dist
}

// At returns the material at a point on the Triangle
func (t *Triangle) At(v Vector3) (norm Direction, mat *Material) {
	return t.Normal(v), t.Mat
}

// SetNormals sets values for each vertex normal
// TODO: figure out proper use of pointers
func (t *Triangle) SetNormals(a, b, c *Direction) {
	if a != nil {
		t.Normals[0] = *a
	}
	if b != nil {
		t.Normals[1] = *b
	}
	if c != nil {
		t.Normals[2] = *c
	}
}

// Normal computes the smoothed normal
func (t *Triangle) Normal(p Vector3) Direction {
	u, v, w := t.Bary(p)
	n0 := t.Normals[0].Scaled(u)
	n1 := t.Normals[1].Scaled(v)
	n2 := t.Normals[2].Scaled(w)
	return n0.Plus(n1).Plus(n2).Unit()
}

// Bary returns the Barycentric coords of Vector p on Triangle t
// TODO: using this in several places; integrate
// https://codeplea.com/triangular-interpolation
func (t *Triangle) Bary(p Vector3) (u, v, w float64) {
	v0 := t.Points[1].Minus(t.Points[0])
	v1 := t.Points[2].Minus(t.Points[0])
	v2 := p.Minus(t.Points[0])
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	d := d00*d11 - d01*d01
	v = (d11*d20 - d01*d21) / d
	w = (d00*d21 - d01*d20) / d
	u = 1 - v - w
	return
}
