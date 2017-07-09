package pbr

import "math/rand"

// Energy is a way to store RGB light energy
type Energy Vector3

// Gained does stuff
func (a Energy) Gained(b Energy, signal Energy) Energy {
	return Energy(Vector3(a).Plus(Vector3(b).By(Vector3(signal))))
}

// Scaled scales energy
func (a Energy) Scaled(n float64) Energy {
	return Energy{a.X * n, a.Y * n, a.Z * n}
}

// Destroy randomly destroys energy
func (a Energy) Destroy(rnd *rand.Rand) (Energy, bool) {
	risk := Vector3(a).Max()
	if rnd.Float64() >= risk {
		return a, true
	}
	return Energy{a.X / risk, a.Y / risk, a.Z / risk}, false
}
