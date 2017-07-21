package collada

import (
	"encoding/xml"
	"io"
)

// Scene describes a set of geometry and materials.
type Scene struct {
	XML *Schema
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
		XML: s,
	}
	return scene, nil
}
