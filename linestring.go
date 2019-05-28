package geom

import (
	"math"

	polyclip "github.com/ctessum/polyclip-go"
)

// LineString is a number of points that make up a path or line.
type LineString Path

// Bounds gives the rectangular extents of the LineString.
func (l LineString) Bounds() *Bounds {
	b := NewBounds()
	b.extendPoints(l)
	return b
}

// Length calculates the length of l.
func (l LineString) Length() float64 {
	length := 0.
	for i := 0; i < len(l)-1; i++ {
		p1 := l[i]
		p2 := l[i+1]
		length += math.Hypot(p2.X-p1.X, p2.Y-p1.Y)
	}
	return length
}

// Within calculates whether l is completely within p or touching its edge.
func (l LineString) Within(p Polygonal) WithinStatus {
	for _, pp := range l {
		if pointInPolygonal(pp, p) == Outside {
			return Outside
		}
	}
	return Inside
}

// Distance calculates the shortest distance from p to the LineString.
func (l LineString) Distance(p Point) float64 {
	d := math.Inf(1)
	for i := 0; i < len(l)-1; i++ {
		segDist := distPointToSegment(p, l[i], l[i+1])
		d = math.Min(d, segDist)
	}
	return d
}

// Clip returns the part of the receiver that falls within the given polygon.
func (l LineString) Clip(p Polygonal) Linear {
	pTemp := Polygon{Path(l)}.op(p, polyclip.CLIPLINE)
	o := make(MultiLineString, len(pTemp))
	for i, pp := range pTemp {
		o[i] = LineString(pp[0 : len(pp)-1])
	}
	return o
}

// Len returns the number of points in the receiver.
func (l LineString) Len() int { return len(l) }

// XY returns the coordinates of point at location i in the receiver.
func (l LineString) XY(i int) (x, y float64) {
	p := l[i]
	return p.X, p.Y
}
