package material

import (
	"math"
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type Sample struct {
	Color    rgb.Energy
	Light    rgb.Energy
	Fresnel  rgb.Energy
	Transmit float64
	Refract  float64
	Rough    float64
	Coat     float64
	Metal    float64
	Thin     bool
}

// Bsdf is an attempt at a new bsdf
// TODO: a real BSDF instead of this procedural one.
// at each hit: choose between transmission, absorption, specular, or diffuse & generate next ray
// then pass incident & resulting directions into the bsdf to find the signal strength
// https://docs.blender.org/manual/en/dev/render/cycles/nodes/types/shaders/principled.html
// https://github.Com/wdas/brdf/blob/master/src/brdfs/disney.brdf#L131
func (s *Sample) Bsdf(norm, inc geom.Direction, dist float64, rnd *rand.Rand) (geom.Direction, rgb.Energy, bool) {
	if inc.Enters(norm) {
		// clear coat TODO: I don't think this is working. Try a red car.
		if rnd.Float64() < s.Coat {
			if rnd.Float64() < schlick(norm, inc, 0.04, 0, 0) {
				return s.shine(norm, inc, rnd)
			}
		}
		reflect := schlick(norm, inc, s.Fresnel.Average(), 0, 0)
		switch {
		// reflect
		case rnd.Float64() < reflect:
			return s.reflect(norm, inc, rnd)
		// transmit (in)
		case rnd.Float64() < s.Transmit:
			return s.transmit(norm, inc, rnd)
		// absorb
		case rnd.Float64() < s.Metal: // TODO: is this extraneous? Should s.Metal just be pre-applied to s.Color?
			return s.absorb(inc)
		// diffuse
		default:
			return s.diffuse(norm, inc, rnd)
		}
	}
	if s.Thin {
		return s.Bsdf(norm.Inv(), inc, dist, rnd)
	}
	// transmit (out)
	return s.exit(norm, inc, dist, rnd)
}

// TODO: integrate with reflect?
func (s *Sample) shine(norm, inc geom.Direction, rnd *rand.Rand) (geom.Direction, rgb.Energy, bool) {
	// TODO: if reflection enters the normal, invert the reflection about the normal
	if refl := inc.Reflected(norm); !refl.Enters(norm) {
		return refl, s.Color, false
	}
	return s.diffuse(norm, inc, rnd)
}

func (s *Sample) reflect(norm, inc geom.Direction, rnd *rand.Rand) (geom.Direction, rgb.Energy, bool) {
	// TODO: if reflection enters the normal, invert the reflection about the normal
	if refl := inc.Reflected(norm).Cone(s.Rough, rnd); !refl.Enters(norm) {
		return refl, rgb.Energy(geom.UnitVector3.Lerp(geom.Vector3(s.Fresnel), s.Metal)), false
	}
	return s.diffuse(norm, inc, rnd)
}

func (s *Sample) transmit(norm, inc geom.Direction, rnd *rand.Rand) (geom.Direction, rgb.Energy, bool) {
	if s.Thin {
		return inc, rgb.Full, false
	}
	if entered, refr := inc.Refracted(norm, 1, s.Refract); entered {
		if spread := refr.Cone(s.Rough, rnd); spread.Enters(norm) {
			return spread, rgb.Full, false
		}
		return refr, rgb.Full, false
	}
	return s.diffuse(norm, inc, rnd)
}

func (s *Sample) exit(norm, inc geom.Direction, dist float64, rnd *rand.Rand) (geom.Direction, rgb.Energy, bool) {
	if s.Transmit == 0 || s.Thin {
		// shallow bounce within margin of error
		// isn't really an intersection, so just keep the ray moving
		return inc, rgb.Full, false
	}
	tint := s.Color.Amplified(s.Transmit)
	absorb := rgb.Energy{
		X: 2 - math.Log10(tint.X*100),
		Y: 2 - math.Log10(tint.Y*100),
		Z: 2 - math.Log10(tint.Z*100),
	}
	if rnd.Float64() >= schlick(norm, inc, 0, s.Refract, 1.0) {
		if exited, refr := inc.Refracted(norm.Inv(), s.Refract, 1); exited {
			if spread := refr.Cone(s.Rough, rnd); !spread.Enters(norm) {
				return spread, beers(dist, absorb), false
			}
			return refr, beers(dist, absorb), false
		}
	}
	return inc.Reflected(norm.Inv()), beers(dist, absorb), false
}

func (s *Sample) diffuse(norm, inc geom.Direction, rnd *rand.Rand) (geom.Direction, rgb.Energy, bool) {
	return norm.RandHemiCos(rnd), s.Color.Amplified(1 / math.Pi), true
}

func (s *Sample) absorb(inc geom.Direction) (geom.Direction, rgb.Energy, bool) {
	return inc, rgb.Energy{}, false
}
