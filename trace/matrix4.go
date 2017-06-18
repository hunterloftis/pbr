package trace

// Matrix4 handles matrix data and operations
// Column-major (as in math and Direct3D)
// https://fgiesen.wordpress.com/2012/02/12/row-major-vs-column-major-row-vectors-vs-column-vectors/
type Matrix4 struct {
	el [4][4]float64
}

// NewMatrix4 constructs a new matrix
func NewMatrix4(a1, a2, a3, a4, b1, b2, b3, b4, c1, c2, c3, c4, d1, d2, d3, d4 float64) (m Matrix4) {
	m.el[0] = [4]float64{a1, b1, c1, d1}
	m.el[1] = [4]float64{a2, b2, c2, d2}
	m.el[2] = [4]float64{a3, b3, c3, d3}
	m.el[3] = [4]float64{a4, b4, c4, d4}
	return
}

// NewIDMatrix4 constructs a new identity matrix
func NewIDMatrix4() (m Matrix4) {
	m.el[0] = [4]float64{1, 0, 0, 0}
	m.el[1] = [4]float64{0, 1, 0, 0}
	m.el[2] = [4]float64{0, 0, 1, 0}
	m.el[3] = [4]float64{0, 0, 0, 1}
	return
}

// NewLookMatrix4 creates a matrix looking from `from` towards `to`
// http://www.cs.virginia.edu/~gfx/courses/1999/intro.fall99.html/lookat.html
// https://www.3dgep.com/understanding-the-view-matrix/#Look_At_Camera
// http://www.codinglabs.net/article_world_view_projection_matrix.aspx
// https://fgiesen.wordpress.com/2012/02/12/row-major-vs-column-major-row-vectors-vs-column-vectors/
func NewLookMatrix4(o Vector3, to Vector3) Matrix4 {
	f := to.Minus(o).Normalize() // forward
	r := yAxis.Cross(f)          // right
	u := f.Cross(r)              // up

	return NewMatrix4(
		r.X, u.X, f.X, o.X,
		r.Y, u.Y, f.Y, o.Y,
		r.Z, u.Z, f.Z, o.Z,
		0, 0, 0, 1,
	)
}

// Mult multiplies by another matrix4
func (a *Matrix4) Mult(b Matrix4) (result Matrix4) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				result.el[i][j] += a.el[i][k] * b.el[k][j]
			}
		}
	}
	return
}

// ApplyPoint multiplies this matrix4 by a vector, including translation
func (a *Matrix4) ApplyPoint(v Vector3) (result Vector3) {
	result.X = v.X*a.el[0][0] + v.Y*a.el[0][1] + v.Z*a.el[0][2] + a.el[0][3]
	result.Y = v.X*a.el[1][0] + v.Y*a.el[1][1] + v.Z*a.el[1][2] + a.el[1][3]
	result.Z = v.X*a.el[2][0] + v.Y*a.el[2][1] + v.Z*a.el[2][2] + a.el[2][3]
	// final row assumed to be [0,0,0,1]
	return
}

// ApplyDir multiplies this matrix4 by a vector, excluding translation
func (a *Matrix4) ApplyDir(v Vector3) (result Vector3) {
	result.X = v.X*a.el[0][0] + v.Y*a.el[0][1] + v.Z*a.el[0][2]
	result.Y = v.X*a.el[1][0] + v.Y*a.el[1][1] + v.Z*a.el[1][2]
	result.Z = v.X*a.el[2][0] + v.Y*a.el[2][1] + v.Z*a.el[2][2]
	return
}
