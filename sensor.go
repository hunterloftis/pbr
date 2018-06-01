package pbr

import (
	"math/rand"

	"github.com/hunterloftis/pbr/geom"
)

// Sensor generates rays from a simulated physical sensor into a Scene.
type Sensor interface {
	// Height returns the height of the Sensor film in pixels.
	Height() int
	// Width returns the width of the Sensor in pixels.
	Width() int

	Ray(x, y float64, rnd *rand.Rand) *geom.Ray3
}
