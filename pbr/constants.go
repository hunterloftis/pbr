package pbr

// Bias is the minimum distance unit
// Applying bias provides more robust processing of geometry
const Bias = 1e-6

// Props is the number of properties stored in each float64 array
// red, green, blue, count, noise
// const Props = 5

// Props stores the various properties
var Props = map[string]int{
	"red":   0,
	"green": 1,
	"blue":  2,
	"count": 3,
	"noise": 4,
}
