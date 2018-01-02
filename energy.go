package pbr

import (
	"math/rand"
)

// Energy stores RGB light energy as a 3D Vector.
type Energy Vector3

// Merged merges energy b into energy a with a given signal strength.
func (a Energy) Merged(b Energy, signal Energy) Energy {
	return Energy{a.X + b.X*signal.X, a.Y + b.Y*signal.Y, a.Z + b.Z*signal.Z}
}

// Amplified returns energy a scaled by n.
func (a Energy) Amplified(n float64) Energy {
	return Energy{a.X * n, a.Y * n, a.Z * n}
}

// RandomGain randomly amplifies or destroys a signal.
// Strong signals are less likely to be destroyed and gain less amplification.
// Weak signals are more likely to be destroyed but gain more amplification.
// This creates greater overall system throughput (higher energy per signal, fewer signals).
func (a Energy) RandomGain(rnd *rand.Rand) Energy {
	max := Vector3(a).Max()
	if rnd.Float64() > max {
		return Energy{}
	}
	return a.Amplified(1 / max)
}

// Strength returns energy a multiplied by energy b.
func (a Energy) Strength(b Energy) Energy {
	return Energy{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// UnmarshalText creates an Energy from a byte array
func (a *Energy) UnmarshalText(b []byte) error {
	v := Vector3(*a)
	return (&v).UnmarshalText(b)
}

// Diff returns the difference in two Energies
func (a *Energy) Diff(b Energy) float64 {
	return Vector3(*a).Minus(Vector3(b)).Len()
}

func (a *Energy) Amount() float64 {
	return Vector3(*a).Len()
}
