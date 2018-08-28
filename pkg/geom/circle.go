package geom

import (
	"math"
	"math/rand"
)

// RandPointInCircle returns a random x, y point within a circle of radius r.
// The point is chosen uniformly and without bias.
// A random number generator must be passed in (rnd).
// https://stackoverflow.com/a/44990593/1911432
func RandPointInCircle(radius float64, rnd *rand.Rand) (x, y float64) {
	angle := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * radius
	x = r * math.Cos(angle)
	y = r * math.Sin(angle)
	return x, y
}
