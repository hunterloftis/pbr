package trace

import (
	"fmt"
	"math"
)

// Scene describes a 3d scene
type Scene struct {
	spheres []Sphere
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray Ray3) bool {
	// var nearest Sphere
	dist := math.MaxFloat64

	for _, sphere := range s.spheres {
		i, d := sphere.Intersect(ray)
		if i && d < dist {
			// nearest = sphere
			dist = d
		}
	}
	if dist == math.MaxFloat64 {
		return false
	}
	fmt.Printf("Hit something\n")
	return true
}

// Add adds a new object to the scene
func (s *Scene) Add(sphere Sphere) {
	s.spheres = append(s.spheres, sphere)
}
