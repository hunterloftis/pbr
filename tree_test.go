package pbr

import "testing"

func TestNew(t *testing.T) {
	m := Plastic(1, 1, 1, 1)
	s := []Surface{UnitSphere(m, Identity())}
	tree := NewTree(s, 0)
	if !tree.leaf {
		t.Error("Expected root node to be a leaf")
	}
}
