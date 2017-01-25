package martinez

import "math"

// A Segment defines a line segment connecting two Points.
type Segment interface {
	Start() Point
	End() Point
}

// A Point defines a location in 2-dimensional space.
type Point interface {
	X() float64
	Y() float64
}

type point struct{ x, y float64 }

// NewPoint creates a new Point at location [X,Y].
func NewPoint(X, Y float64) Point {
	return point{x: X, y: Y}
}

func (p point) X() float64 { return p.x }
func (p point) Y() float64 { return p.y }

// Bounds defines a bounding box in 2-dimensional space.
type Bounds interface {
	// Min returns the minimum corner of the bounding box.
	Min() Point

	// Max returns the maximum corner of the bounding box.
	Max() Point

	// Extend extends the receiver to include the
	// specified Point.
	Extend(Point)
}

type bounds struct {
	min, max point
}

// NewBounds initializes a new Bounds variable.
func NewBounds() Bounds {
	return &bounds{
		min: point{x: math.Inf(1), y: math.Inf(1)},
		max: point{x: math.Inf(-1), y: math.Inf(-1)},
	}
}

// Min returns the minimum point to fulfil the Bounds interface.
func (b *bounds) Min() Point {
	return b.min
}

// Min returns the maximum point to fulfil the Bounds interface.
func (b *bounds) Max() Point {
	return b.max
}

// Extend extends the receiver to include point p to
// fulfil the Bounds interface.
func (b *bounds) Extend(p Point) {
	b.min.x = math.Min(b.min.x, p.X())
	b.min.y = math.Min(b.min.y, p.Y())
	b.max.x = math.Max(b.max.x, p.X())
	b.max.y = math.Max(b.max.y, p.Y())
}

// A Path represents a connected chain of points.
type Path interface {
	Len() int
	At(int) Point
}

type path []Point

func (p path) Len() int       { return len(p) }
func (p path) At(i int) Point { return p[i] }

// NewPath creates a new Path from the specified points.
func NewPath(points ...Point) Path { return path(points) }

// A MultiPath represents a set of related paths.
type MultiPath interface {
	Len() int
	At(int) Path
}

type multiPath []Path

func (p multiPath) Len() int      { return len(p) }
func (p multiPath) At(i int) Path { return p[i] }

// A Polygon is series of related closed paths (i.e., paths where
// the beginning and ending Points are considered to be connected).
type Polygon MultiPath

// NewPolygon creates a new Polygon from the specified paths.
func NewPolygon(paths ...Path) Polygon {
	return multiPath(paths)
}

// A LineString is an open path (i.e., a path where the beginning
// and ending Points are not considered to be connected).
type LineString Path

// NewLineString creates a new LineString from the specified Points.
func NewLineString(points ...Point) LineString { return path(points) }
