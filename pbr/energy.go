package pbr

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
