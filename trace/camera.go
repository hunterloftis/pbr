package trace

import (
	"math"
	"math/rand"
)

// Camera simulates a camera
type Camera struct {
	Width   int
	Height  int
	Lens    float64
	Sensor  float64
	origin  Vector3
	dir     Vector3
	focus   float64
	fStop   float64
	toWorld Matrix4
}

// NewCamera makes a new Full-frame camera
func NewCamera(width, height int, lens float64) Camera {
	return Camera{
		Width:  width,
		Height: height,
		Lens:   lens,
		Sensor: 0.024, // height (24mm, full frame standard)
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

	origin := c.toWorld.Dir(lensPt).Add(c.origin) // TODO: Matrix4.Ray()
	dir := c.toWorld.Dir(refracted)

	return Ray3{Origin: origin, Dir: dir}
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
	c.dir = Vector3{x, y, z}
	c.toWorld = LookMatrix(c.origin, c.dir)
}

// Move positions the camera
func (c *Camera) Move(x, y, z float64) {
	c.origin = Vector3{x, y, z}
	c.toWorld = LookMatrix(c.origin, c.dir)
}

// Focus on a point
func (c *Camera) Focus(x, y, z, fStop float64) {
	c.focus = Vector3{x, y, z}.Minus(c.origin).Length()
	c.fStop = fStop
}
