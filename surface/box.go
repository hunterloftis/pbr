package surface

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
)

type Box struct {
	Min, Max           geom.Vector3
	minArray, maxArray [3]float64
	Center             geom.Vector3
	Radius             float64
}

func NewBox(min, max geom.Vector3) *Box {
	center := min.Plus(max).Scaled(0.5)
	return &Box{
		Min:      min,
		Max:      max,
		minArray: min.Array(),
		maxArray: max.Array(),
		Center:   center,
		Radius:   max.Minus(center).Len(),
	}
}

func MergeBoxes(a, b *Box) *Box {
	return NewBox(a.Min.Min(b.Min), a.Max.Max(b.Max))
}

func BoxAround(surfaces ...Surface) *Box {
	if len(surfaces) == 0 {
		return NewBox(geom.Vector3{}, geom.Vector3{})
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
	left = NewBox(b.Min, geom.ArrayToVector3(maxL))
	right = NewBox(geom.ArrayToVector3(minR), b.Max)
	return left, right
}

// https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-box-intersection
// http://psgraphics.blogspot.com/2016/02/new-simple-ray-box-test-from-andrew.html
func (b *Box) Check(r *geom.Ray3) (ok bool, near, far float64) {
	tmin := bias
	tmax := math.Inf(1)
	for a := 0; a < 3; a++ {
		t0 := (b.minArray[a] - r.OrArray[a]) * r.InvArray[a]
		t1 := (b.maxArray[a] - r.OrArray[a]) * r.InvArray[a]
		if r.InvArray[a] < 0 {
			t0, t1 = t1, t0
		}
		if t0 > tmin {
			tmin = t0
		}
		if t1 < tmax {
			tmax = t1
		}
		if tmax < tmin {
			return false, tmin, tmax
		}
	}
	return true, tmin, tmax
}

func (b *Box) Contains(p geom.Vector3) bool {
	if p.X > b.Max.X || p.X < b.Min.X || p.Y > b.Max.Y || p.Y < b.Min.Y || p.Z > b.Max.Z || p.Z < b.Min.Z {
		return false
	}
	return true
}

// RayFrom inscribes the box within a unit sphere,
// projects a solid angle disc from that sphere towards the origin,
// chooses a random point within that disc,
// and returns a Ray3 from the origin to the random point.
// https://marine.rutgers.edu/dmcs/ms552/2009/solidangle.pdf
func (b *Box) ShadowRay(origin geom.Vector3, normal geom.Direction, rnd *rand.Rand) (*geom.Ray3, float64) {
	forward := origin.Minus(b.Center).Unit()
	x, y := geom.RandPointInCircle(b.Radius, rnd) // TODO: push center back along "forward" axis, away from origin
	right := forward.Cross(geom.Up)
	up := right.Cross(forward)
	point := b.Center.Plus(right.Scaled(x)).Plus(up.Scaled(y))
	ray := geom.NewRay(origin, point.Minus(origin).Unit()) // TODO: this should be a convenience method
	dist := b.Center.Minus(origin).Len()
	cos := ray.Dir.Dot(normal)
	solidAngle := cos * (b.Radius * b.Radius) / (2 * dist * dist) // cosine-weighted ratio of disc surface area to hemisphere surface area
	return ray, solidAngle
}
