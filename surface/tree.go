package surface

import (
	"math"
	"sort"

	"github.com/hunterloftis/pbr/geom"
)

// http://slideplayer.com/slide/7653218/
// https://www.cs.utexas.edu/~ckm/teaching/cs354_f11/lectures/Lecture21.pdf
// https://people.cs.clemson.edu/~dhouse/courses/405/notes/KDtrees-Fussell.pdf
type Tree struct {
	box      *Box
	left     *Tree
	right    *Tree
	axis     int
	wall     float64
	surfaces []Surface
	leaf     bool
}

func NewTree(surfaces []Surface) *Tree {
	return newBranch(BoxAround(surfaces...), surfaces, 0)
}

func newBranch(box *Box, surfaces []Surface, depth int) *Tree {
	t := &Tree{
		surfaces: overlaps(box, surfaces),
		box:      box,
	}
	limit := int(math.Max(1, math.Pow(2, float64(depth-1))))
	if len(t.surfaces) < limit || depth > 12 {
		t.leaf = true
		return t
	}
	t.axis = 0
	max := -1.0
	for i := 0; i < 3; i++ {
		length := box.Max.Axis(i) - box.Min.Axis(i)
		if length > max {
			max = length
			t.axis = i
		}
	}
	t.wall = median(t.surfaces, t.axis)
	lBox, rBox := box.Split(t.axis, t.wall)
	t.left = newBranch(lBox, t.surfaces, depth+1)
	t.right = newBranch(rBox, t.surfaces, depth+1)
	return t
}

func overlaps(box *Box, surfaces []Surface) []Surface {
	matches := make([]Surface, 0)
	for _, s := range surfaces {
		if s.Box().Overlaps(box) {
			matches = append(matches, s)
		}
	}
	return matches
}

// http://slideplayer.com/slide/7653218/
func (t *Tree) Intersect(ray *geom.Ray3) Hit {
	hit, min, max := t.box.Check(ray)
	if !hit {
		return Miss
	}
	if t.leaf {
		hit := t.IntersectSurfaces(ray, max)
		return hit
	}
	var near, far *Tree
	if ray.DirArray[t.axis] >= 0 {
		near, far = t.left, t.right
	} else {
		near, far = t.right, t.left
	}
	split := (t.wall - ray.OrArray[t.axis]) * ray.InvArray[t.axis]
	if min >= split {
		return far.Intersect(ray)
	}
	if max <= split {
		return near.Intersect(ray)
	}
	if nearHit := near.Intersect(ray); nearHit.Ok {
		return nearHit
	}
	return far.Intersect(ray)
}

func (t *Tree) IntersectSurfaces(ray *geom.Ray3, max float64) Hit {
	closest := Miss
	for _, s := range t.surfaces {
		hit := s.Intersect(ray)
		if hit.Ok && hit.Dist <= max {
			max = hit.Dist
			closest = hit
		}
	}
	return closest
}

func median(surfaces []Surface, axis int) float64 {
	vals := make([]float64, 0)
	for _, s := range surfaces {
		b := s.Box()
		vals = append(vals, b.minArray[axis], b.maxArray[axis])
	}
	sort.Float64s(vals)
	return vals[len(vals)/2]
}
