package pbr

import (
	"math/rand"
)

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

// Amplify randomly amplifies or destroys a signal.
// Strong signals get less amplification and are less likely to be destroyed.
// Weak signals are more likely to be destroyed but get more amplification.
// This creates greater overall system throughput (higher energy per signal, fewer signals).
func (a Energy) Amplify(rnd *rand.Rand) Energy {
	if rnd.Float64() > Vector3(a).Max() {
		return Energy{}
	}
	return a.Scaled(1 / Vector3(a).Max())
}

// Strength multiplies one energy by another
func (a Energy) Strength(b Energy) Energy {
	return Energy{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}
