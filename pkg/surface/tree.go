package surface

import (
	"sort"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

const (
	minContents = 8
	maxDepth    = 16
)

type Tree struct {
	branch
	lights []render.Object
}

type branch struct {
	surfaces []render.Surface
	bounds   *geom.Bounds
	left     *branch
	right    *branch
	axis     int
	wall     float64
	leaf     bool
}

func NewTree(ss ...render.Surface) *Tree {
	t := Tree{
		branch: *newBranch(BoundsAround(ss), ss, maxDepth),
	}
	for _, s := range t.branch.surfaces {
		t.lights = append(t.lights, s.Lights()...)
	}
	return &t
}

func (t *Tree) Lights() []render.Object {
	return t.lights
}

func (t *Tree) Bounds() *geom.Bounds {
	return t.bounds
}

func newBranch(bounds *geom.Bounds, surfaces []render.Surface, depth int) *branch {
	b := branch{
		surfaces: overlaps(bounds, surfaces),
		bounds:   bounds,
	}
	if depth <= 0 || len(b.surfaces) <= minContents {
		b.leaf = true
		return &b
	}
	b.axis = 0
	max := -1.0
	for a := 0; a < 3; a++ {
		dist := bounds.Max.Axis(a) - bounds.Min.Axis(a)
		if dist > max {
			max = dist
			b.axis = a
		}
	}
	b.wall = median(b.surfaces, b.axis)
	lBounds, rBounds := bounds.Split(b.axis, b.wall)
	b.left = newBranch(lBounds, b.surfaces, depth-1)
	b.right = newBranch(rBounds, b.surfaces, depth-1)
	return &b
}

// http://slideplayer.com/slide/7653218/
func (b *branch) Intersect(ray *geom.Ray, maxDist float64) (obj render.Object, dist float64) {
	hit, min, max := b.bounds.Check(ray)
	if !hit || min >= maxDist {
		return nil, 0
	}
	if b.leaf {
		return b.IntersectSurfaces(ray, max)
	}
	var near, far *branch
	if ray.DirArray[b.axis] >= 0 {
		near, far = b.left, b.right
	} else {
		near, far = b.right, b.left
	}
	split := (b.wall - ray.OrArray[b.axis]) * ray.InvArray[b.axis]
	if min >= split {
		return far.Intersect(ray, maxDist)
	}
	if max <= split {
		return near.Intersect(ray, maxDist)
	}
	if o, d := near.Intersect(ray, maxDist); o != nil {
		return o, d
	}
	return far.Intersect(ray, maxDist)
}

func (b *branch) IntersectSurfaces(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range b.surfaces {
		if o, d := s.Intersect(r, dist); o != nil {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func overlaps(bounds *geom.Bounds, surfaces []render.Surface) []render.Surface {
	matches := make([]render.Surface, 0)
	for _, s := range surfaces {
		if s.Bounds().Overlaps(bounds) {
			matches = append(matches, s)
		}
	}
	return matches
}

func median(surfaces []render.Surface, axis int) float64 {
	vals := make([]float64, 0)
	for _, s := range surfaces {
		b := s.Bounds()
		vals = append(vals, b.MinArray[axis], b.MaxArray[axis])
	}
	sort.Float64s(vals)
	return vals[len(vals)/2]
}
