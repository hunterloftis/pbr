package trace

import (
	"math"
)

// Scene describes a 3d scene
type Scene struct {
	Spheres []Sphere
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray Ray3) bool {
	dist := math.MaxFloat64

	for _, sphere := range s.Spheres {
		i, d := sphere.Intersect(ray)
		if i && d < dist {
			// nearest = sphere
			dist = d
		}
	}
	if dist == math.MaxFloat64 {
		return false
	}
	return true
}

// Add adds a new object to the scene
func (s *Scene) Add(sphere Sphere) {
	s.Spheres = append(s.Spheres, sphere)
}
