package pbr

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// http://slideplayer.com/slide/7653218/
// https://www.cs.utexas.edu/~ckm/teaching/cs354_f11/lectures/Lecture21.pdf
// https://people.cs.clemson.edu/~dhouse/courses/405/notes/KDtrees-Fussell.pdf
type Tree struct {
	box      *Box
	left     *Tree
	right    *Tree
	surfaces []Surface
	leaf     bool
	name     string
}

func NewTree(surfaces []Surface) *Tree {
	return newBranch(BoxAround(surfaces...), surfaces, 0, "ROOT")
}

func newBranch(box *Box, surfaces []Surface, depth int, suffix string) *Tree {
	t := &Tree{
		surfaces: overlaps(box, surfaces),
		box:      box,
		name:     strconv.Itoa(depth) + suffix,
	}
	t.name += " (" + fmt.Sprintf("%v", &t) + ")"
	limit := int(math.Max(1, math.Pow(2, float64(depth-1))))
	if len(t.surfaces) < limit || depth > 12 {
		t.leaf = true
		return t
	}
	// TODO: try using maximum dimension instead
	axis := depth % 3
	wall := median(t.surfaces, axis)
	lBox, rBox := box.Split(axis, wall)
	t.left = newBranch(lBox, t.surfaces, depth+1, "L")
	t.right = newBranch(rBox, t.surfaces, depth+1, "R")
	return t
}

// TODO: implement as Box.Overlaps(Box) and replace Surface.Bounds() with Surface.Box() which actually returns a Box
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
// TODO: use tmin, tmax, tsplit to optimize this (vs Box.Check)
func (t *Tree) Intersect(ray *Ray3) Hit {
	if t.leaf {
		hit := t.IntersectSurfaces(ray)
		return hit
	}
	left, lDist := t.left.Check(ray)
	right, rDist := t.right.Check(ray)
	if left && right {
		var near, far *Tree
		if lDist <= rDist {
			near, far = t.left, t.right
		} else {
			near, far = t.right, t.left
		}
		if hit := near.Intersect(ray); hit.ok {
			return hit
		}
		return far.Intersect(ray)
	} else if left {
		return t.left.Intersect(ray)
	} else if right {
		return t.right.Intersect(ray)
	}
	return Miss
}

func (t *Tree) IntersectSurfaces(ray *Ray3) Hit {
	closest := Miss
	for _, s := range t.surfaces {
		hit := s.Intersect(ray)
		if hit.ok {
			if t.box.Contains(ray.Moved(hit.dist)) {
				closest = hit.Closer(closest)
			}
		}
	}
	return closest
}

func (t *Tree) Check(ray *Ray3) (ok bool, dist float64) {
	return t.box.Check(ray)
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
