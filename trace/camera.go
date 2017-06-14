package trace

import "math/rand"

// Camera simulates a camera
type Camera struct {
	Origin Vector3
	Width  int
	Height int
}

// Ray creates a Ray originating from the Camera
func (c *Camera) Ray(x, y int) Ray3 {
	aspect := float64(c.Width) / float64(c.Height)
	rx := float64(x) + rand.Float64()
	ry := float64(y) + rand.Float64()
	px := (rx/float64(c.Width) - 0.5) * aspect
	py := ry/float64(c.Height) - 0.5
	projected := Vector3{px, py, -1}
	// dir := projected.Minus(c.origin).Normalize() // why does Go complain with this?
	diff := projected.Minus(c.Origin)
	dir := diff.Normalize()
	return Ray3{Origin: c.Origin, Dir: dir}
}
