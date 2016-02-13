/*
Package geom holds geometry objects and functions to operate on them.
They can be encoded and decoded by other packages in this repository.
*/
package geom

// T is an interface for generic geometry types.
type T interface {
	Bounds() *Bounds
}

// Linear is an interface for types that are linear in nature.
type Linear interface {
	T
	Length() float64
	Simplify(tolerance float64) Polygonal
}

// Polygonal is an interface for types that are polygonal in nature.
type Polygonal interface {
	T
	Polygons() []Polygon
	Intersection(Polygonal) Polygonal
	Area() float64
	Simplify(tolerance float64) Polygonal
}

// PointLike is an interface for types that are pointlike in nature.
type PointLike interface {
	T
	Points() []Point
	Within(Polygonal) bool
	On(l Linear, tolerance float64) bool
}
