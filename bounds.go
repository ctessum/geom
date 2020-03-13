package geom

import (
	"math"
)

// Bounds holds the spatial extent of a geometry.
type Bounds struct {
	Min, Max Point
}

// Extend increases the extent of b1 to include b2.
func (b *Bounds) Extend(b2 *Bounds) {
	if b2 == nil {
		return
	}
	b.extendPoint(b2.Min)
	b.extendPoint(b2.Max)
}

// NewBounds initializes a new bounds object.
func NewBounds() *Bounds {
	return &Bounds{Point{X: math.Inf(1), Y: math.Inf(1)}, Point{X: math.Inf(-1), Y: math.Inf(-1)}}
}

// NewBoundsPoint creates a bounds object from a point.
func NewBoundsPoint(point Point) *Bounds {
	return &Bounds{Point{X: point.X, Y: point.Y}, Point{X: point.X, Y: point.Y}}
}

// Copy returns a copy of b.
func (b *Bounds) Copy() *Bounds {
	return &Bounds{Point{X: b.Min.X, Y: b.Min.Y}, Point{X: b.Max.X, Y: b.Max.Y}}
}

// Empty returns true if b does not contain any points.
func (b *Bounds) Empty() bool {
	return b.Max.X < b.Min.X || b.Max.Y < b.Min.Y
}

func (b *Bounds) extendPoint(point Point) *Bounds {
	b.Min.X = math.Min(b.Min.X, point.X)
	b.Min.Y = math.Min(b.Min.Y, point.Y)
	b.Max.X = math.Max(b.Max.X, point.X)
	b.Max.Y = math.Max(b.Max.Y, point.Y)
	return b
}

func (b *Bounds) extendPoints(points []Point) {
	for _, point := range points {
		b.extendPoint(point)
	}
}

func (b *Bounds) extendPointss(pointss []Path) {
	for _, points := range pointss {
		b.extendPoints(points)
	}
}

// Overlaps returns whether b and b2 overlap.
func (b *Bounds) Overlaps(b2 *Bounds) bool {
	return b.Min.X <= b2.Max.X && b.Min.Y <= b2.Max.Y && b.Max.X >= b2.Min.X && b.Max.Y >= b2.Min.Y
}

// Bounds returns b
func (b *Bounds) Bounds() *Bounds {
	return b
}

// Within calculates whether b is within poly.
func (b *Bounds) Within(poly Polygonal) WithinStatus {
	if bp, ok := poly.(*Bounds); ok {
		if b.Min.Equals(bp.Min) && b.Max.Equals(bp.Max) {
			return OnEdge
		} else if b.Min.X >= bp.Min.X && b.Min.Y >= bp.Min.Y && b.Max.X <= bp.Max.X && b.Max.Y <= bp.Max.Y {
			return Inside
		}
		return Outside
	}
	minIn := pointInPolygonal(b.Min, poly)
	maxIn := pointInPolygonal(b.Max, poly)
	if minIn == Outside || maxIn == Outside {
		return Outside
	}
	return Inside
}

// Len returns the number of points in the receiver (always==5).
func (b *Bounds) Len() int { return 4 }

// Points returns an iterator for the corners of the receiver.
func (b *Bounds) Points() func() Point {
	var i int
	return func() Point {
		defer func() {
			i++
		}()
		switch i {
		case 0:
			return b.Min
		case 1:
			return Point{b.Max.X, b.Min.Y}
		case 2:
			return b.Max
		case 3:
			return Point{b.Min.X, b.Max.Y}
		default:
			panic("out of bounds")
		}
	}
}

// Polygons returns a rectangle polygon
// to fulfill the Polygonal interface.
func (b *Bounds) Polygons() []Polygon {
	return []Polygon{{{b.Min, Point{b.Max.X, b.Min.Y}, b.Max, Point{b.Min.X, b.Max.Y}}}}
}

// Intersection returns the Intersection of the receiver with p.
func (b *Bounds) Intersection(p Polygonal) Polygonal {
	if bp, ok := p.(*Bounds); ok {
		// Special case, other polygon is *Bounds.
		i := &Bounds{
			Min: Point{X: math.Max(b.Min.X, bp.Min.X), Y: math.Max(b.Min.Y, bp.Min.Y)},
			Max: Point{X: math.Min(b.Max.X, bp.Max.X), Y: math.Min(b.Max.Y, bp.Max.Y)},
		}
		if i.Min.Equals(i.Max) {
			return nil
		}
		return i
	}

	bp := p.Bounds()
	if w := bp.Within(b); w == Inside || w == OnEdge {
		// Polygon fully within bounds.
		return p
	} else if bbp, ok := p.(*Bounds); ok {
		// Polygon is bounds.
		if w := b.Within(bbp); w == Inside || w == OnEdge {
			return b
		}
	}
	if !b.Overlaps(bp) {
		return nil
	}
	return b.Polygons()[0].Intersection(p)
}

// Union returns the combination of the receiver and p.
func (b *Bounds) Union(p Polygonal) Polygonal {
	// TODO: optimize
	return b.Polygons()[0].Union(p)
}

// XOr returns the area(s) occupied by either the receiver or p but not both.
func (b *Bounds) XOr(p Polygonal) Polygonal {
	// TODO: optimize
	return b.Polygons()[0].XOr(p)
}

// Difference subtracts p from b.
func (b *Bounds) Difference(p Polygonal) Polygonal {
	// TODO: optimize
	return b.Polygons()[0].Difference(p)
}

// Area returns the area of the reciever.
func (b *Bounds) Area() float64 {
	return (b.Max.X - b.Min.X) * (b.Max.Y - b.Min.Y)
}

// Simplify returns the receiver
// to fulfill the Polygonal interface.
func (b *Bounds) Simplify(tolerance float64) Geom {
	return b
}

func (b *Bounds) Centroid() Point {
	return Point{(b.Min.X + b.Max.X) / 2, (b.Min.Y + b.Max.Y) / 2}
}
