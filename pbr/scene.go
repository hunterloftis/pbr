package pbr

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
	Surfaces []surface
	image    RGBAE
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray Ray3) (hit Hit) {
	var surf surface
	hit.Dist = math.Inf(1)

	for _, s := range s.Surfaces {
		i, d := s.Intersect(ray)
		if i && d < hit.Dist {
			hit.Dist = d
			surf = s
		}
	}
	if !math.IsInf(hit.Dist, 1) {
		hit.Point = ray.Origin.Add(ray.Dir.Scale(hit.Dist))
		hit.Mat = surf.MaterialAt(hit.Point)
		hit.Normal = surf.NormalAt(hit.Point)
	}
	return
}

// Env returns the light value from the environment map
// http://gl.ict.usc.edu/Data/HighResProbes/
func (s *Scene) Env(ray Ray3) Vector3 {
	if s.image.Width > 0 && s.image.Height > 0 {
		u := 1 + math.Atan2(ray.Dir.X, -ray.Dir.Z)/math.Pi
		v := math.Acos(ray.Dir.Y) / math.Pi
		x := int(u * float64(s.image.Width))
		y := int(v * float64(s.image.Height))
		index := ((y*s.image.Width + x) * 3) % len(s.image.Data)
		r := float64(s.image.Data[index])
		g := float64(s.image.Data[index+1])
		b := float64(s.image.Data[index+2])
		return Vector3{r, g, b}.Scale(s.image.Expose)
	}
	return Vector3{0, 0, 0}
}

// Add adds a new object to the scene
func (s *Scene) Add(surf surface) {
	s.Surfaces = append(s.Surfaces, surf)
}

// SetEnv sets the environment map
func (s *Scene) SetEnv(file string, expose float64) {
	fi, _ := os.Open(file)
	defer fi.Close()
	width, height, data, _ := rgbe.Decode(fi)
	s.image = RGBAE{Width: width, Height: height, Data: data, Expose: expose}
}
