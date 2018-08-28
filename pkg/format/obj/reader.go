package obj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hunterloftis/pbr2/pkg/format/mtl"
	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/material"
	"github.com/hunterloftis/pbr2/pkg/surface"
)

// TODO: make robust
// TODO: support smoothing groups (s)

func ReadFile(filename string, recursive bool) (*Mesh, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open scene %v, %v", filename, err)
	}
	defer f.Close()
	mesh := Read(f, filepath.Dir(filename))
	if recursive {
		ReadMaterials(mesh)
	}
	return mesh, nil
}

func ReadMaterials(mesh *Mesh) {
	lib := make(map[string]*material.Mapped)
	for _, t := range mesh.Triangles {
		if m, ok := t.Mat.(*Material); ok {
			if lib[m.Name] == nil {
				readLibraries(lib, m.Files)
			}
			if lib[m.Name] != nil {
				t.Mat = lib[m.Name]
			}
		}
	}
}

func readLibraries(lib map[string]*material.Mapped, files []string) {
	for _, f := range files {
		mats, err := mtl.ReadFile(f, true)
		if err != nil {
			fmt.Println(err)
		}
		for name, mat := range mats {
			lib[name] = mat
		}
	}
}

type tablegroup struct {
	vv []geom.Vec
	nn []geom.Dir
	tt []geom.Vec
}

func (t *tablegroup) vert(i int) geom.Vec {
	if i < 1 {
		return t.vv[len(t.vv)+i]
	}
	return t.vv[i-1]
}

func (t *tablegroup) norm(i int) geom.Dir {
	if i < 1 {
		return t.nn[len(t.nn)+i]
	}
	return t.nn[i-1]
}

func (t *tablegroup) tex(i int) geom.Vec {
	if i < 1 {
		return t.tt[len(t.tt)+i]
	}
	return t.tt[i-1]
}

func Read(r io.Reader, dir string) *Mesh {
	const (
		vertex   = "v"
		normal   = "vn"
		texture  = "vt"
		face     = "f"
		library  = "mtllib"
		material = "usemtl"
	)

	mesh := NewMesh()
	table := &tablegroup{}
	mat := &Material{}
	mats := make(map[string]*Material)
	libs := make([]string, 0)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		key, args := f[0], f[1:]

		switch key {
		case vertex:
			v, err := newVert(args)
			if err != nil {
				panic(err)
			}
			table.vv = append(table.vv, v)
		case normal:
			n, err := newNorm(args)
			if err != nil {
				panic(err)
			}
			table.nn = append(table.nn, n)
		case texture:
			t, err := newTex(args)
			if err != nil {
				panic(err)
			}
			table.tt = append(table.tt, t)
		case face:
			tris, err := newTriangles(args, table, mat)
			if err != nil {
				panic(err)
			}
			mesh.Triangles = append(mesh.Triangles, tris...)
		case library:
			libs = append(libs, strings.Join(args, " "))
		case material:
			mat = newMaterial(args, mats)
		}
	}

	for _, mat := range mats {
		mat.Files = make([]string, len(libs))
		for i, lib := range libs {
			mat.Files[i], _ = filepath.Abs(filepath.Join(dir, lib))
		}
	}

	return mesh
}

func newVert(args []string) (geom.Vec, error) {
	str := strings.Join(args[0:3], ",")
	return geom.ParseVec(str)
}

func newNorm(args []string) (geom.Dir, error) {
	str := strings.Join(args, ",")
	return geom.ParseDirection(str)
}

func newTex(args []string) (geom.Vec, error) {
	for len(args) < 3 {
		args = append(args, "0")
	}
	str := strings.Join(args, ",")
	return geom.ParseVec(str)
}

func newMaterial(args []string, mats map[string]*Material) *Material {
	if len(args) == 0 {
		return &Material{}
	}
	name := args[0]
	if mats[name] == nil {
		mats[name] = &Material{Name: name}
	}
	return mats[name]
}

func newTriangles(args []string, table *tablegroup, mat *Material) ([]*surface.Triangle, error) {
	size := len(args)
	if size < 3 {
		return nil, fmt.Errorf("face requires at least 3 vertices (contains %v)", size)
	}
	verts := make([]geom.Vec, 0)
	norms := make([]geom.Dir, 0)
	texes := make([]geom.Vec, 0)
	for _, arg := range args {
		fields := strings.Split(arg, "/")
		if i, err := parseInt(fields[0]); err == nil {
			verts = append(verts, table.vert(i))
		}
		if len(fields) < 2 {
			continue
		}
		if i, err := parseInt(fields[1]); err == nil {
			texes = append(texes, table.tex(i))
		}
		if len(fields) < 3 {
			continue
		}
		if i, err := parseInt(fields[2]); err == nil {
			norms = append(norms, table.norm(i))
		}
	}
	if len(verts) != size {
		return nil, fmt.Errorf("face vertex size != arg list size")
	}
	tris := make([]*surface.Triangle, 0)
	for i := 2; i < size; i++ {
		tri := surface.NewTriangle(verts[0], verts[i-1], verts[i], mat)
		if len(norms) == size {
			tri.SetNormals(norms[0], norms[i-1], norms[i])
		}
		if len(texes) == size {
			tri.SetTexture(texes[0], texes[i-1], texes[i])
		}
		tris = append(tris, tri)
	}
	return tris, nil
}

func parseInt(str string) (int, error) {
	i, err := strconv.ParseInt(str, 0, 0)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}
