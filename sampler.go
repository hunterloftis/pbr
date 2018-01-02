package pbr

import (
	"math/rand"
	"time"
)

// Sampler samples pixels for a Scene by tracing Rays from a Camera.
type Sampler struct {
	Width  int
	Height int
	SamplerConfig

	cam   *Camera
	scene *Scene
	rnd   *rand.Rand
}

// SamplerConfig configures a Sampler.
type SamplerConfig struct {
	Bounces  int
	Direct   uint // TODO
	Indirect uint // TODO
}

type result struct {
	index  uint
	energy Energy
}

// NewSampler constructs a new Sampler for a given Camera and Scene.
// The Sampler samples Rays from the Camera into the Scene.
// bounces specifies the maximum number of times a Ray can bounce around the scene (eg, 10).
// adapt specifies how adaptive sampling should be to noise (0 = none, 3 = medium, 4 = high).
func NewSampler(cam *Camera, scene *Scene, config ...SamplerConfig) *Sampler {
	conf := SamplerConfig{}
	if len(config) > 0 {
		conf = config[0]
	}
	if conf.Bounces == 0 {
		conf.Bounces = 10 // Reasonable default
	}
	return &Sampler{
		Width:         cam.Width,
		Height:        cam.Height,
		SamplerConfig: conf,
		cam:           cam,
		scene:         scene,
		rnd:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Sampler) Sample(in <-chan uint, out chan<- result) {
	size := uint(s.cam.Width * s.cam.Height)
	go func() {
		for {
			if p, ok := <-in; ok {
				i := p % size
				x, y := s.pixelAt(i)
				sample := s.trace(x, y)
				out <- result{i, sample}
			} else {
				return
			}
		}
	}()
}

func (s *Sampler) trace(x, y float64) Energy {
	ray := s.cam.ray(x, y, s.rnd)
	signal := Energy{1, 1, 1}
	energy := Energy{0, 0, 0}

	for bounce := 0; bounce < s.Bounces; bounce++ {
		hit, surface, dist := s.scene.Intersect(ray)
		if !hit {
			energy = energy.Merged(s.scene.Env(ray), signal)
			break
		}
		point := ray.Moved(dist)
		normal, mat := surface.At(point)
		energy = energy.Merged(mat.Emit(normal, ray.Dir), signal)
		signal = signal.RandomGain(s.rnd) // "Russian Roulette"
		if signal == (Energy{}) {
			break
		}
		if next, dir, str := mat.Bsdf(normal, ray.Dir, dist, s.rnd); next {
			signal = signal.Strength(str)
			ray = Ray3{point, dir}
		} else {
			break
		}
	}
	return energy
}

func (s *Sampler) pixelAt(i uint) (x, y float64) {
	return float64(i % uint(s.Width)), float64(i / uint(s.Width))
}
