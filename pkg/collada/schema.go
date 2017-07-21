package collada

// Schema is the top-level Collada XML schema.
// http://planet5.cat-v.org/
// https://github.com/GlenKelley/go-collada/blob/master/import.go
// https://larry-price.com/blog/2015/12/04/xml-parsing-in-go
// http://htmlpreview.github.io/?https://github.com/utensil/lol-model-format/blob/master/references/Collada_Tutorial_1.htm
// https://www.khronos.org/files/collada_reference_card_1_4.pdf
type Schema struct {
	Version  string      `xml:"version,attr"`
	Geometry []XGeometry `xml:"library_geometries>geometry"`
	Material []XMaterial `xml:"library_materials>material"`
	Effect   []XEffect   `xml:"library_effects>effect"`
}

// XGeometry holds scene geometry information.
type XGeometry struct {
	Source    []XSource    `xml:"mesh>source"`
	Triangles []XTriangles `xml:"mesh>triangles"`
	Input     []XInput     `xml:"mesh>vertices>input"`
}

// XMaterial links named materials to the InstanceEffects that describe them.
type XMaterial struct {
	Name           string          `xml:"name,attr"`
	InstanceEffect XInstanceEffect `xml:"instance_effect"`
}

// XEffect describes a material (color, opacity).
type XEffect struct {
	ID    string `xml:"id,attr"`
	Color string `xml:"profile_COMMON>technique>lambert>diffuse>color"`
}

// XSource stores a flattened array of floats which map to sets of parameters (like X, Y, and Z).
type XSource struct {
	ID         string      `xml:"id,attr"`
	FloatArray XFloatArray `xml:"float_array"`
	Param      []XParam    `xml:"technique_common>accessor>param"`

	floats []float64
	params string
}

// XTriangles references the named material of a triangle and the indices of the sources that describe its three points.
type XTriangles struct {
	Count    int      `xml:"count,attr"`
	Material string   `xml:"material,attr"`
	Input    []XInput `xml:"input"`
	Data     string   `xml:"p"`
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
