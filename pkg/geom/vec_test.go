package geom

import "testing"

func TestScaled(t *testing.T) {
	v := Vec{0, 1, 2}.Scaled(2)
	if v.X != 0 {
		t.Error("Expected 0, got", v.X)
	}
	if v.Y != 2 {
		t.Error("Expected 2, got", v.Y)
	}
	if v.Z != 4 {
		t.Error("Expected 4, got", v.Z)
	}
}

func TestBy(t *testing.T) {
	a := Vec{0, 1, 2}
	b := Vec{1, 2, 3}
	c := a.By(b)
	if c.X != 0 {
		t.Error("Expected 0, got", c.X)
	}
	if c.Y != 2 {
		t.Error("Expected 2, got", c.Y)
	}
	if c.Z != 6 {
		t.Error("Expected 6, got", c.Z)
	}
}
