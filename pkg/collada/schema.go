package collada

// http://planet5.cat-v.org/
// https://github.com/GlenKelley/go-collada/blob/master/import.go
// https://larry-price.com/blog/2015/12/04/xml-parsing-in-go
// http://htmlpreview.github.io/?https://github.com/utensil/lol-model-format/blob/master/references/Collada_Tutorial_1.htm
type Schema struct {
	Version  string     `xml:"attr"`
	Geometry []Geometry `xml:"library_geometries>geometry"`
	Material []Material `xml:"library_materials>material"`
	Effect   []Effect   `xml:"library_effects>effect"`
}

type Geometry struct {
	Source    []Source    `xml:"mesh>source"`
	Triangles []Triangles `xml:"mesh>triangles"`
	Input     []Input     `xml:"mesh>vertices>input"`
}

type Material struct {
	Name           string         `xml:"name,attr"`
	InstanceEffect InstanceEffect `xml:"instance_effect"`
}

type Effect struct {
	ID    string `xml:"id,attr"`
	Color string `xml:"profile_COMMON>technique>lambert>diffuse>color"`
}

type Source struct {
	FloatArray FloatArray `xml:"float_array"`
	Param      []Param    `xml:"technique_common>accessor>param"` // Order (usually X,Y,Z)
}

type Triangles struct {
	Count    int     `xml:"count,attr"`
	Material string  `xml:"material,attr"`
	Input    []Input `xml:"input"`
	Data     string  `xml:"p"` // Indices => corresponding Sources
}

type Input struct {
	Semantic string `xml:"semantic,attr"`
	Source   string `xml:"source,attr"`
	Offset   int    `xml:"offset,attr"`
}

type FloatArray struct {
	ID    string `xml:"id,attr"`
	Count int    `xml:"count,attr"`
	Data  string `xml:",chardata"`
}

type Param struct {
	Name string `xml:"name,attr"`
}

type InstanceEffect struct {
	URL string `xml:"url,attr"`
}
