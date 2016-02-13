package geom

// Point is a holder for 2D coordinates X and Y.
type Point struct {
	X, Y float64
}

// NewPoint returns a new point with coordinates x and y.
func NewPoint(x, y float64) *Point {
	return &Point{X: x, Y: y}
}

// Bounds gives the rectangular extents of the Point.
func (p Point) Bounds() *Bounds {
	return NewBoundsPoint(p)
}

// Within calculates whether p is within poly or touching its edge.
func (p Point) Within(poly Polygonal) bool {
	return pointInPolygonal(p, poly)
}
