package pbr

// http://planet5.cat-v.org/
// https://github.com/GlenKelley/go-collada/blob/master/import.go
// https://larry-price.com/blog/2015/12/04/xml-parsing-in-go
type Collada struct {
	Version  string     `xml:"attr"`
	Geometry []Geometry `xml:"library_geometries>geometry"`
}

type Geometry struct {
	Mesh Mesh `xml:"mesh"`
	test int
}

type Mesh struct {
	Source []Source `xml:"source"`
}

type Source struct {
	FloatArray FloatArray `xml:"float_array"`
}

type FloatArray struct {
	ID    string `xml:"id,attr"`
	Count int    `xml:"count,attr"`
	Data  string `xml:",chardata"`
}
