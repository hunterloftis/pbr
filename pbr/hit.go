package pbr

// Hit describes an intersection between a Ray3 and a Surface
type Hit struct {
	Normal Vector3  // The Normal of the Surface at the intersection point
	Mat    Material // The Material of the Surface at the intersection point
	Dist   float64  // The Distance along the Ray3 where the intersection occurs
	Point  Vector3  // The point of intersection in world space
}
