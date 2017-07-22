package collada

import (
	"strconv"
	"strings"
)

// mapping is a map of IDs to Schema objects in a collada Schema.
type mapping struct {
	sources   map[string]*XSource
	vertices  map[string]*XVertices
	materials map[string]*XMaterial
	effects   map[string]*XEffect
	instances map[string]*XInstanceMaterial
}

// mapped maps IDs to schema objects in a Schema.
func (s *Schema) mapped() *mapping {
	m := &mapping{
		sources:   make(map[string]*XSource),
		vertices:  make(map[string]*XVertices),
		materials: make(map[string]*XMaterial),
		effects:   make(map[string]*XEffect),
		instances: make(map[string]*XInstanceMaterial),
	}

	for i := 0; i < len(s.Geometry); i++ {
		for j := 0; j < len(s.Geometry[i].Source); j++ {
			id := s.Geometry[i].Source[j].ID
			m.sources[id] = &s.Geometry[i].Source[j]
		}
		for j := 0; j < len(s.Geometry[i].Vertices); j++ {
			id := s.Geometry[i].Vertices[j].ID
			m.vertices[id] = &s.Geometry[i].Vertices[j]
		}
	}

	for i := 0; i < len(s.Material); i++ {
		id := s.Material[i].ID
		m.materials[id] = &s.Material[i]
	}

	for i := 0; i < len(s.Effect); i++ {
		id := s.Effect[i].ID
		m.effects[id] = &s.Effect[i]
	}

	for i := 0; i < len(s.VisualScene); i++ {
		for j := 0; j < len(s.VisualScene[i].InstanceGeometry); j++ {
			for k := 0; k < len(s.VisualScene[i].InstanceGeometry[j].InstanceMaterial); k++ {
				mat := &s.VisualScene[i].InstanceGeometry[j].InstanceMaterial[k]
				m.instances[mat.Symbol] = mat
			}
		}
	}

	return m
}

// material returns a Material instance given a material symbol.
func (m *mapping) material(symbol string) *Material {
	instance := m.instances[symbol]
	material := m.materials[instance.Target[1:]]
	effect := m.effects[material.InstanceEffect.URL[1:]]
	color := stringToFloats(effect.Color)
	return &Material{
		Name: material.Name,
		R:    color[0],
		G:    color[1],
		B:    color[2],
		A:    color[3],
	}
}

func (m *mapping) source(in *XInput, s string) (*XSource, bool) {
	vID := in.Source[1:]
	v := m.vertices[vID]
	for i := 0; i < len(v.Input); i++ {
		if v.Input[i].Semantic == s {
			id := v.Input[i].Source[1:]
			return m.sources[id], true
		}
	}
	return nil, false
}

// stringToFloats converts a space-delimited string of floats into a slice of float64.
func stringToFloats(s string) []float64 {
	fields := strings.Fields(s)
	floats := make([]float64, len(fields))
	for i := 0; i < len(fields); i++ {
		floats[i], _ = strconv.ParseFloat(fields[i], 64)
	}
	return floats
}
