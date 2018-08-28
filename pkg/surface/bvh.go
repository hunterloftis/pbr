package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

// TODO: construct Bounding Volume Hierarchy
type BVH struct {
	surfs  []render.Surface
	lights []render.Object
	bounds *geom.Bounds
}

func NewBVH(ss ...render.Surface) *BVH {
	b := BVH{
		surfs:  ss,
		bounds: BoundsAround(ss),
	}
	for _, s := range b.surfs {
		b.lights = append(b.lights, s.Lights()...)
	}
	return &b
}

func (b *BVH) Intersect(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range b.surfs {
		if o, d := s.Intersect(r, dist); o != nil {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func (b *BVH) Lights() []render.Object {
	return b.lights
}

func (b *BVH) Bounds() *geom.Bounds {
	return b.bounds
}
