package collada

// Triangle describes a 3D triangle's position, normal, and material.
type Triangle struct {
	Pos  [3]Vector3
	Norm [3]Vector3
	Mat  *Material
}