package pbr

// BIAS is the minimum distance unit.
// Applying bias provides more robust processing of geometry.
const BIAS = 1e-10

// AIR is the refractive index of air
const AIR = 1.00029

// UP is the unit vector orienting towards the sky
var UP = Direction{0, 1, 0}
