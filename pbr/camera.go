package pbr

import (
	"math"
	"math/rand"
)

// Camera generates rays from a simulated physical camera into a Scene.
// The rays produced are determined by position,
// orientation, sensor type, focus, exposure, and lens selection.
type Camera struct {
	Width  int
	Height int
	CameraConfig

	focus float64
	pos   *Matrix4
}

// CameraConfig configures a camera
type CameraConfig struct {
	Lens     float64
	Sensor   float64
	Position Vector3
	Target   *Vector3
	Focus    *Vector3
	FStop    float64
}

// NewCamera makes a new Full-frame (35mm) camera.
func NewCamera(width, height int, config ...CameraConfig) *Camera {
	conf := config[0]
	if conf.Lens == 0 {
		conf.Lens = 0.050 // 50mm focal length
	}
	if conf.Sensor == 0 {
		conf.Sensor = 0.024 // height (36mm x 24mm, 35mm full frame standard)
	}
	if conf.Target == nil {
		target := conf.Position.Plus(Vector3{0, 0, -1})
		conf.Target = &target
	}
	if conf.Focus == nil {
		conf.Focus = conf.Target
	}
	if conf.FStop == 0 {
		conf.FStop = 4
	}
	return &Camera{
		Width:        width,
		Height:       height,
		CameraConfig: conf,
		focus:        conf.Focus.Minus(conf.Position).Len(),
		pos:          LookMatrix(conf.Position, *conf.Target),
	}
}

// TODO: precompute N rays for each x, y pixel & then remove Camera.focus
func (c *Camera) ray(x, y float64, rnd *rand.Rand) Ray3 {
	rx := x + rnd.Float64()
	ry := y + rnd.Float64()
	px := rx / float64(c.Width)
	py := ry / float64(c.Height)
	sensorPt := c.sensorPoint(px, py)
	straight := Vector3{}.Minus(sensorPt).Unit()
	focalPt := Vector3(straight).Scaled(c.focus)
	lensPt := c.aperturePoint(rnd)
	refracted := focalPt.Minus(lensPt).Unit()
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
	d := c.Lens / c.FStop
	t := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * d * 0.5
	x := r * math.Cos(t)
	y := r * math.Sin(t)
	return Vector3{x, y, 0}
}
