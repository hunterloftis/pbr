package trace

// Camera simulates a camera
type Camera struct {
}

// Ray creates a Ray originating from the Camera
func (c *Camera) Ray(x, y int) *Ray3 {
	return &Ray3{}
}
