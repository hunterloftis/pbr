package pbr

// Mesh describes a triangle mesh
type Mesh struct {
	Pos  *Matrix4
	Mat  *Material
	Tris []Triangle
}

// Intersect returns whether the ray intersects and where
// TODO: implement next
func (m *Mesh) Intersect(ray Ray3) (bool, float64, int) {
	return false, 0, 0
}

// At returns the material at a point on the mesh
// TODO: implement after Intersect
func (m *Mesh) At(v Vector3, id int) (normal Direction, material *Material) {
	return Vector3{0, 1, 0}.Unit(), m.Mat
}
