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
		f, _ := strconv.ParseFloat(fields[i], 64)
		floats = append(floats, f)
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
			var source *XSource
			for k := 0; k < inputs; k++ {
				if triangles.Input[k].Semantic == "VERTEX" {
					vertexOffset = triangles.Input[k].Offset
					vID := triangles.Input[k].Source[1:]
					fmt.Println("vertex source id:", vID)
					v := vertices[vID]
					for l := 0; l < len(v.Input); l++ {
						if v.Input[l].Semantic == "POSITION" {
							sID := v.Input[l].Source[1:]
							source = sources[sID]
						}
					}
				}
			}
			fmt.Println("source:", source)
			// if source == nil {
			// 	return nil, fmt.Errorf("collada: no VERTEX Source found")
			// }
			// if source.params != "XYZ" {
			// 	return nil, fmt.Errorf("collada: expected params XYZ, got %v", source.params)
			// }
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
					// triangle.Vert[l].X = source.floats[index]
					// triangle.Vert[l].Y = source.floats[index+1]
					// triangle.Vert[l].Z = source.floats[index+2]
				}
				scene.Triangles = append(scene.Triangles, triangle)
			}
		}
	}
	return scene, nil
}

// Triangle describes a 3D triangle's position, normal, and material.
type Triangle struct {
	Vert     [3]Vector3
	Normal   Vector3
	Material Material
}

// Vector3 describes a 3D point in space.
type Vector3 struct {
	X, Y, Z float64
}

// Material describes the name, color, and opacity of a material.
type Material struct {
	Name         string
	R, G, B      float64
	Transparency float64
}
