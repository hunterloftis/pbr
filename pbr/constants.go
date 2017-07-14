package pbr

// Bias is the minimum distance unit.
// Applying bias provides more robust processing of geometry.
const Bias = 1e-6

type index int // TODO: should this single-line type go in its own file? Or something?

// Pixel elements are stored in specific offsets.
// These constants allow easy access, eg `someFloat64Array[i + Blue]`
// TODO: https://splice.com/blog/iota-elegant-constants-golang/
const (
	Red index = iota
	Green
	Blue
	Count
	Noise
	Stride
)

// Air is the refractive index of air
const Air = 1.00029

// Up is the unit vector orienting towards the sky
var Up = Direction{0, 1, 0}
