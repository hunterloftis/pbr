package pbr

import (
	"math"
	"math/rand"
	"time"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type sampler struct {
	adapt   float64
	bounces int
	direct  int
	branch  int
	camera  *Camera
	scene   *Scene
}

type sample struct {
	row   int
	count int
}

func (s *sampler) start(buffer *rgb.Framebuffer, in <-chan int, done chan<- sample) {
	width := uint(s.camera.Width())
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for y := range in {
			total := 0
			for x := 0; x < int(width); x++ {
				i := uint(y*int(width) + x)
				count := s.adapted(buffer, uint(i))
				for c := 0; c < count; c++ {
					buffer.Add(i, s.trace(x, y, rnd))
				}
				total += count
			}
			done <- sample{y, total}
		}
	}()
}

func tangentMatrix(view, normal geom.Direction) (to, from *geom.Matrix4) {
	t := view.Cross(normal)
	b := t.Cross(normal)
	n := normal
	m := geom.NewMatrix4(
		t.X, b.X, n.X, 0,
		t.Y, b.Y, n.Y, 0,
		t.Z, b.Z, n.Z, 0,
		0, 0, 0, 1,
	)
	return m, m.Inverse()
}

func (s *sampler) trace(x, y int, rnd *rand.Rand) (energy rgb.Energy) {
	ray := s.camera.ray(float64(x), float64(y), rnd)
	strength := rgb.Energy{1, 1, 1}
	// lights := s.scene.Lights() TODO: direct lighting

	for i := 0; i < 9; i++ {
		hit := s.scene.Intersect(ray)
		if !hit.Ok {
			energy = energy.Plus(s.scene.EnvAt(ray.Dir).Times(strength))
			break
		}
		point := ray.Moved(hit.Dist)
		normal, mat := hit.Surface.At(point)
		if !mat.Light.Zero() {
			energy = energy.Plus(mat.Light.Times(strength))
			break
		}
		bsdf := mat.BSDF()
		view := ray.Dir.Inv()
		toTangent, fromTangent := tangentMatrix(view, normal)
		wo := toTangent.MultDir(view)
		wi := bsdf.Sample(wo, rnd)
		weight := wi.Dot(geom.Up) / bsdf.PDF(wi, wo)
		strength = strength.Times(bsdf.Eval(wi, wo)).Scaled(weight)
		ray = geom.NewRay(point, fromTangent.MultDir(wi))
	}
	return energy
}

// TODO: sample Specular reflections from direct light sources and weight results by their BSDF towards the light
// Or, better, sample lights directly in general and pass that through a unified BSDF
func (s *sampler) tracePrimary2(x, y int, rnd *rand.Rand) (energy rgb.Energy) {
	ray := s.camera.ray(float64(x), float64(y), rnd)
	hit := s.scene.Intersect(ray)
	if !hit.Ok {
		return s.scene.EnvAt(ray.Dir)
	}
	point := ray.Moved(hit.Dist)
	normal, mat := hit.Surface.At(point)
	energy = energy.Plus(mat.Light)
	branch := 1 + int(float64(s.branch)*(mat.Rough+0.25)/1.25)
	sum := rgb.Energy{}
	lights := s.scene.Lights()
	for i := 0; i < branch; i++ {
		dir, signal, diffused := mat.Bsdf2(normal, ray.Dir, hit.Dist, rnd)
		if diffused && lights > 0 {
			direct, coverage := s.traceDirect(lights, point, normal, rnd)
			sum = sum.Plus(direct.Times(mat.Color))
			signal = signal.Scaled(1 - coverage)
		}
		next := geom.NewRay(point, dir)
		sum = sum.Plus(s.traceIndirect(next, 1, signal, rnd))
	}
	average := sum.Scaled(1 / float64(branch))
	return energy.Plus(average)
}

func (s *sampler) traceIndirect(ray *geom.Ray3, depth int, signal rgb.Energy, rnd *rand.Rand) (energy rgb.Energy) {
	if depth >= s.bounces {
		return energy
	}
	if signal = signal.RandomGain(rnd); signal.Zero() {
		return energy
	}
	hit := s.scene.Intersect(ray)
	if !hit.Ok {
		energy = energy.Merged(s.scene.EnvAt(ray.Dir), signal)
		return energy
	}
	point := ray.Moved(hit.Dist)
	normal, mat := hit.Surface.At(point)
	energy = energy.Merged(mat.Light, signal)
	dir, strength, diffused := mat.Bsdf2(normal, ray.Dir, hit.Dist, rnd)
	if lights := s.scene.Lights(); diffused && lights > 0 {
		direct, coverage := s.traceDirect(lights, point, normal, rnd)
		energy = energy.Merged(direct.Times(mat.Color), signal)
		signal = signal.Scaled(1 - coverage)
	}
	next := geom.NewRay(point, dir)
	return energy.Plus(s.traceIndirect(next, depth+1, signal.Times(strength), rnd))
}

func (s *sampler) traceDirect(num int, point geom.Vector3, normal geom.Direction, rnd *rand.Rand) (energy rgb.Energy, coverage float64) {
	limit := int(math.Min(float64(s.direct), float64(num)))
	for i := 0; i < limit; i++ {
		light := s.scene.Light(rnd)
		ray, solidAngle := light.Box().ShadowRay(point, normal, rnd)
		if solidAngle <= 0 {
			break
		}
		coverage += solidAngle
		hit := s.scene.Intersect(ray)
		if !hit.Ok {
			break
		}
		e := hit.Surface.Material().Emit().Scaled(solidAngle / math.Pi)
		energy = energy.Plus(e)
	}
	return energy, coverage
}

// http://gfx.cs.princeton.edu/pubs/DeCoro_2010_DOR/outliers.pdf
// TODO: backgrounds should be basically completely black on the heatmap
func (s *sampler) adapted(buffer *rgb.Framebuffer, i uint) int {
	if s.adapt == 0 {
		return 1
	}
	count := buffer.Count(i)
	if count < 3 {
		return 1
	}
	needs := 1
	brightness := buffer.Average(i).Mean()
	midtones := (((255 - math.Min(brightness, 255)) / 255) + 3) / 4
	noise := buffer.Noise(i)
	varMean, countMean := buffer.Variance()
	ratio := (noise + 1) / (varMean + 1)
	targetCount := math.Ceil(ratio * countMean * midtones)
	correction := targetCount - count
	limited := math.Max(0, math.Min(s.adapt, correction))
	needs += int(limited)
	return needs
}
