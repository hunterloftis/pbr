package collada

import (
	"encoding/xml"
	"io"
)

// Scene describes a set of geometry and materials.
type Scene struct {
	XML       *Schema
	Triangles []*Triangle
}

// ReadScene creates a Scene from the Collada XML read from a given Reader.
func ReadScene(r io.Reader) (*Scene, error) {
	d := xml.NewDecoder(r)
	s := &Schema{}

	err := d.Decode(s)
	if err != nil {
		return nil, err
	}

	// m := s.mapped()
	t := make([]*Triangle, 0)

	for i := 0; i < len(s.Geometry); i++ {
		for j := 0; j < len(s.Geometry[i].Triangles); j++ {
			triangles := s.Geometry[i].Triangles[j].lookup(s)
			// material := m.material(triangles.Material)
			// indices := triangles.indices()
			// input, _ := triangles.input("VERTEX")
			// sourcePos, _ := m.source(input, "POSITION")
			// sourceNorm, _ := m.source(input, "NORMAL")
			// stride := len(triangles.Input) * 3
			for k := 0; k < triangles.el.Count; k++ {
				triangle := &Triangle{
					Pos:  triangles.vertices("POSITION", k),
					Norm: triangles.vertices("NORMAL", k),
					Mat:  triangles.material,
				}
				// start := k * stride
				// for l := 0; l < 3; l++ {
				// 	position := start + l + input.Offset
				// 	index := indices[position] * 3
				// 	triangle.Pos[l] = sourcePos.vector3(index)
				// 	triangle.Norm[l] = sourceNorm.vector3(index)
				// }
				t = append(t, triangle)
			}
		}
	}
	return &Scene{XML: s, Triangles: t}, nil
}
