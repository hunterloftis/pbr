package pbr

// Surface is an interface to all surface types (sphere, cube, etc).
// A surface is anything that can be intersected by a Ray.
type Surface interface {
	Intersect(Ray3) (hit bool, dist float64, id int)
	At(point Vector3, id int) (normal Direction, material *Material)
}
