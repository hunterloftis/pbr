package trace

import "math/rand"

// Ray3 holds a 3d origin, unit direction
type Ray3 struct {
	Origin Vector3
	Dir    Vector3
}

// Intersect intersects with the scene
func (r *Ray3) Intersect() bool {
	if rand.Float64() < 0.1 {
		return true
	}
	return false
}
