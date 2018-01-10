package obj

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hunterloftis/pbr"
)

type Scanner struct {
	scanner *bufio.Scanner
	next    []pbr.Surface
	err     error
	v       []pbr.Vector3
	vn      []pbr.Direction
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
		next:    make([]pbr.Surface, 0),
		v:       make([]pbr.Vector3, 0),
		vn:      make([]pbr.Direction, 0),
	}
}

func (s *Scanner) Scan() bool {
	if len(s.next) > 0 {
		return true
	}
	for s.scanner.Scan() {
		mat := pbr.Plastic(1, 1, 1, 0.7)
		line := s.scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key := fields[0]
		args := fields[1:]

		switch key {
		case "v":
			v, err := pbr.ParseVector3(strings.Join(args, ","))
			if err != nil {
				s.err = err
				return false
			}
			s.v = append(s.v, v)
		case "vn":
			vn, err := pbr.ParseDirection(strings.Join(args, ","))
			if err != nil {
				s.err = err
				return false
			}
			s.vn = append(s.vn, vn)
		case "f":
			size := len(args)
			if size < 3 || size > 4 {
				s.err = fmt.Errorf("face must contain 3 or 4 vertices, but contains %v", size)
			}
			v := make([]pbr.Vector3, size)
			n := make([]*pbr.Direction, size)
			var err error
			for i := 0; i < size; i++ {
				v[i], n[i], err = s.vertex(args[i])
				if err != nil {
					s.err = err
					return false
				}
			}
			if size == 3 {
				t := pbr.NewTriangle(v[0], v[1], v[2], mat)
				t.SetNormals(n[0], n[1], n[2])
				s.next = append(s.next, t)
				return true
			}
			if size == 4 {
				t1 := pbr.NewTriangle(v[0], v[1], v[2], mat)
				t1.SetNormals(n[0], n[1], n[2])
				t2 := pbr.NewTriangle(v[0], v[2], v[3], mat)
				t2.SetNormals(n[0], n[2], n[3])
				s.next = append(s.next, t1, t2)
				return true
			}
		}
	}
	return false
}

func (s *Scanner) Surface() (next pbr.Surface) {
	next = s.next[0]
	s.next = s.next[1:]
	return next
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) vertex(val string) (v pbr.Vector3, n *pbr.Direction, err error) {
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
