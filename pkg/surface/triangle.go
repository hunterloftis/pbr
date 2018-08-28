package surface

import (
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Triangle describes a triangle
type Triangle struct {
	Points  [3]geom.Vec // TODO: private fields?
	Normals [3]geom.Dir
	Texture [3]geom.Vec
	Mat     Material
	edge1   geom.Vec
	edge2   geom.Vec
	bounds  *geom.Bounds
}

// NewTriangle creates a new triangle
func NewTriangle(a, b, c geom.Vec, m ...Material) *Triangle {
	edge1 := b.Minus(a)
	edge2 := c.Minus(a)
	n, _ := edge1.Cross(edge2).Unit()
	t := &Triangle{
		Points:  [3]geom.Vec{a, b, c},
		Normals: [3]geom.Dir{n, n, n},
		Mat:     &DefaultMaterial{},
		edge1:   edge1,
		edge2:   edge2,
	}
	if len(m) > 0 {
		t.Mat = m[0]
	}
	min := t.Points[0].Min(t.Points[1]).Min(t.Points[2])
	max := t.Points[0].Max(t.Points[1]).Max(t.Points[2])
	t.bounds = geom.NewBounds(min, max)
	return t
}

func (t *Triangle) Transformed(mtx *geom.Mtx) *Triangle {
	t2 := &Triangle{
		Mat: t.Mat,
	}
	for i := 0; i < 3; i++ {
		t2.Points[i] = mtx.MultPoint(t.Points[i])
		t2.Normals[i] = mtx.MultDir(t.Normals[i])
		t2.Texture[i] = t.Texture[i]
	}
	t2.edge1 = t2.Points[1].Minus(t2.Points[0])
	t2.edge2 = t2.Points[2].Minus(t2.Points[0])
	min := t2.Points[0].Min(t2.Points[1]).Min(t2.Points[2])
	max := t2.Points[0].Max(t2.Points[1]).Max(t2.Points[2])
	t2.bounds = geom.NewBounds(min, max)
	return t2
}

func (t *Triangle) Bounds() *geom.Bounds {
	return t.bounds
}

// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (t *Triangle) Intersect(ray *geom.Ray, max float64) (obj render.Object, dist float64) {
	if ok, near, _ := t.bounds.Check(ray); !ok || near >= max {
		return nil, 0
	}
	h := geom.Vec(ray.Dir).Cross(t.edge2)
	a := t.edge1.Dot(h)
	if a > -bias && a < bias {
		return nil, 0
	}
	f := 1 / a
	s := ray.Origin.Minus(t.Points[0])
	u := f * s.Dot(h)
	if u < 0 || u > 1 {
		return nil, 0
	}
	q := s.Cross(t.edge1)
	v := f * geom.Vec(ray.Dir).Dot(q)
	if v < 0 || u+v > 1 {
		return nil, 0
	}
	dist = f * t.edge2.Dot(q)
	if dist <= bias || dist >= max {
		return nil, 0
	}
	return t, dist
}

// At returns the material at a point on the Triangle
// https://stackoverflow.com/questions/21210774/normal-mapping-on-procedural-sphere
func (t *Triangle) At(pt geom.Vec, in geom.Dir, rnd *rand.Rand) (geom.Dir, render.BSDF) {
	u, v, w := t.Bary(pt)
	n := t.normal(u, v, w)
	texture := t.texture(u, v, w)
	n2, bsdf := t.Mat.At(texture.X, texture.Y, in, n, rnd)
	// TODO: compute binormal and combine texture normal with n to return actual normal
	_ = n2
	normal := n
	return normal, bsdf
}

func (t *Triangle) Lights() []render.Object {
	if !t.Mat.Light().Zero() {
		return []render.Object{t}
	}
	return nil
}

func (t *Triangle) Light() rgb.Energy {
	return t.Mat.Light()
}

func (t *Triangle) Transmit() rgb.Energy {
	return t.Mat.Transmit()
}

// SetNormals sets values for each vertex normal
func (t *Triangle) SetNormals(a, b, c geom.Dir) {
	t.Normals[0] = a
	t.Normals[1] = b
	t.Normals[2] = c
}

func (t *Triangle) SetTexture(a, b, c geom.Vec) {
	t.Texture[0] = a
	t.Texture[1] = b
	t.Texture[2] = c
}

// Normal computes the smoothed normal
func (t *Triangle) normal(u, v, w float64) geom.Dir { // TODO: instead of separate u, v, w just use a Vec and multiply
	n0 := t.Normals[0].Scaled(u)
	n1 := t.Normals[1].Scaled(v)
	n2 := t.Normals[2].Scaled(w)
	n, _ := n0.Plus(n1).Plus(n2).Unit()
	return n
}

func (t *Triangle) texture(u, v, w float64) geom.Vec {
	tex0 := t.Texture[0].Scaled(u)
	tex1 := t.Texture[1].Scaled(v)
	tex2 := t.Texture[2].Scaled(w)
	return tex0.Plus(tex1).Plus(tex2)
}

// Bary returns the Barycentric coords of Vector p on Triangle t
// https://codeplea.com/triangular-interpolation
func (t *Triangle) Bary(p geom.Vec) (u, v, w float64) {
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
