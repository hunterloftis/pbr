package pbr

import (
	"math"
	"math/rand"
)

// Camera simulates a camera
type Camera struct {
	Width  int
	Height int
	Lens   float64
	Sensor float64
	origin Vector3
	target Vector3
	focus  float64
	fStop  float64
	pos    Matrix4
}

// Camera35mm makes a new Full-frame camera
func Camera35mm(width, height int, lens float64) Camera {
	return Camera{
		Width:  width,
		Height: height,
		Lens:   lens,
		Sensor: 0.024, // height (36mm x 24mm, 35mm full frame standard)
	}
}

// Ray creates a Ray originating from the Camera
func (c *Camera) Ray(x, y float64, rnd *rand.Rand) Ray3 {
	rx := x + rnd.Float64()
	ry := y + rnd.Float64()
	px := rx / float64(c.Width)
	py := ry / float64(c.Height)
	sensorPt := c.sensorPoint(px, py)
	straight := Vector3{}.Minus(sensorPt).Normalize()
	focalPt := straight.Scale(c.focus)
	lensPt := c.aperturePoint(rnd)
	refracted := focalPt.Minus(lensPt).Normalize()
	ray := Ray3{Origin: lensPt, Dir: refracted}
	return c.pos.MultRay(ray)
}

func (c *Camera) sensorPoint(u, v float64) Vector3 {
	aspect := float64(c.Width) / float64(c.Height)
	w := c.Sensor * aspect
	h := c.Sensor
	z := 1 / ((1 / c.Lens) - (1 / c.focus))
	x := (u - 0.5) * w
	y := (v - 0.5) * h
	return Vector3{-x, y, z}
}

// https://stackoverflow.com/questions/5837572/generate-a-random-point-within-a-circle-uniformly
func (c *Camera) aperturePoint(rnd *rand.Rand) Vector3 {
	d := c.Lens / c.fStop
	t := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * d * 0.5
	x := r * math.Cos(t)
	y := r * math.Sin(t)
	return Vector3{x, y, 0}
}

// LookAt orients the camera
func (c *Camera) LookAt(x, y, z float64) {
	c.target = Vector3{x, y, z}
	c.pos = LookMatrix(c.origin, c.target)
}

// MoveTo positions the camera
func (c *Camera) MoveTo(x, y, z float64) {
	c.origin = Vector3{x, y, z}
	c.pos = LookMatrix(c.origin, c.target)
}

// Focus on a point
func (c *Camera) Focus(x, y, z, fStop float64) {
	c.focus = Vector3{x, y, z}.Minus(c.origin).Length()
	c.fStop = fStop
}
