package pbr

// Ray3 describes a 3-dimensional ray with an origin and a unit direction Vector3.
type Ray3 struct {
	Origin Vector3
	Dir    Vector3
}

// Move returns the point dist distance along the ray
func (r Ray3) Move(dist float64) Vector3 {
	return r.Origin.Plus(r.Dir.Scaled(dist))
}
