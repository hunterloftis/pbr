// Package pbr renders physically-based 3D scenes with a Monte Carlo path tracer.
//
// Basics
//
// A Surface describes surfaces (like spheres, cubes, and triangles).
// Surfaces can be created programmatically or loaded from .obj files.
// A Scene contains various Surfaces.
// A Camera specifies a point-of-view for rendering a Scene.
// A Renderer samples light from the camera into 2D images.
package pbr
