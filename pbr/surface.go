package pbr

type surface interface {
	Intersect(Ray3) (bool, float64)
	NormalAt(Vector3) Vector3
	MaterialAt(Vector3) Material
}
