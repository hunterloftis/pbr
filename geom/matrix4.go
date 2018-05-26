package geom

import (
	"math"
)

var yAxis = Direction{0, 1, 0}

// Matrix4 handles matrix data and operations
// Column-major (as in math and Direct3D)
// https://fgiesen.wordpress.com/2012/02/12/row-major-vs-column-major-row-vectors-vs-column-vectors/
// TODO: Rename to Matrix?
type Matrix4 struct {
	el  [4][4]float64
	inv *Matrix4
}

// NewMatrix4 constructs a new matrix
func NewMatrix4(a1, a2, a3, a4, b1, b2, b3, b4, c1, c2, c3, c4, d1, d2, d3, d4 float64) *Matrix4 {
	m := Matrix4{
		el: [4][4]float64{
			[4]float64{a1, b1, c1, d1},
			[4]float64{a2, b2, c2, d2},
			[4]float64{a3, b3, c3, d3},
			[4]float64{a4, b4, c4, d4},
		},
	}
	return &m
}

// Identity creates a new identity matrix
func Identity() *Matrix4 {
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
func LookMatrix(o Vector3, to Vector3) *Matrix4 {
	f, _ := o.Minus(to).Unit() // forward
	r, _ := yAxis.Cross(f)     // right
	u, _ := f.Cross(r)         // up
	orient := NewMatrix4(
		r.X, u.X, f.X, 0,
		r.Y, u.Y, f.Y, 0,
		r.Z, u.Z, f.Z, 0,
		0, 0, 0, 1,
	)
	return Trans(o.X, o.Y, o.Z).Mult(orient)
}

// Trans creates a new translation matrix
func Trans(x, y, z float64) *Matrix4 {
	return NewMatrix4(
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	)
}

// Scale creates a new scaling matrix
func Scale(x, y, z float64) *Matrix4 {
	return NewMatrix4(
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	)
}

// Rot creates a rotation matrix from an angle-axis Vector representation
// http://www.euclideanspace.com/maths/geometry/rotations/conversions/angleToMatrix/
func Rot(v Vector3) *Matrix4 {
	a := v.Len()
	c := math.Cos(a)
	s := math.Sin(a)
	t := 1 - c
	n, _ := v.Unit()
	x, y, z := n.X, n.Y, n.Z
	return NewMatrix4(
		t*x*x+c, t*x*y-z*s, t*x*z+y*s, 0,
		t*x*y+z*s, t*y*y+c, t*y*z-x*s, 0,
		t*x*z-y*s, t*y*z+x*s, t*z*z+c, 0,
		0, 0, 0, 1,
	)
}

// Mult multiplies by another matrix4
// TODO: a and b might be flipped here
func (a *Matrix4) Mult(b *Matrix4) *Matrix4 {
	m := Matrix4{}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				m.el[j][i] += a.el[k][i] * b.el[j][k]
			}
		}
	}
	return &m
}

// Equals tests whether two Matrices have equal values
func (a *Matrix4) Equals(b *Matrix4) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if a.el[i][j] != b.el[i][j] {
				return false
			}
		}
	}
	return true
}

// MultPoint multiplies this matrix4 by a vector, including translation
func (a *Matrix4) MultPoint(v Vector3) (result Vector3) {
	result.X = v.X*a.el[0][0] + v.Y*a.el[1][0] + v.Z*a.el[2][0] + a.el[3][0]
	result.Y = v.X*a.el[0][1] + v.Y*a.el[1][1] + v.Z*a.el[2][1] + a.el[3][1]
	result.Z = v.X*a.el[0][2] + v.Y*a.el[1][2] + v.Z*a.el[2][2] + a.el[3][2]
	// final row assumed to be [0,0,0,1]
	return
}

// MultDist multiplies this matrix4 by a vector, excluding translation
func (a *Matrix4) MultDist(v Vector3) (result Vector3) {
	result.X = v.X*a.el[0][0] + v.Y*a.el[1][0] + v.Z*a.el[2][0]
	result.Y = v.X*a.el[0][1] + v.Y*a.el[1][1] + v.Z*a.el[2][1]
	result.Z = v.X*a.el[0][2] + v.Y*a.el[1][2] + v.Z*a.el[2][2]
	return
}

// MultDir multiplies this matrix4 by a direction, renormalizing the result
func (a *Matrix4) MultDir(v Direction) (result Direction) {
	dir, _ := a.MultDist(Vector3(v)).Unit()
	return dir
}

// MultRay multiplies this matrix by a ray
// https://gamedev.stackexchange.com/questions/72440/the-correct-way-to-transform-a-ray-with-a-matrix
func (a *Matrix4) MultRay(r *Ray3) *Ray3 {
	return NewRay(a.MultPoint(r.Origin), a.MultDir(r.Dir))
}

// Inverse returns the inverse of this matrix
// https://www.gamedev.net/forums/topic/648190-algorithm-for-4x4-matrix-inverse/
// https://stackoverflow.com/questions/1148309/inverting-a-4x4-matrix
func (a *Matrix4) Inverse() *Matrix4 {
	if a.inv != nil {
		return a.inv
	}
	i := Identity()
	e := a.el
	i.el[0][0] = e[1][1]*e[2][2]*e[3][3] - e[1][1]*e[2][3]*e[3][2] - e[2][1]*e[1][2]*e[3][3] + e[2][1]*e[1][3]*e[3][2] + e[3][1]*e[1][2]*e[2][3] - e[3][1]*e[1][3]*e[2][2]
	i.el[1][0] = e[1][0]*e[2][3]*e[3][2] - e[1][0]*e[2][2]*e[3][3] + e[2][0]*e[1][2]*e[3][3] - e[2][0]*e[1][3]*e[3][2] - e[3][0]*e[1][2]*e[2][3] + e[3][0]*e[1][3]*e[2][2]
	i.el[2][0] = e[1][0]*e[2][1]*e[3][3] - e[1][0]*e[2][3]*e[3][1] - e[2][0]*e[1][1]*e[3][3] + e[2][0]*e[1][3]*e[3][1] + e[3][0]*e[1][1]*e[2][3] - e[3][0]*e[1][3]*e[2][1]
	i.el[3][0] = e[1][0]*e[2][2]*e[3][1] - e[1][0]*e[2][1]*e[3][2] + e[2][0]*e[1][1]*e[3][2] - e[2][0]*e[1][2]*e[3][1] - e[3][0]*e[1][1]*e[2][2] + e[3][0]*e[1][2]*e[2][1]
	i.el[0][1] = e[0][1]*e[2][3]*e[3][2] - e[0][1]*e[2][2]*e[3][3] + e[2][1]*e[0][2]*e[3][3] - e[2][1]*e[0][3]*e[3][2] - e[3][1]*e[0][2]*e[2][3] + e[3][1]*e[0][3]*e[2][2]
	i.el[1][1] = e[0][0]*e[2][2]*e[3][3] - e[0][0]*e[2][3]*e[3][2] - e[2][0]*e[0][2]*e[3][3] + e[2][0]*e[0][3]*e[3][2] + e[3][0]*e[0][2]*e[2][3] - e[3][0]*e[0][3]*e[2][2]
	i.el[2][1] = e[0][0]*e[2][3]*e[3][1] - e[0][0]*e[2][1]*e[3][3] + e[2][0]*e[0][1]*e[3][3] - e[2][0]*e[0][3]*e[3][1] - e[3][0]*e[0][1]*e[2][3] + e[3][0]*e[0][3]*e[2][1]
	i.el[3][1] = e[0][0]*e[2][1]*e[3][2] - e[0][0]*e[2][2]*e[3][1] - e[2][0]*e[0][1]*e[3][2] + e[2][0]*e[0][2]*e[3][1] + e[3][0]*e[0][1]*e[2][2] - e[3][0]*e[0][2]*e[2][1]
	i.el[0][2] = e[0][1]*e[1][2]*e[3][3] - e[0][1]*e[1][3]*e[3][2] - e[1][1]*e[0][2]*e[3][3] + e[1][1]*e[0][3]*e[3][2] + e[3][1]*e[0][2]*e[1][3] - e[3][1]*e[0][3]*e[1][2]
	i.el[1][2] = e[0][0]*e[1][3]*e[3][2] - e[0][0]*e[1][2]*e[3][3] + e[1][0]*e[0][2]*e[3][3] - e[1][0]*e[0][3]*e[3][2] - e[3][0]*e[0][2]*e[1][3] + e[3][0]*e[0][3]*e[1][2]
	i.el[2][2] = e[0][0]*e[1][1]*e[3][3] - e[0][0]*e[1][3]*e[3][1] - e[1][0]*e[0][1]*e[3][3] + e[1][0]*e[0][3]*e[3][1] + e[3][0]*e[0][1]*e[1][3] - e[3][0]*e[0][3]*e[1][1]
	i.el[3][2] = e[0][0]*e[1][2]*e[3][1] - e[0][0]*e[1][1]*e[3][2] + e[1][0]*e[0][1]*e[3][2] - e[1][0]*e[0][2]*e[3][1] - e[3][0]*e[0][1]*e[1][2] + e[3][0]*e[0][2]*e[1][1]
	i.el[0][3] = e[0][1]*e[1][3]*e[2][2] - e[0][1]*e[1][2]*e[2][3] + e[1][1]*e[0][2]*e[2][3] - e[1][1]*e[0][3]*e[2][2] - e[2][1]*e[0][2]*e[1][3] + e[2][1]*e[0][3]*e[1][2]
	i.el[1][3] = e[0][0]*e[1][2]*e[2][3] - e[0][0]*e[1][3]*e[2][2] - e[1][0]*e[0][2]*e[2][3] + e[1][0]*e[0][3]*e[2][2] + e[2][0]*e[0][2]*e[1][3] - e[2][0]*e[0][3]*e[1][2]
	i.el[2][3] = e[0][0]*e[1][3]*e[2][1] - e[0][0]*e[1][1]*e[2][3] + e[1][0]*e[0][1]*e[2][3] - e[1][0]*e[0][3]*e[2][1] - e[2][0]*e[0][1]*e[1][3] + e[2][0]*e[0][3]*e[1][1]
	i.el[3][3] = e[0][0]*e[1][1]*e[2][2] - e[0][0]*e[1][2]*e[2][1] - e[1][0]*e[0][1]*e[2][2] + e[1][0]*e[0][2]*e[2][1] + e[2][0]*e[0][1]*e[1][2] - e[2][0]*e[0][2]*e[1][1]
	det := 1.0 / (e[0][0]*i.el[0][0] + e[0][1]*i.el[1][0] + e[0][2]*i.el[2][0] + e[0][3]*i.el[3][0])
	for j := 0; j < 4; j++ {
		for k := 0; k < 4; k++ {
			i.el[j][k] *= det
		}
	}
	a.inv, i.inv = i, a
	return i
}

func (a *Matrix4) At(col, row int) float64 {
	return a.el[col-1][row-1]
}

func (a *Matrix4) Transpose() *Matrix4 {
	m := &Matrix4{}
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			m.el[row][col] = a.el[col][row]
		}
	}
	return m
}
