package pbr

// Hit describes an intersection between a Ray3 and a Surface
type Hit struct {
	Normal   Vector3   // The Normal of the Surface at the intersection point
	Incident Vector3   // The incident direction
	Dist     float64   // The Distance along the Ray3 where the intersection occurs
	Mat      *Material // The Material of the Surface at the intersection point
	Point    Vector3   // The point of intersection in world space
}
