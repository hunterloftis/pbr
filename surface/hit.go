package surface

import "math"

type Hit struct {
	Ok      bool
	Surface Surface
	Dist    float64
}

var Miss Hit = Hit{false, nil, math.Inf(1)}

func NewHit(Surface Surface, Dist float64) Hit {
	return Hit{true, Surface, Dist}
}

func (h Hit) Closer(h2 Hit) Hit {
	if h.Dist <= h2.Dist {
		return h
	}
	return h2
}
