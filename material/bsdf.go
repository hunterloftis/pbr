package material

import (
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/rgb"
)

type BSDF interface {
	Sample(out geom.Direction, rnd *rand.Rand) (in geom.Direction)
	PDF(in, out geom.Direction) float64
	Eval(in, out geom.Direction) rgb.Energy
}

// func (s *Sample) BSDF() BSDF {
// 	if s.Metal == 1 {
// 		// copper: http://www.cs.cornell.edu/courses/cs5625/2013sp/lectures/Lec2ShadingModelsWeb.pdf
// 		return Microfacet{
// 			F0:        rgb.Energy{0.95, 0.64, 0.54},
// 			Roughness: 0.2,
// 		}
// 	}
// 	if s.Transmit == 1 {
// 		return Lambert{1, 0.3, 0.3}
// 	}
// 	return Lambert{0.9, 0.9, 0.9}
// }
