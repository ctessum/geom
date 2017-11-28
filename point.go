package geom

import (
	"fmt"
	"math"
)

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

// Within calculates whether p is within poly.
func (p Point) Within(poly Polygonal) WithinStatus {
	return pointInPolygonal(p, poly)
}

// Equals returns whether p is equal to p2.
func (p Point) Equals(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

// Buffer returns a circle with the specified radius
// centered at the receiver location. The circle is represented
// as a polygon with the specified number of segments.
func (p Point) Buffer(radius float64, segments int) Polygon {
	if segments < 3 {
		panic(fmt.Errorf("geom: invalid number of segments %d", segments))
	}
	if radius < 0 {
		panic(fmt.Errorf("geom: invalid radius %g", radius))
	}
	dTheta := math.Pi * 2 / float64(segments)
	o := make(Polygon, 1)
	o[0] = make([]Point, segments)
	for i := 0; i < segments; i++ {
		theta := float64(i) * dTheta
		o[0][i] = Point{
			X: p.X + radius*math.Cos(theta),
			Y: p.Y + radius*math.Sin(theta),
		}
	}
	return o
}

// Len returns the number of points in the receiver (always==1).
func (p Point) Len() int { return 1 }

// Points returns an iterator for the points in the receiver (there will only
// be one point).
func (p Point) Points() func() Point {
	return func() Point { return p }
}
