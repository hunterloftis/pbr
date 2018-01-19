package pbr

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/Opioid/rgbe"
	"github.com/hunterloftis/pbr/obj"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
	"github.com/hunterloftis/pbr/surface"
)

// Scene describes a 3d scene
type Scene struct {
	Surfaces []surface.Surface // TODO: make private
	tree     *surface.Tree
	ambient  rgb.Energy
	lights   []surface.Surface
	env      *RGBAE
}

// RGBAE Describes an rgbae (hdr) image
type RGBAE struct {
	Width  int
	Height int
	Data   []float32
	Expose float64
}

// NewScene creates and returns a pointer to an empty Scene.
func NewScene(surfaces ...surface.Surface) *Scene {
	scene := &Scene{}
	for _, s := range surfaces {
		scene.Add(s)
	}
	return scene
}

// Intersect tests whether a ray hits any objects in the scene
// TODO: implement a BSP tree and compare
func (s *Scene) Intersect(ray *geom.Ray3) surface.Hit {
	return s.tree.Intersect(ray)
	// return s.tree.IntersectSurfaces(ray) // unoptimized
}

// Env returns the light value from the environment map.
// http://gl.ict.usc.edu/Data/HighResProbes/
// http://cseweb.ucsd.edu/~wychang/cse168/
// http://www.pauldebevec.com/Probes/
func (s *Scene) EnvAt(ray *geom.Ray3) rgb.Energy {
	if s.env != nil {
		u := 1 + math.Atan2(ray.Dir.X, -ray.Dir.Z)/math.Pi
		v := math.Acos(ray.Dir.Y) / math.Pi
		x := int(u * float64(s.env.Width))
		y := int(v * float64(s.env.Height))
		index := ((y*s.env.Width + x) * 3) % len(s.env.Data)
		r := float64(s.env.Data[index])
		g := float64(s.env.Data[index+1])
		b := float64(s.env.Data[index+2])
		return rgb.Energy(geom.Vector3{r, g, b}.Scaled(s.env.Expose))
	}
	vertical := math.Max(0, (ray.Dir.Cos(geom.Direction{0, 1, 0})+0.5)/1.5)
	return rgb.Energy{}.Blend(s.ambient, vertical)
}

// Add adds new Surfaces to the scene.
func (s *Scene) Add(surfaces ...surface.Surface) {
	s.Surfaces = append(s.Surfaces, surfaces...)
}

func (s *Scene) SetAmbient(a rgb.Energy) *Scene {
	s.ambient = a
	return s
}

func (s *Scene) ReadHdrFile(file string, expose float64) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	width, height, data, err := rgbe.Decode(f)
	if err != nil {
		return err
	}
	s.SetEnv(width, height, data, expose)
	return nil
}

// TODO: implement an obj.Reader that simplifies this?
func (s *Scene) ReadObjFile(file string, thin bool) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open scene %v, %v", file, err)
	}
	defer f.Close()
	scanner := obj.NewScanner(f)
	for scanner.Scan() {
		if m := scanner.Material(); len(m) > 0 {
			matfile := filepath.Join(filepath.Dir(file), m)
			f, err := os.Open(matfile)
			if err != nil {
				continue
			}
			defer f.Close()
			err = scanner.ReadMaterials(f, thin)
			if err != nil {
				return err
			}
			continue
		}
		s.Add(scanner.Surface())
	}
	return scanner.Err()
}

func (s *Scene) Info() (box *surface.Box, surfaces []surface.Surface) {
	b := surface.BoxAround(s.Surfaces...)
	return b, s.Surfaces
}

// TODO: figure out how to deal with triangle mesh lights; ignore them? group them into a higher-level Mesh abstraction?
// too expensive to store & check each of them as a light.
func (s *Scene) Prepare() {
	s.tree = surface.NewTree(s.Surfaces)
	s.lights = make([]surface.Surface, 0)
	for _, surf := range s.Surfaces {
		if surf.Material().Emit().Average() > 0 {
			s.lights = append(s.lights, surf)
		}
	}
}

func (s *Scene) Light(rnd *rand.Rand) surface.Surface {
	i := rnd.Intn(len(s.lights))
	return s.lights[i]
}

// SetPano sets the environment to an HDR (radiance) panoramic mapping.
func (s *Scene) SetEnv(width, height int, data []float32, expose float64) {
	s.env = &RGBAE{Width: width, Height: height, Data: data, Expose: expose}
}
