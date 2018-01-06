package pbr

import "math"

// TODO: these should all be public
type Hit struct {
	ok      bool
	surface Surface
	dist    float64
}

var Miss Hit = Hit{false, nil, math.Inf(1)}

func NewHit(surface Surface, dist float64) Hit {
	return Hit{true, surface, dist}
}

func (h Hit) Closer(h2 Hit) Hit {
	if h.dist <= h2.dist {
		return h
	}
	return h2
}

func (h Hit) Dist() float64 {
	return h.dist
}

func (h Hit) Surface() Surface {
	return h.surface
}
