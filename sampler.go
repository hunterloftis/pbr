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
					buffer.Add(i, s.tracePrimary(x, y, rnd))
				}
				total += count
			}
			done <- sample{y, total}
		}
	}()
}

// TODO: sample Specular reflections from direct light sources and weight results by their BSDF towards the light
// Or, better, sample lights directly in general and pass that through a unified BSDF
func (s *sampler) tracePrimary(x, y int, rnd *rand.Rand) (energy rgb.Energy) {
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
		dir, signal, diffused := mat.Bsdf(normal, ray.Dir, hit.Dist, rnd)
		if diffused && lights > 0 {
			direct, coverage := s.traceDirect(lights, point, normal, rnd)
			sum = sum.Plus(direct.Strength(mat.Color))
			signal = signal.Amplified(1 - coverage)
		}
		next := geom.NewRay(point, dir)
		sum = sum.Plus(s.traceIndirect(next, 1, signal, rnd))
	}
	average := sum.Amplified(1 / float64(branch))
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
	dir, strength, diffused := mat.Bsdf(normal, ray.Dir, hit.Dist, rnd)
	if lights := s.scene.Lights(); diffused && lights > 0 {
		direct, coverage := s.traceDirect(lights, point, normal, rnd)
		energy = energy.Merged(direct.Strength(mat.Color), signal)
		signal = signal.Amplified(1 - coverage)
	}
	next := geom.NewRay(point, dir)
	return energy.Plus(s.traceIndirect(next, depth+1, signal.Strength(strength), rnd))
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
		e := hit.Surface.Material().Emit().Amplified(solidAngle / math.Pi)
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
	brightness := buffer.Average(i).Average()
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
