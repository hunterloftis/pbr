package pbr

import (
	"bufio"
	"io"
	"math"
	"strconv"
	"strings"

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

func (s *Scene) Info() (box *Box, center Vector3, surfaces int) {
	c := Vector3{}
	for _, s := range s.Surfaces {
		c = c.Plus(s.Center())
	}
	surfaces = len(s.Surfaces)
	center = c.Scaled(1 / float64(surfaces))
	return s.tree.box, center, surfaces
}

// TODO: make this called automatically by anything that depends on it instead of forcing that onto the visible API
func (s *Scene) Prepare() {
	s.tree = NewTree(BoxAround(s.Surfaces...), s.Surfaces, 0, "ROOT")
}

// SetPano sets the environment to an HDR (radiance) panoramic mapping.
func (s *Scene) SetPano(r io.Reader, expose float64) {
	width, height, data, _ := rgbe.Decode(r)
	s.pano = &RGBAE{Width: width, Height: height, Data: data, Expose: expose}
}

// ImportObj imports the meshes and materials from a .obj file
// http://paulbourke.net/dataformats/obj/
// https://en.wikipedia.org/wiki/Wavefront_.obj_file
// https://stackoverflow.com/questions/23723993/converting-quadriladerals-in-an-obj-file-into-triangles
// TODO: make robust
// TODO: make work with quads, not just tris
func (s *Scene) ImportObj(r io.Reader) {
	vs := make([]Vector3, 0, 1024)
	vns := make([]Direction, 0, 1024)
	scanner := bufio.NewScanner(r)
	mat := Plastic(1, 1, 1, 0.5)
	tris := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key := fields[0]
		args := fields[1:]
		switch key {
		case "v":
			v := Vector3{}
			_ = v.Set(strings.Join(args[0:3], ","))
			vs = append(vs, v)
		case "vn":
			v := Vector3{}
			_ = v.Set(strings.Join(args[0:3], ","))
			vns = append(vns, v.Unit())
		case "f":
			v1, n1 := vertex(args[0], vs, vns)
			v2, n2 := vertex(args[1], vs, vns)
			v3, n3 := vertex(args[2], vs, vns)
			if len(args) == 3 {
				t := NewTriangle(v1, v2, v3, mat)
				t.SetNormals(n1, n2, n3)
				s.Add(t)
				tris++
			} else if len(args) == 4 {
				v4, n4 := vertex(args[3], vs, vns)
				t1 := NewTriangle(v1, v2, v3, mat)
				t2 := NewTriangle(v1, v3, v4, mat)
				t1.SetNormals(n1, n2, n3)
				t2.SetNormals(n1, n3, n4)
				s.Add(t1, t2)
				tris += 2
			}
		}
	}
}

// TODO: make robust
// https://codeplea.com/triangular-interpolation
func vertex(s string, vs []Vector3, vns []Direction) (v Vector3, n *Direction) {
	fields := strings.Split(s, "/")
	vi, err := strconv.ParseInt(fields[0], 0, 0)
	if err != nil {
		panic(err)
	}
	hasNormal := len(fields) >= 3
	if hasNormal {
		ni, err := strconv.ParseInt(fields[2], 0, 0)
		if err != nil {
			panic(err)
		}
		if ni > 0 {
			n = &vns[ni-1]
		} else {
			n = &vns[len(vns)+int(ni)]
		}
	}
	if vi > 0 {
		v = vs[vi-1]
	} else {
		v = vs[len(vs)+int(vi)]
	}
	return v, n
}
