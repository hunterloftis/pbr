package pbr

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

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
	skyUp    *Energy // TODO: these should be Energy
	skyDown  *Energy
}

// NewScene creates and returns a pointer to an empty Scene.
func NewScene(up, down *Energy) *Scene {
	return &Scene{
		skyUp:   up,
		skyDown: down,
	}
}

// ColladaScene creates a Scene from collada xml data
func ColladaScene(xml []byte) *Scene {
	// fmt.Println("xml:", string(xml))
	s := Scene{}
	return &s
}

// Intersect tests whether a ray hits any objects in the scene
func (s *Scene) Intersect(ray Ray3) (hit bool, surf Surface, dist float64) {
	dist = math.Inf(1)
	for _, s := range s.Surfaces {
		h, d := s.Intersect(ray)
		if h && d < dist {
			hit, dist, surf = true, d, s // TODO: this should be an Intersection struct
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
	vertical := (ray.Dir.Cos(UP) + 1) / 2.0
	return s.skyDown.Blend(*s.skyUp, vertical)
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

// ImportObj imports the meshes and materials from a .obj file
// TODO: make robust
func (s *Scene) ImportObj(r io.Reader) {
	vs := make([]Vector3, 0, 1024)
	vns := make([]Direction, 0, 1024)
	scanner := bufio.NewScanner(r)
	mat := Glass(0.2, 1, 0.1, 0.95)
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
			t := NewTriangle(v1, v2, v3, mat)
			t.SetNormals(n1, n2, n3)
			s.Add(&t)
			tris++
		}
	}
	fmt.Println("Imported mesh with", tris, "triangles.")
}

// TODO: make robust
// https://codeplea.com/triangular-interpolation
func vertex(s string, vs []Vector3, vns []Direction) (v Vector3, n *Direction) {
	fields := strings.Split(s, "/")
	vi, err := strconv.ParseInt(fields[0], 0, 0)
	if err != nil {
		panic(err)
	}
	if len(fields) >= 3 {
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
