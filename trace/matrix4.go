package trace

var yAxis Vector3

func init() {
	yAxis = Vector3{0, 1, 0}
}

// Matrix4 handles matrix data and operations
// Column-major (as in math and Direct3D)
// https://fgiesen.wordpress.com/2012/02/12/row-major-vs-column-major-row-vectors-vs-column-vectors/
type Matrix4 struct {
	el [4][4]float64
}

// NewMatrix4 constructs a new matrix
func NewMatrix4(a1, a2, a3, a4, b1, b2, b3, b4, c1, c2, c3, c4, d1, d2, d3, d4 float64) (m Matrix4) {
	m.el = [4][4]float64{
		[4]float64{a1, b1, c1, d1},
		[4]float64{a2, b2, c2, d2},
		[4]float64{a3, b3, c3, d3},
		[4]float64{a4, b4, c4, d4},
	}
	return
}

// Identity creates a new identity matrix
func Identity() Matrix4 {
	return NewMatrix4(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
}

// LookMatrix creates a matrix looking from `from` towards `to`
// http://www.cs.virginia.edu/~gfx/courses/1999/intro.fall99.html/lookat.html
// https://www.3dgep.com/understanding-the-view-matrix/#Look_At_Camera
// http://www.codinglabs.net/article_world_view_projection_matrix.aspx
// https://fgiesen.wordpress.com/2012/02/12/row-major-vs-column-major-row-vectors-vs-column-vectors/
func LookMatrix(o Vector3, to Vector3) Matrix4 {
	f := o.Minus(to).Normalize()    // forward
	r := yAxis.Cross(f).Normalize() // right
	u := f.Cross(r).Normalize()     // up

	return NewMatrix4(
		r.X, u.X, f.X, 0,
		r.Y, u.Y, f.Y, 0,
		r.Z, u.Z, f.Z, 0,
		0, 0, 0, 1,
	)
}

// Translate creates a new translation matrix
func Translate(x, y, z float64) Matrix4 {
	return NewMatrix4(
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	)
}

// Scale creates a new scaling matrix
func Scale(x, y, z float64) Matrix4 {
	return NewMatrix4(
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	)
}

// Trans is a chaining translation
func (a Matrix4) Trans(x, y, z float64) Matrix4 {
	return a.Mult(Translate(x, y, z))
}

// Scale is a chaining scale
func (a Matrix4) Scale(x, y, z float64) Matrix4 {
	return a.Mult(Scale(x, y, z))
}

// Mult multiplies by another matrix4
func (a Matrix4) Mult(b Matrix4) (result Matrix4) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				result.el[j][i] += a.el[k][i] * b.el[j][k]
			}
		}
	}
	return
}

// Point multiplies this matrix4 by a vector, including translation
func (a Matrix4) Point(v Vector3) (result Vector3) {
	result.X = v.X*a.el[0][0] + v.Y*a.el[1][0] + v.Z*a.el[2][0] + a.el[3][0]
	result.Y = v.X*a.el[0][1] + v.Y*a.el[1][1] + v.Z*a.el[2][1] + a.el[3][1]
	result.Z = v.X*a.el[0][2] + v.Y*a.el[1][2] + v.Z*a.el[2][2] + a.el[3][2]
	// final row assumed to be [0,0,0,1]
	return
}

// Dir multiplies this matrix4 by a vector, excluding translation
func (a Matrix4) Dir(v Vector3) (result Vector3) {
	result.X = v.X*a.el[0][0] + v.Y*a.el[1][0] + v.Z*a.el[2][0]
	result.Y = v.X*a.el[0][1] + v.Y*a.el[1][1] + v.Z*a.el[2][1]
	result.Z = v.X*a.el[0][2] + v.Y*a.el[1][2] + v.Z*a.el[2][2]
	return
}

// Ray multiplies this matrix by a ray
// https://gamedev.stackexchange.com/questions/72440/the-correct-way-to-transform-a-ray-with-a-matrix
func (a Matrix4) Ray(r Ray3) (result Ray3) {
	result.Origin = a.Point(r.Origin)
	result.Dir = a.Dir(r.Dir).Normalize()
	return
}

// Inverse returns the inverse of this matrix
func (a Matrix4) Inverse() Matrix4 {
	// TODO: invert a 4x4 matrix
	m := Identity()
	return m
}
