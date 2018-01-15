package pbr

import (
	"io"
	"math"

	"github.com/Opioid/rgbe"
)

// Scene describes a 3d scene
type Scene struct {
	Surfaces []Surface
	pano     *RGBAE
	skyUp    *Energy // TODO: these should be Energy
	skyDown  *Energy
	tree     *Tree
}

// RGBAE Describes an rgbae (hdr) image
type RGBAE struct {
	Width  int
	Height int
	Data   []float32
	Expose float64
}

// NewScene creates and returns a pointer to an empty Scene.
func NewScene(up, down *Energy) *Scene {
	return &Scene{
		skyUp:   up,
		skyDown: down,
	}
}

// Intersect tests whether a ray hits any objects in the scene
// TODO: implement a BSP tree and compare
func (s *Scene) Intersect(ray *Ray3) Hit {
	return s.tree.Intersect(ray)
	// return s.tree.IntersectSurfaces(ray) // unoptimized
}

func (s *Scene) Tree() *Tree {
	return s.tree
}

// Env returns the light value from the environment map.
// http://gl.ict.usc.edu/Data/HighResProbes/
// http://cseweb.ucsd.edu/~wychang/cse168/
// http://www.pauldebevec.com/Probes/
func (s *Scene) Env(ray *Ray3) Energy {
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
	vertical := math.Max(0, (ray.Dir.Cos(UP)+0.5)/1.5)
	return s.skyDown.Blend(*s.skyUp, vertical)
}

// Add adds new Surfaces to the scene.
func (s *Scene) Add(surfaces ...Surface) {
	s.Surfaces = append(s.Surfaces, surfaces...)
}

func (s *Scene) Info() (box *Box, center Vector3, surfaces []Surface) {
	b := BoxAround(s.Surfaces...)
	center = b.Min.Plus(b.Max).Scaled(0.5)
	return b, center, s.Surfaces
}

// TODO: make this called automatically by anything that depends on it instead of forcing that onto the visible API
func (s *Scene) Prepare() {
	s.tree = NewTree(s.Surfaces)
}

// SetPano sets the environment to an HDR (radiance) panoramic mapping.
func (s *Scene) SetPano(r io.Reader, expose float64) {
	width, height, data, _ := rgbe.Decode(r)
	s.pano = &RGBAE{Width: width, Height: height, Data: data, Expose: expose}
}
