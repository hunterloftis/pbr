package geom

// Ray describes a 3-dimensional ray with an origin and a unit direction Vector3.
type Ray struct {
	Origin   Vec
	Dir      Dir
	OrArray  [3]float64
	DirArray [3]float64
	InvArray [3]float64
}

func NewRay(origin Vec, dir Dir) *Ray {
	return &Ray{
		Origin:   origin,
		Dir:      dir,
		OrArray:  origin.Array(),
		DirArray: Vec(dir).Array(),
		InvArray: [3]float64{1 / dir.X, 1 / dir.Y, 1 / dir.Z},
	}
}

// Moved returns the point dist distance along the ray
func (r *Ray) Moved(dist float64) Vec {
	return r.Origin.Plus(r.Dir.Scaled(dist))
}
