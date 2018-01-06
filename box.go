package pbr

import "math"

type Box struct {
	min, max           Vector3
	minArray, maxArray [3]float64
}

func NewBox(min, max Vector3) *Box {
	return &Box{
		min:      min,
		max:      max,
		minArray: min.Array(),
		maxArray: max.Array(),
	}
}

func MergeBoxes(a, b *Box) *Box {
	return NewBox(a.min.Min(b.min), a.max.Max(b.max))
}

func BoxAround(surfaces ...Surface) *Box {
	if len(surfaces) == 0 {
		return NewBox(Vector3{}, Vector3{})
	}
	box := surfaces[0].Box()
	for _, s := range surfaces {
		box = MergeBoxes(box, s.Box())
	}
	return box
}

// TODO: should these receivers be pointers?
func (b *Box) Overlaps(b2 *Box) bool {
	if b.min.X > b2.max.X || b.max.X < b2.min.X || b.min.Y > b2.max.Y || b.max.Y < b2.min.Y || b.min.Z > b2.max.Z || b.max.Z < b2.min.Z {
		return false
	}
	return true
}

func (b *Box) Split(axis int, val float64) (left, right *Box) {
	maxL := b.max.Array()
	minR := b.min.Array()
	maxL[axis] = val
	minR[axis] = val
	left = NewBox(b.min, ArrayToVector3(maxL))
	right = NewBox(ArrayToVector3(minR), b.max)
	return left, right
}

// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection
// http://psgraphics.blogspot.com/2016/02/new-simple-ray-box-test-from-andrew.html
func (b *Box) Check(r *Ray3) (ok bool, dist float64) {
	if b.Contains(r.Origin) {
		return true, BIAS
	}
	tmin := BIAS
	tmax := math.Inf(1)
	for a := 0; a < 3; a++ {
		t0 := (b.minArray[a] - r.orArray[a]) * r.invArray[a]
		t1 := (b.maxArray[a] - r.orArray[a]) * r.invArray[a]
		if r.invArray[a] < 0 {
			t0, t1 = t1, t0
		}
		if t0 > tmin {
			tmin = t0
		}
		if t1 < tmax {
			tmax = t1
		}
		if tmax < tmin {
			return false, math.Inf(1)
		}
	}
	return true, tmin
}

func (b *Box) Contains(p Vector3) bool {
	if p.X > b.max.X || p.X < b.min.X || p.Y > b.max.Y || p.Y < b.min.Y || p.Z > b.max.Z || p.Z < b.min.Z {
		return false
	}
	return true
}
