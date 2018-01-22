package geom

// Ray3 describes a 3-dimensional ray with an origin and a unit direction Vector3.
// TODO: Rename to Ray?
type Ray3 struct {
	Origin   Vector3
	Dir      Direction
	OrArray  [3]float64
	DirArray [3]float64
	InvArray [3]float64
}

func NewRay(origin Vector3, dir Direction) *Ray3 {
	return &Ray3{
		Origin:   origin,
		Dir:      dir,
		OrArray:  origin.Array(),
		DirArray: Vector3(dir).Array(),
		InvArray: [3]float64{1 / dir.X, 1 / dir.Y, 1 / dir.Z},
	}
}

// Moved returns the point dist distance along the ray
func (r *Ray3) Moved(dist float64) Vector3 {
	return r.Origin.Plus(r.Dir.Scaled(dist))
}
