// Package pbr renders physically-based 3D scenes with a Monte Carlo path tracer.
//
// Basics
//
// A Surface describes surfaces (like Spheres, Cubes, and Triangles).
// Surfaces can be created programmatically or loaded from .obj files.
// A Scene contains various Surfaces to be rendered together.
// A Sensor specifies a point-of-view for rendering a Scene.
// A Render samples light from the camera to create 2D images.
package pbr
