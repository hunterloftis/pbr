package trace

// Hit describes an intersection
type Hit struct {
	Normal Vector3
	Mat    Material
	Dist   float64
	Point  Vector3
}
