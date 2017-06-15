package trace

import (
	"math"
	"os"

	"github.com/Opioid/rgbe"
)

// RGBAE Describes an rgbae (hdr) image
type RGBAE struct {
	Width  int
	Height int
	Data   []float32
	Expose float64
}

// Scene describes a 3d scene
type Scene struct {
	Spheres []Sphere
	image   RGBAE
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

// Env returns the light value from the environment map
// http://gl.ict.usc.edu/Data/HighResProbes/
func (s *Scene) Env(ray Ray3) Vector3 {
	if s.image.Width > 0 && s.image.Height > 0 {
		u := 1 + math.Atan2(ray.Dir.X, -ray.Dir.Z)/math.Pi
		v := 1 - math.Acos(ray.Dir.Y)/math.Pi
		x := int(u * float64(s.image.Width))
		y := int(v * float64(s.image.Height))
		index := (y*s.image.Width + x) * 3
		r := float64(s.image.Data[index])
		g := float64(s.image.Data[index+1])
		b := float64(s.image.Data[index+2])
		return Vector3{r, g, b}.Scale(s.image.Expose)
	}
	return Vector3{0, 0, 0}
}

// Add adds a new object to the scene
func (s *Scene) Add(sphere Sphere) {
	s.Spheres = append(s.Spheres, sphere)
}

// SetEnv sets the environment map
func (s *Scene) SetEnv(file string, expose float64) {
	fi, _ := os.Open(file)
	defer fi.Close()
	width, height, data, _ := rgbe.Decode(fi)
	s.image = RGBAE{Width: width, Height: height, Data: data, Expose: expose}
}
