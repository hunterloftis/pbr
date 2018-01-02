package pbr

// BIAS is the minimum distance unit.
// Applying bias provides more robust processing of geometry.
const BIAS = 1e-6

// Pixel elements are stored in specific offsets.
// These constants allow easy access, eg `someFloat64Array[i + Blue]`
const (
	Red uint = iota
	Green
	Blue
	Count
	Noise
	Stride
)

// AIR is the refractive index of air
const AIR = 1.00029

// UP is the unit vector orienting towards the sky
var UP = Direction{0, 1, 0}
