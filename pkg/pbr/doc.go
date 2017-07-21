// Package pbr implements Physically-Based Rendering with a Monte Carlo path tracer.
//
// Surface describes surfaces like spheres and cubes.
// Scene contains various Surfaces.
// Camera specifies a point-of-view for rendering a Scene.
// Sampler samples light energy in a Scene from a particular Camera.
// Renderer renders sampled light onto Image pixels.
package pbr

// TODO: reduce public API to only what's really necessary
