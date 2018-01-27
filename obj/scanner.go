package obj

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/material"
	"github.com/hunterloftis/pbr/surface"
)

type MatName string
type TexName string

type Scanner struct {
	scanner *bufio.Scanner
	next    []surface.Surface
	mtllib  MatName
	err     error
	v       []geom.Vector3
	vn      []geom.Direction
	vt      []geom.Vector3
	lib     map[string]*material.Map
	mat     *material.Map
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
		next:    make([]surface.Surface, 0),
		v:       make([]geom.Vector3, 0),
		vn:      make([]geom.Direction, 0),
		vt:      make([]geom.Vector3, 0),
		lib:     make(map[string]*material.Map),
		mat:     material.Plastic(1, 1, 1, 0.7),
	}
}

// http://paulbourke.net/dataformats/obj/
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
		case "vt":
			vt, err := geom.ParseVector3(args[0] + "," + args[1] + ",0") // todo: geom.Point?
			if err != nil {
				s.err = err
				return false
			}
			s.vt = append(s.vt, vt)
		case "f":
			size := len(args)
			if size < 3 {
				s.err = fmt.Errorf("face requires at least 3 vertices (contains %v)", size)
			}
			v := make([]geom.Vector3, size)
			t := make([]geom.Vector3, size)
			n := make([]*geom.Direction, size)
			var err error
			for i := 0; i < size; i++ {
				v[i], t[i], n[i], err = s.vertex(args[i])
				if err != nil {
					s.err = err
					return false
				}
			}
			for i := 2; i < size; i++ {
				tri := surface.NewTriangle(v[0], v[i-1], v[i], s.mat)
				tri.SetNormals(n[0], n[i-1], n[i])
				tri.SetTexture(t[0], t[i-1], t[i])
				s.next = append(s.next, tri)
			}
			return true
		case "mtllib":
			rest := strings.TrimSpace(line[7:])
			s.mtllib = MatName(rest)
			return true
		case "usemtl":
			if m, ok := s.lib[args[0]]; ok {
				s.mat = m
			}
		}
	}
	return false
}

func (s *Scanner) Next() interface{} {
	if len(s.mtllib) > 0 {
		lib := s.mtllib
		s.mtllib = ""
		return lib
	}
	if len(s.next) > 0 {
		surf := s.next[0]
		s.next = s.next[1:]
		return surf
	}
	return errors.New("no scan results")
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) AddMaterials(mats []*material.Map) {
	for _, mat := range mats {
		s.lib[mat.Name()] = mat
	}
}

func (s *Scanner) vertex(val string) (v, t geom.Vector3, n *geom.Direction, err error) {
	fields := strings.Split(val, "/")
	const position, texture, normal = 0, 1, 2
	if len(fields) > position {
		vIndex, err := index(fields[position], len(s.v))
		if err != nil {
			return v, t, nil, err
		}
		v = s.v[vIndex]
	} else {
		return v, t, nil, fmt.Errorf("vertex position not found: %v", val)
	}
	if len(fields) > texture {
		if len(fields[texture]) > 0 {
			tIndex, err := index(fields[texture], len(s.vt))
			if err != nil {
				return v, t, nil, err
			}
			t = s.vt[tIndex]
		}
	}
	if len(fields) > normal {
		nIndex, err := index(fields[normal], len(s.vn))
		if err != nil {
			return v, t, nil, err
		}
		n = &s.vn[nIndex]
	}
	return v, t, n, err
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
