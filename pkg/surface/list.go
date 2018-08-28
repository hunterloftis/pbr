package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

type List struct {
	surfs  []render.Surface
	lights []render.Object
	bounds *geom.Bounds
}

func NewList(ss ...render.Surface) *List {
	l := List{
		surfs:  ss,
		bounds: BoundsAround(ss),
	}
	for _, s := range l.surfs {
		l.lights = append(l.lights, s.Lights()...)
	}
	return &l
}

func (l *List) Intersect(r *geom.Ray, max float64) (obj render.Object, dist float64) {
	dist = max
	for _, s := range l.surfs {
		if o, d := s.Intersect(r, dist); o != nil {
			obj, dist = o, d
		}
	}
	return obj, dist
}

func (l *List) Lights() []render.Object {
	return l.lights
}

func (l *List) Bounds() *geom.Bounds {
	return l.bounds
}

func (l *List) Surfaces() []render.Surface {
	return l.surfs
}
