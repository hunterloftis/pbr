package collada

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Scene describes a set of geometry and materials.
type Scene struct {
	XML       *Schema
	Triangles []*Triangle
}

// StringToFloats converts a space-delimited string of floats into a slice of float64.
func StringToFloats(s string) []float64 {
	fields := strings.Fields(s)
	floats := make([]float64, len(fields))
	for i := 0; i < len(fields); i++ {
		floats[i], _ = strconv.ParseFloat(fields[i], 64)
	}
	return floats
}

// StringToInts converts a space-delimited string of floats into a slice of float64.
// TODO: abstract StringTo_ with a callback conversion function?
func StringToInts(s string) []int {
	fields := strings.Fields(s)
	ints := make([]int, len(fields))
	for i := 0; i < len(fields); i++ {
		ints[i], _ = strconv.Atoi(fields[i])
	}
	return ints
}

// ReadScene creates a Scene from the Collada XML read from a given Reader.
func ReadScene(r io.Reader) (*Scene, error) {
	d := xml.NewDecoder(r)
	s := &Schema{}
	err := d.Decode(s)
	if err != nil {
		return nil, err
	}

	scene := &Scene{
		XML:       s,
		Triangles: make([]*Triangle, 0),
	}

	sources := make(map[string]*XSource)
	vertices := make(map[string]*XVertices)
	for i := 0; i < len(s.Geometry); i++ {
		for j := 0; j < len(s.Geometry[i].Source); j++ {
			id := s.Geometry[i].Source[j].ID
			sources[id] = &s.Geometry[i].Source[j]
			sources[id].floats = StringToFloats(sources[id].FloatArray.Data)
			for l := 0; l < len(sources[id].Param); l++ {
				sources[id].params += sources[id].Param[l].Name
			}
		}
		for j := 0; j < len(s.Geometry[i].Vertices); j++ {
			id := s.Geometry[i].Vertices[j].ID
			vertices[id] = &s.Geometry[i].Vertices[j]
		}
	}

	for i := 0; i < len(s.Geometry); i++ {
		for j := 0; j < len(s.Geometry[i].Triangles); j++ {
			triangles := &s.Geometry[i].Triangles[j]
			fmt.Println("geometry", i, "triangles", j, "data:", triangles.Data)
			indices := StringToInts(triangles.Data)
			fmt.Println("triangle indices:", indices)
			inputs := len(triangles.Input)
			vertexOffset := 0
			var sourcePos, sourceNorm *XSource
			for k := 0; k < inputs; k++ {
				if triangles.Input[k].Semantic == "VERTEX" {
					vertexOffset = triangles.Input[k].Offset
					vID := triangles.Input[k].Source[1:]
					fmt.Println("vertex source id:", vID)
					v := vertices[vID]
					for l := 0; l < len(v.Input); l++ {
						sID := v.Input[l].Source[1:]
						switch v.Input[l].Semantic {
						case "POSITION":
							sourcePos = sources[sID]
						case "NORMAL":
							sourceNorm = sources[sID]
						}
					}
				}
			}
			fmt.Println("source.floats:", sourcePos.floats)
			if sourcePos == nil {
				return nil, fmt.Errorf("collada: no position source found")
			}
			if sourceNorm == nil {
				return nil, fmt.Errorf("collada: no normal source found")
			}
			if sourcePos.params != "XYZ" {
				return nil, fmt.Errorf("collada: expected params XYZ, got %v", sourcePos.params)
			}
			if sourceNorm.params != "XYZ" {
				return nil, fmt.Errorf("collada: expected params XYZ, got %v", sourceNorm.params)
			}
			stride := inputs * 3
			fmt.Println("stride:", stride)
			for k := 0; k < triangles.Count; k++ {
				triangle := &Triangle{}
				start := k * stride
				for l := 0; l < 3; l++ {
					position := start + l + vertexOffset
					fmt.Println("Index stored at position", position)
					index := indices[position]
					fmt.Println("Triangle", k, "point", l, "index:", index)
					triangle.Vert[l].X = sourcePos.floats[index] // TODO: sketchup uses Z for vertical and Y for depth?
					triangle.Vert[l].Y = sourcePos.floats[index+1]
					triangle.Vert[l].Z = sourcePos.floats[index+2]
					triangle.Norm[l].X = sourceNorm.floats[index]
					triangle.Norm[l].Y = sourceNorm.floats[index+1]
					triangle.Norm[l].Z = sourceNorm.floats[index+2]
				}
				fmt.Println("Triangle", k, "verts:", triangle.Vert)
				scene.Triangles = append(scene.Triangles, triangle)
			}
		}
	}
	return scene, nil
}

// Triangle describes a 3D triangle's position, normal, and material.
type Triangle struct {
	Vert [3]Vector3
	Norm [3]Vector3
	Mat  Material
}

// Vector3 describes a 3D point in space.
type Vector3 struct {
	X, Y, Z float64
}

// Material describes the name, color, and opacity of a material.
type Material struct {
	Name       string
	R, G, B, A float64
}
