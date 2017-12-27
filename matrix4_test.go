package pbr

import "testing"

func TestInverse(t *testing.T) {
	a := Ident()
	b := a.Inverse()
	if !b.Equals(a) {
		t.Error("Identity Inverse() should be Identity.")
	}
}
