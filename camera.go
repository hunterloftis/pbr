package pbr

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/surface"
)

// Camera generates rays from a simulated physical camera into a Scene.
// The rays produced are determined by position,
// orientation, sensor type, focus, exposure, and lens selection.
type Camera struct {
	width     int
	height    int
	sensor    float64
	lens      float64
	focusDist float64
	fstop     float64
	trans     *geom.Matrix4
	position  geom.Vector3
	target    geom.Vector3
	focus     geom.Vector3
}

// NewCamera makes a new camera.
func NewCamera(width, height int) *Camera {
	c := &Camera{
		width:    width,
		height:   height,
		lens:     0.050, // 50mm focal length
		sensor:   0.024, // height (36mm x 24mm, 35mm full frame standard)
		fstop:    4,
		position: geom.Vector3{0, 0, 0},
		target:   geom.Vector3{0, 0, -1},
		focus:    geom.Vector3{0, 0, -1},
	}
	c.transform()
	return c
}

func (c *Camera) transform() {
	c.trans = geom.LookMatrix(c.position, c.target)
	c.focusDist = c.target.Minus(c.position).Len()
}

func (c *Camera) FrameDistance(box *surface.Box) float64 {
	fov := 2 * math.Atan(c.sensor/(2*c.lens))
	theta := fov / 2
	return box.Radius / math.Tan(theta)
}

func (c *Camera) LookAt(target, focus geom.Vector3) *Camera {
	c.target = target
	c.focus = focus
	c.transform()
	return c
}

// TODO: inconsistent APIs. Some take Vector3s, some take 3 floats. Unify.
func (c *Camera) MoveTo(x, y, z float64) *Camera {
	c.position = geom.Vector3{x, y, z}
	c.transform()
	return c
}

func (c *Camera) SetLens(l float64) *Camera {
	c.lens = l / 1000
	return c
}

func (c *Camera) SetStop(s float64) *Camera {
	c.fstop = s
	return c
}

func (c *Camera) Orientation() (pos, target, focus geom.Vector3) {
	return c.position, c.target, c.focus
}

// TODO: precompute N rays for each x, y pixel & then remove Camera.focus
func (c *Camera) ray(x, y float64, rnd *rand.Rand) *geom.Ray3 {
	rx := x + rnd.Float64()
	ry := y + rnd.Float64()
	px := rx / float64(c.width)
	py := ry / float64(c.height)
	sensorPt := c.sensorPoint(px, py)
	straight := geom.Vector3{}.Minus(sensorPt).Unit()
	focalPt := geom.Vector3(straight).Scaled(c.focusDist)
	lensPt := c.aperturePoint(rnd)
	refracted := focalPt.Minus(lensPt).Unit()
	ray := geom.NewRay(lensPt, refracted)
	return c.trans.MultRay(ray)
}

func (c *Camera) sensorPoint(u, v float64) geom.Vector3 {
	aspect := float64(c.width) / float64(c.height)
	w := c.sensor * aspect
	h := c.sensor
	z := 1 / ((1 / c.lens) - (1 / c.focusDist))
	x := (u - 0.5) * w
	y := (v - 0.5) * h
	return geom.Vector3{-x, y, z}
}

// https://stackoverflow.com/questions/5837572/generate-a-random-point-within-a-circle-uniformly
func (c *Camera) aperturePoint(rnd *rand.Rand) geom.Vector3 {
	d := c.lens / c.fstop
	t := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * d * 0.5
	x := r * math.Cos(t)
	y := r * math.Sin(t)
	return geom.Vector3{x, y, 0}
}

func (c *Camera) Width() int {
	return c.width
}

func (c *Camera) Height() int {
	return c.height
}
