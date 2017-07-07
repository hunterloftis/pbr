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
	sky      Vector3
}

// EmptyScene creates and returns a pointer to an empty Scene.
func EmptyScene() *Scene {
	return &Scene{}
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray Ray3) (hit Hit) {
	var surf Surface
	hit.Dist = math.Inf(1)

	for _, s := range s.Surfaces {
		i, d := s.Intersect(ray)
		if i && d < hit.Dist {
			hit.Dist = d
			surf = s
		}
	}
	if !math.IsInf(hit.Dist, 1) {
		hit.Point = ray.Origin.Plus(ray.Dir.Scaled(hit.Dist))
		hit.Mat = surf.MaterialAt(hit.Point)
		hit.Normal = surf.NormalAt(hit.Point)
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
	return s.sky
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
func (s *Scene) SetSky(r, g, b float64) {
	s.sky = Vector3{r, g, b}
}
