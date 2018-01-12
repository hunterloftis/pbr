package pbr

// Surface is an interface to all surface types (sphere, cube, etc).
// A surface is anything that can be intersected by a Ray.
type Surface interface {
	Intersect(*Ray3) Hit
	At(point Vector3, dir Direction) (normal Direction, material *Material)
	Box() *Box
	Center() Vector3
}
