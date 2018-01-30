package pbr

import (
	"math"
	"math/rand"
	"os"
	"sync/atomic"

	"github.com/Opioid/rgbe"
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/obj"
	"github.com/hunterloftis/pbr/rgb"
	"github.com/hunterloftis/pbr/surface"
)

const maxEnvEnergy = 10000

// Scene contains the elements that compose a 3D scene (Surfaces, lights, an environment map).
// A Scene can test for intersections with a Ray to see if any Scene objects were hit.
// Scene objects can be added programmatically or loaded from files.
type Scene struct {
	surfaces []surface.Surface
	tree     *surface.Tree
	ambient  rgb.Energy
	lights   []surface.Surface
	env      *rgbae
	rays     uint64
}

// rgbae describes an rgbae (hdr) image
type rgbae struct {
	width  int
	height int
	data   []float32
	expose float64
}

// NewScene constructs a Scene containing any passed Surfaces.
func NewScene(surfaces ...surface.Surface) *Scene {
	scene := &Scene{}
	for _, s := range surfaces {
		scene.Add(s)
	}
	return scene
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray *geom.Ray3) surface.Hit {
	atomic.AddUint64(&s.rays, 1)
	return s.tree.Intersect(ray)
	// return s.tree.IntersectSurfaces(ray, math.Inf(1))
}

// Rays returns the total count of Ray/Scene intersections tested since the Scene was created.
func (s *Scene) Rays() uint64 {
	return atomic.LoadUint64(&s.rays)
}

// Env returns the background energy in a given direction
// http://gl.ict.usc.edu/Data/HighResProbes/
// http://cseweb.ucsd.edu/~wychang/cse168/
// http://www.pauldebevec.com/Probes/
// TODO: correct for (aspect ratio?) so it doesn't get bendy
func (s *Scene) EnvAt(dir geom.Direction) rgb.Energy {
	if s.env != nil {
		u := 1 + math.Atan2(dir.X, -dir.Z)/math.Pi
		v := math.Acos(dir.Y) / math.Pi
		x := int(u * float64(s.env.width))
		y := int(v * float64(s.env.height))
		index := ((y*s.env.width + x) * 3) % len(s.env.data)
		r := float64(s.env.data[index])
		g := float64(s.env.data[index+1])
		b := float64(s.env.data[index+2])
		e := rgb.Energy(geom.Vector3{r, g, b}.Scaled(s.env.expose))
		return e.Limit(maxEnvEnergy)
	}
	vertical := math.Max(0, (dir.Cos(geom.Direction{0, 1, 0})+0.5)/1.5)
	return rgb.Energy{}.Blend(s.ambient, vertical)
}

// Add adds Surfaces to the scene.
func (s *Scene) Add(surfaces ...surface.Surface) {
	s.surfaces = append(s.surfaces, surfaces...)
}

// SetAmbient sets the ambient background lighting.
func (s *Scene) SetAmbient(a rgb.Energy) *Scene {
	s.ambient = a
	return s
}

// ReadHdr opens and reads a panoramic HDR image to be used as the environment map.
// You can alternatively directly set the HDR data with SetEnv().
func (s *Scene) ReadHdr(file string, expose float64) error {
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

func (s *Scene) SetEnv(width, height int, data []float32, expose float64) {
	s.env = &rgbae{width: width, height: height, data: data, expose: expose}
}

func (s *Scene) ReadObj(file string, thin bool) error {
	surfaces, err := obj.ReadFile(file, thin)
	if err != nil {
		return err
	}
	s.Add(surfaces...)
	return nil
}

// Info returns a bounding Box of all the scene's Surfaces and a list of those Surfaces.
func (s *Scene) Info() (box *surface.Box, surfaces []surface.Surface) {
	b := surface.BoxAround(s.surfaces...)
	return b, s.surfaces
}

// Light returns a random light (surface with Emit() > 0) from the Scene.
func (s *Scene) Light(rnd *rand.Rand) surface.Surface {
	i := rnd.Intn(len(s.lights))
	return s.lights[i]
}

// Lights returns the number of lights in the Scene.
func (s *Scene) Lights() int {
	return len(s.lights)
}

// TODO: figure out how to deal with triangle mesh lights; ignore them? group them into a higher-level Mesh abstraction?
// too expensive to store & check each of them as a light.
func (s *Scene) prepare() {
	s.tree = surface.NewTree(s.surfaces)
	s.lights = make([]surface.Surface, 0)
	for _, surf := range s.surfaces { // TODO: change this to a Light interface
		if surf.Material().Emit().Average() > 0 {
			s.lights = append(s.lights, surf)
		}
	}
}
