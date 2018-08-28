package surface

import (
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/render"
)

// bias is the minimum distance unit.
// Applying bias provides more robust processing of geometry.
const bias = 1e-12

func BoundsAround(ss []render.Surface) *geom.Bounds {
	if len(ss) == 0 {
		return geom.NewBounds(geom.Vec{}, geom.Vec{})
	}
	Bounds := ss[0].Bounds()
	for _, s := range ss {
		Bounds = geom.MergeBounds(Bounds, s.Bounds())
	}
	return Bounds
}
