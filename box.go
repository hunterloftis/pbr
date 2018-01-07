package pbr

import "math"

type Box struct {
	Min, Max           Vector3
	minArray, maxArray [3]float64
}

func NewBox(min, max Vector3) *Box {
	return &Box{
		Min:      min,
		Max:      max,
		minArray: min.Array(),
		maxArray: max.Array(),
	}
}

func MergeBoxes(a, b *Box) *Box {
	return NewBox(a.Min.Min(b.Min), a.Max.Max(b.Max))
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
	if b.Min.X > b2.Max.X || b.Max.X < b2.Min.X || b.Min.Y > b2.Max.Y || b.Max.Y < b2.Min.Y || b.Min.Z > b2.Max.Z || b.Max.Z < b2.Min.Z {
		return false
	}
	return true
}

func (b *Box) Split(axis int, val float64) (left, right *Box) {
	maxL := b.Max.Array()
	minR := b.Min.Array()
	maxL[axis] = val
	minR[axis] = val
	left = NewBox(b.Min, ArrayToVector3(maxL))
	right = NewBox(ArrayToVector3(minR), b.Max)
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
	if p.X > b.Max.X || p.X < b.Min.X || p.Y > b.Max.Y || p.Y < b.Min.Y || p.Z > b.Max.Z || p.Z < b.Min.Z {
		return false
	}
	return true
}
