package geom

import (
	"math"
	"math/rand"
)

// https://stackoverflow.com/a/44990593/1911432
func RandPointInCircle(radius float64, rnd *rand.Rand) (x, y float64) {
	angle := 2 * math.Pi * rnd.Float64()
	r := math.Sqrt(rnd.Float64()) * radius
	x = r * math.Cos(angle)
	y = r * math.Sin(angle)
	return x, y
}
