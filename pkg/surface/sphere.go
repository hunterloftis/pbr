package surface

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

// Sphere describes a 3d sphere
type Sphere struct {
	mtx    *geom.Mtx
	mat    Material
	bounds *geom.Bounds
}

// UnitSphere returns a pointer to a new 1x1x1 Sphere Surface with a given material and optional transforms.
func UnitSphere(m ...Material) *Sphere {
	s := &Sphere{
		mtx: geom.Identity(),
		mat: &DefaultMaterial{},
	}
	if len(m) > 0 {
		s.mat = m[0]
	}
	return s.transform(geom.Identity())
}

// TODO: unify with cube.transform AABB calc
func (s *Sphere) transform(t *geom.Mtx) *Sphere {
	s.mtx = s.mtx.Mult(t)
	min := s.mtx.MultPoint(geom.Vec{})
	max := s.mtx.MultPoint(geom.Vec{})
	for x := -0.5; x <= 0.5; x += 1 {
		for y := -0.5; y <= 0.5; y += 1 {
			for z := -0.5; z <= 0.5; z += 1 {
				pt := s.mtx.MultPoint(geom.Vec{x, y, z})
				min = min.Min(pt)
				max = max.Max(pt)
			}
		}
	}
	s.bounds = geom.NewBounds(min, max)
	return s
}

func (s *Sphere) Shift(v geom.Vec) *Sphere {
	return s.transform(geom.Shift(v))
}

func (s *Sphere) Scale(v geom.Vec) *Sphere {
	return s.transform(geom.Scale(v))
}

func (s *Sphere) Rotate(v geom.Vec) *Sphere {
	return s.transform(geom.Rotate(v))
}

func (s *Sphere) Center() geom.Vec {
	return s.mtx.MultPoint(geom.Vec{})
}

func (s *Sphere) Bounds() *geom.Bounds {
	return s.bounds
}

// Intersect tests whether the sphere intersects a given ray.
// http://tfpsly.free.fr/english/index.html?url=http://tfpsly.free.fr/english/3d/Raytracing.html
// TODO: http://kylehalladay.com/blog/tutorial/math/2013/12/24/Ray-Sphere-Intersection.html
func (s *Sphere) Intersect(ray *geom.Ray, max float64) (obj render.Object, dist float64) {
	if ok, near, _ := s.bounds.Check(ray); !ok || near >= max {
		return nil, 0
	}
	i := s.mtx.Inverse()
	r := i.MultRay(ray)
	op := geom.Vec{}.Minus(r.Origin)
	b := op.Dot(geom.Vec(r.Dir))
	det := b*b - op.Dot(op) + 0.5*0.5
	if det < 0 {
		return nil, 0
	}
	root := math.Sqrt(det)
	t1 := b - root
	if t1 > 0 {
		dist := s.mtx.MultDist(r.Dir.Scaled(t1)).Len()
		if dist > bias {
			return s, dist
		}
	}
	t2 := b + root
	if t2 > 0 {
		dist := s.mtx.MultDist(r.Dir.Scaled(t2)).Len()
		if dist > bias {
			return s, dist
		}
	}
	return nil, 0
}

// At returns the surface normal given a point on the surface.
func (s *Sphere) At(pt geom.Vec, in geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf render.BSDF) {
	i := s.mtx.Inverse()
	p := i.MultPoint(pt)
	pu, _ := p.Unit()
	n := s.mtx.MultDir(pu)
	n2, bsdf := s.mat.At(0, 0, in, n, rnd)
	_ = n2
	normal = n // TODO: compute normal by combining n and n2 (and a bitangent)
	return normal, bsdf
}

func (s *Sphere) Light() rgb.Energy {
	return s.mat.Light()
}

func (s *Sphere) Transmit() rgb.Energy {
	return s.mat.Transmit()
}

func (s *Sphere) Lights() []render.Object {
	if !s.mat.Light().Zero() {
		return []render.Object{s}
	}
	return nil
}
