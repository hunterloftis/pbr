package trace

import "math/rand"

var yAxis Vector3

func init() {
	yAxis = Vector3{0, 1, 0}
}

// Camera simulates a camera
type Camera struct {
	Width   int
	Height  int
	Origin  Vector3
	toWorld Matrix4
}

// Ray creates a Ray originating from the Camera
func (c *Camera) Ray(x, y int, rnd *rand.Rand) Ray3 {
	aspect := float64(c.Width) / float64(c.Height)
	rx := float64(x) + rnd.Float64()
	ry := float64(y) + rnd.Float64()
	px := (rx/float64(c.Width) - 0.5) * aspect
	py := ry/float64(c.Height) - 0.5
	projected := Vector3{px, py, 1}
	dir := projected.Minus(c.Origin).Normalize()
	dirWorld := c.toWorld.ApplyDir(dir)
	return Ray3{Origin: c.Origin, Dir: dirWorld}
}

// LookAt orients the camera
func (c *Camera) LookAt(x, y, z float64) {
	c.toWorld = NewLookMatrix4(c.Origin, Vector3{x, y, z})
}

// Move positions the camera
func (c *Camera) Move(x, y, z float64) {
	c.Origin = Vector3{x, y, z}
}
