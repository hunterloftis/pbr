package obj

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
	"github.com/hunterloftis/pbr/surface"
	"github.com/hunterloftis/pbr/surface/material"
)

type Scanner struct {
	scanner *bufio.Scanner
	next    []surface.Surface
	mtllib  string
	err     error
	v       []geom.Vector3
	vn      []geom.Direction
	lib     map[string]*material.Material
	mat     *material.Material
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
		next:    make([]surface.Surface, 0),
		v:       make([]geom.Vector3, 0),
		vn:      make([]geom.Direction, 0),
		lib:     make(map[string]*material.Material),
		mat:     material.Plastic(1, 1, 1, 0.7),
	}
}

// https://stackoverflow.com/a/18680899/1911432
// https://www.opengl.org/discussion_boards/showthread.php/198728-loading-obj-file-how-to-triangularize-polygons
func (s *Scanner) Scan() bool {
	if len(s.next) > 0 {
		return true
	}
	if len(s.mtllib) > 0 {
		return true
	}
	for s.scanner.Scan() {
		line := s.scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := fields[0]
		args := fields[1:]

		switch key {
		case "v":
			v, err := geom.ParseVector3(strings.Join(args, ","))
			if err != nil {
				s.err = err
				return false
			}
			s.v = append(s.v, v)
		case "vn":
			vn, err := geom.ParseDirection(strings.Join(args, ","))
			if err != nil {
				s.err = err
				return false
			}
			s.vn = append(s.vn, vn)
		case "f":
			size := len(args)
			if size < 3 {
				s.err = fmt.Errorf("face requires at least 3 vertices (contains %v)", size)
			}
			v := make([]geom.Vector3, size)
			n := make([]*geom.Direction, size)
			var err error
			for i := 0; i < size; i++ {
				v[i], n[i], err = s.vertex(args[i])
				if err != nil {
					s.err = err
					return false
				}
			}
			for i := 2; i < size; i++ {
				t := surface.NewTriangle(v[0], v[i-1], v[i], s.mat)
				t.SetNormals(n[0], n[i-1], n[i])
				s.next = append(s.next, t)
			}
			return true
		case "mtllib":
			s.mtllib = args[0]
			return true
		case "usemtl":
			if m, ok := s.lib[args[0]]; ok {
				s.mat = m
			}
		}
	}
	return false
}

type phong struct {
	name string
	kd   rgb.Energy // diffuse color
	tr   float64    // transmission
	ns   float64    // specular exponent
	ks   rgb.Energy // specular color
	ke   rgb.Energy // emissive color
	ni   float64    // refractive index
	pm   float64    // metal percent
}

func (s *Scanner) ReadMaterials(r io.Reader, thin bool) (err error) {
	scanner := bufio.NewScanner(r)
	mat := phong{}
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := fields[0]
		args := fields[1:]
		switch key {
		case "newmtl":
			s.addMaterial(mat, thin)
			mat = phong{name: args[0]}
		case "Kd":
			if mat.kd, err = rgb.ParseEnergy(strings.Join(args, ",")); err != nil {
				return err
			}
		case "Tr":
			if mat.tr, err = strconv.ParseFloat(args[0], 64); err != nil {
				return err
			}
		case "d":
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			mat.tr = 1 - d
		case "Ns":
			if mat.ns, err = strconv.ParseFloat(args[0], 64); err != nil {
				return err
			}
		case "Ks":
			if mat.ks, err = rgb.ParseEnergy(strings.Join(args, ",")); err != nil {
				return err
			}
		case "Ke":
			if mat.ke, err = rgb.ParseEnergy(strings.Join(args, ",")); err != nil {
				return err
			}
		case "Ni":
			if mat.ni, err = strconv.ParseFloat(args[0], 64); err != nil {
				return err
			}
		case "Pm":
			if mat.pm, err = strconv.ParseFloat(args[0], 64); err != nil {
				return err
			}
		}
	}
	s.addMaterial(mat, thin)
	return nil
}

// https://github.com/AnalyticalGraphicsInc/obj2gltf#material-types
// http://exocortex.com/blog/extending_wavefront_mtl_to_support_pbr
// TODO: refractive index (ni) => .Fresnel
func (s *Scanner) addMaterial(mat phong, thin bool) {
	if len(mat.name) == 0 {
		return
	}
	d := material.MaterialDesc{
		Color:    mat.kd,
		Transmit: mat.tr,
		Gloss:    mat.ns / 1000,
		Light:    mat.ke,
		Metal:    mat.pm,
		Thin:     thin,
	}
	if mat.tr > 0 { // TODO: don't assume all transparent objects are glass?
		if d.Thin {
			d.Transmit = mat.tr
		} else {
			d.Transmit = 1
			d.Color = d.Color.Amplified(mat.tr)
		}
		d.Fresnel = rgb.Energy{0.042, 0.042, 0.042} // Glass
	} else if d.Metal > 0 {
		d.Fresnel = mat.kd
		d.Color = mat.ks
	}
	s.lib[mat.name] = material.New(d)
}

func (s *Scanner) Surface() surface.Surface {
	next := s.next[0]
	s.next = s.next[1:]
	return next
}

func (s *Scanner) Material() string {
	lib := s.mtllib
	s.mtllib = ""
	return lib
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) vertex(val string) (v geom.Vector3, n *geom.Direction, err error) {
	fields := strings.Split(val, "/")
	const position, texture, normal = 0, 1, 2
	if len(fields) > position {
		vIndex, err := index(fields[position], len(s.v))
		if err != nil {
			return v, nil, err
		}
		v = s.v[vIndex]
	} else {
		return v, nil, fmt.Errorf("vertex position not found: %v", val)
	}
	if len(fields) > texture {
		// TODO: parse texture coords
	}
	if len(fields) > normal {
		nIndex, err := index(fields[normal], len(s.vn))
		if err != nil {
			return v, nil, err
		}
		n = &s.vn[nIndex]
	}
	return v, n, err
}

func index(s string, size int) (n int, err error) {
	i, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return 0, err
	}
	if i > 0 {
		n = int(i - 1)
	} else {
		n = size + int(i)
	}
	if n < 0 || n > size-1 {
		return 0, fmt.Errorf("index out of bounds: %v (size: %v, string: %v)", n, size, s)
	}
	return n, nil
}
