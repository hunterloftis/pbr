package pbr

import (
	"encoding/xml"
	"fmt"
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
	skyUp    Vector3 // TODO: these should be Energy
	skyDown  Vector3
}

// EmptyScene creates and returns a pointer to an empty Scene.
func EmptyScene() *Scene {
	return &Scene{}
}

// ColladaScene creates a Scene from collada xml data
func ColladaScene(r io.Reader) *Scene {
	// fmt.Println("xml:", string(xml))
	d := xml.NewDecoder(r)
	collada := &Collada{}
	_ = d.Decode(collada)
	fmt.Println(collada)
	s := Scene{}
	return &s
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
func (s *Scene) Env(ray Ray3) Energy {
	if s.pano != nil {
		u := 1 + math.Atan2(ray.Dir.X, -ray.Dir.Z)/math.Pi
		v := math.Acos(ray.Dir.Y) / math.Pi
		x := int(u * float64(s.pano.Width))
		y := int(v * float64(s.pano.Height))
		index := ((y*s.pano.Width + x) * 3) % len(s.pano.Data)
		r := float64(s.pano.Data[index])
		g := float64(s.pano.Data[index+1])
		b := float64(s.pano.Data[index+2])
		return Energy(Vector3{r, g, b}.Scaled(s.pano.Expose))
	}
	vertical := math.Max((ray.Dir.Cos(Up)+0.5)/1.5, 0)
	return Energy(s.skyDown.Lerp(s.skyUp, vertical))
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

// SetSky sets the top (up) and bottom (down) sky color.
func (s *Scene) SetSky(up, down Vector3) {
	s.skyUp = up
	s.skyDown = down
}
