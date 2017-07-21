package collada

// Map is a map of IDs to Schema objects in a collada Schema.
type Map struct {
	sources   map[string]*XSource
	vertices  map[string]*XVertices
	materials map[string]*XMaterial
	effects   map[string]*XEffect
	instances map[string]*XInstanceMaterial
}

// NewMap maps IDs to schema objects in a Schema.
func NewMap(s *Schema) Map {
	m := Map{
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
			m.sources[id].floats = StringToFloats(m.sources[id].FloatArray.Data)
			for l := 0; l < len(m.sources[id].Param); l++ {
				m.sources[id].params += m.sources[id].Param[l].Name
			}
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
