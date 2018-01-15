package pbr

import (
	"math/rand"
)

// Energy stores RGB light energy as a 3D Vector.
type Energy Vector3

var Energy1 = Energy{1, 1, 1}

// Merged merges energy b into energy a with a given signal strength.
func (a Energy) Merged(b Energy, signal Energy) Energy {
	return Energy{a.X + b.X*signal.X, a.Y + b.Y*signal.Y, a.Z + b.Z*signal.Z}
}

// Amplified returns energy a scaled by n.
// TODO: "amplified" is misleading, should be "scaled"
func (a Energy) Amplified(n float64) Energy {
	return Energy{a.X * n, a.Y * n, a.Z * n}
}

// RandomGain randomly amplifies or destroys a signal.
// Strong signals are less likely to be destroyed and gain less amplification.
// Weak signals are more likely to be destroyed but gain more amplification.
// This creates greater overall system throughput (higher energy per signal, fewer signals).
func (a Energy) RandomGain(rnd *rand.Rand) Energy {
	greatest := Vector3(a).Greatest()
	if rnd.Float64() > greatest {
		return Energy{}
	}
	return a.Amplified(1 / greatest)
}

// Strength returns energy a multiplied by energy b.
func (a Energy) Strength(b Energy) Energy {
	return Energy{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Diff returns the difference in two Energies
func (a Energy) Variance(b Energy) float64 {
	d := Vector3(a).Minus(Vector3(b))
	return d.X*d.X + d.Y*d.Y + d.Z*d.Z
}

func (a Energy) Average() float64 {
	return (a.X + a.Y + a.Z) / 3
}

func (a Energy) Blend(b Energy, n float64) Energy {
	return Energy(Vector3(a).Lerp(Vector3(b), n))
}

func (a *Energy) Set(b Energy) {
	a.X = b.X
	a.Y = b.Y
	a.Z = b.Z
}

func (a *Energy) UnmarshalText(b []byte) error {
	v, err := ParseVector3(string(b))
	if err != nil {
		return err
	}
	a.Set(Energy(v))
	return nil
}

func ParseEnergy(s string) (e Energy, err error) {
	v, err := ParseVector3(s)
	return Energy(v), err
}
