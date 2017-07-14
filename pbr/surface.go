package pbr

// Surface is an interface to all surface types (sphere, cube, etc).
// A surface is anything that can be intersected by a Ray.
type Surface interface {
	Intersect(Ray3) (bool, float64)
	At(Vector3) (normal Direction, material *Material)
}
