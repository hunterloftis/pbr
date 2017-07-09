package pbr

// Bias is the minimum distance unit.
// Applying bias provides more robust processing of geometry.
const Bias = 1e-6

// Pixel elements are stored in specific offsets.
// These constants allow easy access, eg `someFloat64Array[i + Blue]`
const (
	Red      = 0
	Green    = 1
	Blue     = 2
	Count    = 3
	Noise    = 4
	Elements = 5
)

// Air is the refractive index of air
const Air = 1.00029

// Up is the unit vector orienting towards the sky
var Up = Direction{0, 1, 0}
