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
func (point Point) Bounds() *Bounds {
	return NewBoundsPoint(point)
}
