package trace

// Camera simulates a camera
type Camera struct {
	origin Vector3
}

// Ray creates a Ray originating from the Camera
func (c *Camera) Ray(x, y int) Ray3 {
	px := float64(x)/960.0 - 0.5
	py := float64(y)/540.0 - 0.5
	projected := Vector3{px, py, -1}
	// dir := projected.Minus(c.origin).Normalize() // why does Go complain with this?
	diff := projected.Minus(c.origin)
	dir := diff.Normalize()
	return Ray3{Origin: c.origin, Dir: dir}
}
