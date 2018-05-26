package pbr

import (
	"fmt"
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
				count := 1
				for c := 0; c < count; c++ {
					buffer.Add(i, s.trace(x, y, rnd))
				}
				total += count
			}
			done <- sample{y, total}
		}
	}()
}

// clamping weight: https://www.solidangle.com/research/physically_based_shader_design_in_arnold.pdf
func (s *sampler) trace(x, y int, rnd *rand.Rand) (energy rgb.Energy) {
	ray := s.camera.ray(float64(x), float64(y), rnd)
	strength := rgb.Energy{1, 1, 1}
	lights := s.scene.Lights()

	for i := 0; i < 7; i++ {
		if math.IsNaN(strength.X) || math.IsNaN(strength.Y) || math.IsNaN(strength.Z) {
			fmt.Println("starting loop with NaN strength:", strength)
			panic("starting at NaN strength")
		}
		if i > 1 {
			oldStr := strength
			if strength = strength.RandomGain(rnd); strength.Zero() {
				break
			}
			if math.IsNaN(strength.X) || math.IsNaN(strength.Y) || math.IsNaN(strength.Z) {
				fmt.Println("strength:", strength)
				fmt.Println("oldStr:", oldStr)
				fmt.Println("randomgain:", oldStr.RandomGain(rnd))
				panic("Russian roulette made strength NaN")
			}
		}
		hit := s.scene.Intersect(ray)
		if !hit.Ok {
			energy = energy.Plus(s.scene.EnvAt(ray.Dir).Times(strength))
			if math.IsNaN(energy.X) || math.IsNaN(energy.Y) || math.IsNaN(energy.Z) {
				fmt.Println("scene energy:", s.scene.EnvAt(ray.Dir))
				fmt.Println("strength:", strength)
				panic("NaN in miss")
			}
			break
		}
		point := ray.Moved(hit.Dist)
		normal, mat := hit.Surface.At(point)
		if mat.Emission > 0 {
			energy = energy.Plus(mat.Light().Times(strength))
			if math.IsNaN(energy.X) || math.IsNaN(energy.Y) || math.IsNaN(energy.Z) {
				fmt.Println("mat:", mat)
				fmt.Println("mat.Light():", mat.Light())
				fmt.Println("strength:", strength)
				panic("NaN in emission > 0")
			}
			break
		}

		toTangent, fromTangent := tangentMatrix(normal)
		wo := toTangent.MultDir(ray.Dir.Inv())
		bsdf := mat.BSDF(rnd)

		direct := 0.0
		for j := 0; j < len(lights); j++ {
			light := lights[j]
			shadow, solidAngle := light.Box().ShadowRay(point, normal, rnd)
			if solidAngle <= 0 {
				continue
			}
			direct += solidAngle
			hit := s.scene.Intersect(shadow)
			if !hit.Ok {
				continue
			}
			_, mat := hit.Surface.At(shadow.Moved(hit.Dist))
			wid := toTangent.MultDir(shadow.Dir)
			weightD := solidAngle / math.Pi
			reflectance := bsdf.Eval(wid, wo).Scaled(weightD).Times(strength)
			lightEnergy := mat.Light().Times(reflectance)
			energy = energy.Plus(lightEnergy)
		}

		wi, pdf := bsdf.Sample(wo, rnd)
		indirect := (1 - direct)
		cos := wi.Dot(geom.Up)
		weight := math.Min(30, indirect*cos/pdf)
		reflectance := bsdf.Eval(wi, wo).Scaled(weight)
		strength = strength.Times(reflectance)

		if math.IsNaN(strength.X) || math.IsNaN(strength.Y) || math.IsNaN(strength.Z) {
			fmt.Println("weight:", weight)
			fmt.Println("direct:", direct)
			fmt.Println("indirect:", indirect)
			fmt.Println("reflectance:", reflectance)
			fmt.Println("strength:", strength)
			fmt.Println("ray:", ray)
			fmt.Println("energy:", energy)
			panic("damn it, NaN")
		}

		ray = geom.NewRay(point, fromTangent.MultDir(wi))
	}
	return energy
}

// TODO: precompute on surfaces
func tangentMatrix(normal geom.Direction) (to, from *geom.Matrix4) {
	if geom.Vector3(normal).Equals(geom.Vector3(geom.Up)) {
		return geom.Identity(), geom.Identity()
	}
	angle := math.Acos(normal.Dot(geom.Up))
	axis := normal.Cross(geom.Up)
	angleAxis := axis.Scaled(angle)
	m := geom.Rot(angleAxis)
	return m, m.Inverse()
}
