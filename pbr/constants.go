package pbr

// Bias is the minimum distance unit
// Applying bias provides more robust processing of geometry
const Bias = 1e-6

// The various Elements of a pixel
// are each stored in a specific offset in flattened arrays
const (
	Red      = 0
	Green    = 1
	Blue     = 2
	Count    = 3
	Noise    = 4
	Elements = 5
)
