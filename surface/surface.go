package surface

import (
	"github.com/hunterloftis/pbr/geom"
	"github.com/hunterloftis/pbr/surface/material"
)

// Surface is an interface to all surface types (sphere, cube, etc).
// A surface is anything that can be intersected by a Ray.
type Surface interface {
	Intersect(*geom.Ray3) Hit
	At(point geom.Vector3) (normal geom.Direction, mat *material.Material)
	Box() *Box
	Center() geom.Vector3
	Material() *material.Material
}
