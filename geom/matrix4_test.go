package geom

import "testing"

func TestInverse(t *testing.T) {
	a := Identity()
	b := a.Inverse()
	if !b.Equals(a) {
		t.Error("Identity Inverse() should be Identity.")
	}
}
