package pbr

import (
	"io"
	"math"

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
	Surfaces []Surface
	pano     *RGBAE
	skyUp    Vector3
	skyDown  Vector3
}

// EmptyScene creates and returns a pointer to an empty Scene.
func EmptyScene() *Scene {
	return &Scene{}
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray Ray3) (hit bool, surf Surface, dist float64) {
	dist = math.Inf(1)
	for _, s := range s.Surfaces {
		i, d := s.Intersect(ray)
		if i && d < dist {
			hit, dist, surf = true, d, s
		}
	}
	return
}

// Env returns the light value from the environment map.
// http://gl.ict.usc.edu/Data/HighResProbes/
func (s *Scene) Env(ray Ray3) Vector3 {
	if s.pano != nil {
		u := 1 + math.Atan2(ray.Dir.X, -ray.Dir.Z)/math.Pi
		v := math.Acos(ray.Dir.Y) / math.Pi
		x := int(u * float64(s.pano.Width))
		y := int(v * float64(s.pano.Height))
		index := ((y*s.pano.Width + x) * 3) % len(s.pano.Data)
		r := float64(s.pano.Data[index])
		g := float64(s.pano.Data[index+1])
		b := float64(s.pano.Data[index+2])
		return Vector3{r, g, b}.Scaled(s.pano.Expose)
	}
	vertical := math.Max((ray.Dir.Cos(Up)+0.5)/1.5, 0)
	return s.skyDown.Lerp(s.skyUp, vertical)
}

// Add adds new Surfaces to the scene.
func (s *Scene) Add(surfaces ...Surface) {
	s.Surfaces = append(s.Surfaces, surfaces...)
}

// SetPano sets the environment to an HDR (radiance) panoramic mapping.
func (s *Scene) SetPano(r io.Reader, expose float64) {
	width, height, data, _ := rgbe.Decode(r)
	s.pano = &RGBAE{Width: width, Height: height, Data: data, Expose: expose}
}

// SetSky sets the sky color
func (s *Scene) SetSky(up, down Vector3) {
	s.skyUp = up
	s.skyDown = down
}
