package collada

import (
	"strconv"
	"strings"
)

// Schema is the top-level Collada XML schema.
// http://planet5.cat-v.org/
// https://github.com/GlenKelley/go-collada/blob/master/import.go
// https://larry-price.com/blog/2015/12/04/xml-parsing-in-go
// http://htmlpreview.github.io/?https://github.com/utensil/lol-model-format/blob/master/references/Collada_Tutorial_1.htm
// https://www.khronos.org/files/collada_reference_card_1_4.pdf
type Schema struct {
	Version     string         `xml:"version,attr"`
	Geometry    []XGeometry    `xml:"library_geometries>geometry"`
	Material    []XMaterial    `xml:"library_materials>material"`
	Effect      []XEffect      `xml:"library_effects>effect"`
	VisualScene []XVisualScene `xml:"library_visual_scenes>visual_scene"`
}

func (s *Schema) vertices(id string) (*XVertices, bool) {
	for i := 0; i < len(s.Geometry); i++ {
		for j := 0; j < len(s.Geometry[i].Vertices); j++ {
			if s.Geometry[i].Vertices[j].ID == id {
				return &s.Geometry[i].Vertices[j], true
			}
		}
	}
	return nil, false
}

func (s *Schema) source(id string) (*XSource, bool) {
	for i := 0; i < len(s.Geometry); i++ {
		for j := 0; j < len(s.Geometry[i].Source); j++ {
			if s.Geometry[i].Source[j].ID == id {
				return &s.Geometry[i].Source[j], true
			}
		}
	}
	return nil, false
}

// XVisualScene does something.
type XVisualScene struct {
	InstanceGeometry []XInstanceGeometry `xml:"node>instance_geometry"`
}

// XInstanceGeometry does something.
type XInstanceGeometry struct {
	InstanceMaterial []XInstanceMaterial `xml:"bind_material>technique_common>instance_material"`
}

// XInstanceMaterial maps material symbol names to material ids
type XInstanceMaterial struct {
	Symbol string `xml:"symbol,attr"`
	Target string `xml:"target,attr"`
}

// XGeometry holds scene geometry information.
type XGeometry struct {
	Source    []XSource    `xml:"mesh>source"`
	Triangles []XTriangles `xml:"mesh>triangles"`
	Vertices  []XVertices  `xml:"mesh>vertices"`
}

// XMaterial links named materials to the InstanceEffects that describe them.
type XMaterial struct {
	ID             string          `xml:"id,attr"`
	Name           string          `xml:"name,attr"`
	InstanceEffect XInstanceEffect `xml:"instance_effect"`
}

// XEffect describes a material (color, opacity).
type XEffect struct {
	ID    string `xml:"id,attr"`
	Color string `xml:"profile_COMMON>technique>lambert>diffuse>color"`
}

// XVertices holds vertex information (like position and normal) in XInput children.
type XVertices struct {
	ID    string   `xml:"id,attr"`
	Input []XInput `xml:"input"`
}

func (v *XVertices) input(semantic string) (*XInput, bool) {
	for _, in := range v.Input {
		if in.Semantic == semantic {
			return &in, true
		}
	}
	return nil, false
}

// XSource stores a flattened array of floats which map to sets of parameters (like X, Y, and Z).
type XSource struct {
	ID         string      `xml:"id,attr"`
	FloatArray XFloatArray `xml:"float_array"`
	Param      []XParam    `xml:"technique_common>accessor>param"`
}

func (s *XSource) vector3(i int) Vector3 {
	floats := stringToFloats(s.FloatArray.Data) // TODO: this is probably run a ton
	return Vector3{
		X: floats[i+offX],
		Y: floats[i+offY],
		Z: floats[i+offZ],
	}
}

// XTriangles references the named material of a triangle and the indices of the sources that describe its three points.
type XTriangles struct {
	Count    int      `xml:"count,attr"`
	Material string   `xml:"material,attr"`
	Input    []XInput `xml:"input"`
	Data     string   `xml:"p"`
}

// TrianglesLookup answers queries about triangles.
type TrianglesLookup struct {
	el       *XTriangles
	root     *Schema
	indices  []int
	inputs   map[string]*XInput
	material *Material
}

func (t *XTriangles) lookup(root *Schema) *TrianglesLookup {
	l := &TrianglesLookup{
		el:      t,
		root:    root,
		indices: stringToInts(t.Data),
		inputs:  make(map[string]*XInput), // TODO: necessary?
	}
	for i := 0; i < len(t.Input); i++ {
		l.inputs[t.Input[i].Semantic] = &t.Input[i]
	}
	symbol := l.el.Material
	var instance *XInstanceMaterial
	for i := 0; i < len(root.VisualScene); i++ {
		for j := 0; j < len(root.VisualScene[i].InstanceGeometry); j++ {
			for k := 0; k < len(root.VisualScene[i].InstanceGeometry[j].InstanceMaterial); k++ {
				mat := &root.VisualScene[i].InstanceGeometry[j].InstanceMaterial[k]
				if mat.Symbol == symbol {
					instance = mat
				}
			}
		}
	}
	var material *XMaterial
	for i := 0; i < len(root.Material); i++ {
		if root.Material[i].ID == instance.Target[1:] {
			material = &root.Material[i]
		}
	}
	var effect *XEffect
	for i := 0; i < len(root.Effect); i++ {
		if root.Effect[i].ID == material.InstanceEffect.URL[1:] {
			effect = &root.Effect[i]
		}
	}
	color := stringToFloats(effect.Color)
	l.material = &Material{
		Name: material.Name,
		R:    color[0],
		G:    color[1],
		B:    color[2],
		A:    color[3],
	}
	return l
}

func (l *TrianglesLookup) vertices(semantic string, triangle int) (v [3]Vector3) {
	input0 := l.inputs["VERTEX"]
	verts, _ := l.root.vertices(input0.Source[1:])
	input1, _ := verts.input(semantic)
	source, _ := l.root.source(input1.Source[1:])
	stride := len(l.el.Input) * 3
	start := triangle*stride + input0.Offset
	for i := 0; i < 3; i++ {
		pos := start + i
		index := l.indices[pos] * 3
		v[i] = source.vector3(index)
	}
	return
}

// stringToInts converts a space-delimited string of floats into a slice of float64.
func stringToInts(s string) []int {
	fields := strings.Fields(s)
	ints := make([]int, len(fields))
	for i := 0; i < len(fields); i++ {
		ints[i], _ = strconv.Atoi(fields[i])
	}
	return ints
}

// XInput links named meanings (like "Position") to XSource IDs (like "#ID5").
type XInput struct {
	Semantic string `xml:"semantic,attr"`
	Source   string `xml:"source,attr"`
	Offset   int    `xml:"offset,attr"`
}

// XFloatArray stores arrays of floats attached to an ID string.
type XFloatArray struct {
	ID    string `xml:"id,attr"`
	Count int    `xml:"count,attr"`
	Data  string `xml:",chardata"`
}

// XParam arrays associate an XFloatArray's data with a set of attributes (like X,Y,Z).
type XParam struct {
	Name string `xml:"name,attr"`
}

// XInstanceEffect maps a named material to a given material effect (like Lambert-diffuse) by ID.
type XInstanceEffect struct {
	URL string `xml:"url,attr"`
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
