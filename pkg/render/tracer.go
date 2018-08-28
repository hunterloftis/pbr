package render

import (
	"math"
	"math/rand"
	"time"

	"github.com/hunterloftis/pbr2/pkg/geom"
	"github.com/hunterloftis/pbr2/pkg/rgb"
)

const (
	maxWeight = 10
	maxEnergy = 2000
)

var (
	infinity = math.Inf(1)
)

type Camera interface {
	Ray(x, y, width, height float64, rnd *rand.Rand) *geom.Ray
}

type Environment interface {
	At(geom.Dir) rgb.Energy
}

type Surface interface {
	Intersect(r *geom.Ray, max float64) (obj Object, dist float64)
	Lights() []Object
	Bounds() *geom.Bounds
}

type Object interface {
	At(pt geom.Vec, dir geom.Dir, rnd *rand.Rand) (normal geom.Dir, bsdf BSDF)
	Bounds() *geom.Bounds
	Light() rgb.Energy    // TODO: rename to Emit()? Lumens()? <-- would need to actually be lumens in that case
	Transmit() rgb.Energy // TODO: rename to Absorb() and precompute logarithms?
}

type BSDF interface {
	Sample(wo geom.Dir, rnd *rand.Rand) (wi geom.Dir, pdf float64, shadow bool)
	Eval(wi, wo geom.Dir) rgb.Energy
}

type tracer struct {
	scene  *Scene
	out    chan *Sample
	active toggle
	rnd    *rand.Rand
	local  *Sample
	bounce int
	direct bool
}

func newTracer(s *Scene, o chan *Sample, w, h, bounce int, direct bool) *tracer {
	return &tracer{
		scene:  s,
		out:    o,
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		local:  NewSample(w, h),
		bounce: bounce,
		direct: direct,
	}
}

func (t *tracer) start() {
	if t.active.Set(true) {
		go t.process()
	}
}

func (t *tracer) stop() {
	t.active.Set(false)
}

func (t *tracer) process() {
	width := t.local.Width
	height := t.local.Height
	camera := t.scene.Camera
	for t.active.State() {
		s := NewSample(width, height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				rx := float64(x) + t.rnd.Float64()
				ry := float64(y) + t.rnd.Float64()
				r := camera.Ray(rx, ry, float64(width), float64(height), t.rnd)
				energy := t.trace(r, t.bounce).Limit(maxEnergy)
				s.Add(x, y, energy)
			}
		}
		t.local.Merge(s)
		t.out <- s
	}
}

func (t *tracer) trace(ray *geom.Ray, depth int) rgb.Energy {
	energy := rgb.Black
	signal := rgb.White

	for d := 0; d < depth; d++ {
		obj, dist := t.scene.Surface.Intersect(ray, infinity)

		if obj == nil {
			env := t.scene.Env.At(ray.Dir).Times(signal)
			energy = energy.Plus(env)
			break
		}
		if l := obj.Light(); !l.Zero() {
			energy = energy.Plus(l.Times(signal))
			break
		}

		pt := ray.Moved(dist)
		normal, bsdf := obj.At(pt, ray.Dir, t.rnd)

		if !ray.Dir.Enters(normal) {
			t := obj.Transmit()
			// if t.Zero() {	// TODO: should this be removed again now that the triangle distance bug is fixed?
			// 	ray = geom.NewRay(pt, ray.Dir)
			// 	continue
			// }
			transmittance := beers(dist, t)
			signal = signal.Times(transmittance)
		}

		toTan, fromTan := geom.Tangent(normal)
		wo := toTan.MultDir(ray.Dir.Inv())
		indirect := 1.0

		wi, pdf, shadow := bsdf.Sample(wo, t.rnd)

		if t.direct && shadow {
			dir, light, coverage := t.shadow(pt, normal)
			wiDirect := toTan.MultDir(dir)
			if coverage > 0 {
				reflectance := bsdf.Eval(wiDirect, wo).Scaled(coverage)
				e := light.Times(reflectance).Times(signal)
				energy = energy.Plus(e)
				indirect -= coverage
			}
		}

		weight := math.Min(maxWeight, indirect/pdf)
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		bounce := fromTan.MultDir(wi)
		signal = signal.Times(reflectance).RandomGain(t.rnd)

		if signal.Zero() {
			break
		}

		ray = geom.NewRay(pt, bounce)
	}

	return energy
}

// https://blog.yiningkarlli.com/2013/04/importance-sampled-direct-lighting.html
func (t *tracer) shadow(pt geom.Vec, normal geom.Dir) (wi geom.Dir, energy rgb.Energy, coverage float64) {
	lights := t.scene.Surface.Lights()
	if len(lights) < 1 {
		return geom.Up, rgb.Black, 0
	}
	i := t.rnd.Intn(len(lights))
	l := lights[i]

	ray, coverage := l.Bounds().ShadowRay(pt, normal, t.rnd)
	if coverage <= 0 {
		return geom.Up, rgb.Black, 0
	}

	obj, _ := t.scene.Surface.Intersect(ray, infinity)
	if obj == nil {
		return geom.Up, rgb.Black, 0
	}

	return ray.Dir, obj.Light(), coverage
}

// Beer's Law.
// http://www.epolin.com/converting-absorbance-transmittance
// https://en.wikipedia.org/wiki/Optical_depth
func beers(dist float64, transmit rgb.Energy) rgb.Energy {
	// Avoid edge cases
	if dist == 0 || transmit.Zero() {
		return rgb.White
	}
	// TODO: precompute this on materials, use absorption instead of transmission?
	absorb := rgb.Energy{
		X: 2 - math.Log10(transmit.X*100),
		Y: 2 - math.Log10(transmit.Y*100),
		Z: 2 - math.Log10(transmit.Z*100),
	}
	r := math.Exp(-absorb.X * dist)
	g := math.Exp(-absorb.Y * dist)
	b := math.Exp(-absorb.Z * dist)
	return rgb.Energy{r, g, b}
}
