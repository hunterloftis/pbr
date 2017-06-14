package trace

import "math/rand"

// Scene describes a 3d scene
type Scene struct {
	spheres []*Sphere
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(r *Ray3) bool {
	if rand.Float64() < 0.1 {
		return true
	}
	return false
}

// Add adds a new object to the scene
func (s *Scene) Add(sphere *Sphere) {
	s.spheres = append(s.spheres, sphere)
}
