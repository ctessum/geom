package geom

import (
	"math"

	"github.com/ctessum/polyclip-go"
)

// A Polygon is a series of closed rings. The inner rings should be nested
// inside of the outer ring.
type Polygon [][]Point

// Bounds gives the rectangular extents of the polygon.
func (p Polygon) Bounds() *Bounds {
	b := NewBounds()
	b.extendPointss(p)
	return b
}

// Polygons returns []{p} to fulfill the Polygonal interface.
func (p Polygon) Polygons() []Polygon {
	return []Polygon{p}
}

// Intersection returns the area(s) shared by p and p2.
func (p Polygon) Intersection(p2 Polygonal) Polygon {
	return p.op(p2, polyclip.INTERSECTION)
}

// Union returns the combination of p and p2.
func (p Polygon) Union(p2 Polygonal) Polygon {
	return p.op(p2, polyclip.UNION)
}

// XOr returns the area(s) occupied by either p or p2 but not both.
func (p Polygon) XOr(p2 Polygonal) Polygon {
	return p.op(p2, polyclip.XOR)
}

// Difference subtracts p2 from p.
func (p Polygon) Difference(p2 Polygonal) Polygon {
	return p.op(p2, polyclip.DIFFERENCE)
}

func (p Polygon) op(p2 Polygonal, op polyclip.Op) Polygon {
	pp := p.toPolyClip()
	var pp2 polyclip.Polygon
	for _, pp2x := range p2.Polygons() {
		pp2 = append(pp2, pp2x.toPolyClip()...)
	}
	return polyClipToPolygon(pp.Construct(op, pp2))
}

func (p Polygon) toPolyClip() polyclip.Polygon {
	o := make(polyclip.Polygon, len(p))
	for i, r := range p {
		o[i] = make(polyclip.Contour, len(r))
		for j, pp := range r {
			o[i][j] = polyclip.Point(pp)
		}
	}
	return o
}

func polyClipToPolygon(p polyclip.Polygon) Polygon {
	pp := make(Polygon, len(p))
	for i, r := range p {
		pp[i] = make([]Point, len(r)+1)
		for j, ppp := range r {
			pp[i][j] = Point(ppp)
		}
		// Close the ring as per OGC standard.
		pp[i][len(r)] = pp[i][0]
	}
	return pp
}

// Area returns the area of p, assuming
// that nested polygons have alternating winding directions.
func (p Polygon) Area() float64 {
	a := 0.
	for _, r := range p {
		a += area(r)
	}
	return math.Abs(a)
}

// see http://www.mathopenref.com/coordpolygonarea2.html
func area(polygon []Point) float64 {
	highI := len(polygon) - 1
	A := (polygon[highI].X +
		polygon[0].X) * (polygon[0].Y - polygon[highI].Y)
	for i := 0; i < highI; i++ {
		A += (polygon[i].X +
			polygon[i+1].X) * (polygon[i+1].Y - polygon[i].Y)
	}
	return A / 2.
}

// Centroid calculates the centroid of p, from
// wikipedia: http://en.wikipedia.org/wiki/Centroid#Centroid_of_polygon.
// The polygon can have holes, but each ring must be closed (i.e.,
// p[0] == p[n-1], where the ring has n points) and must not be
// self-intersecting.
// The algorithm will not check to make sure the holes are
// actually inside the outer rings.
func (p Polygon) Centroid() Point {
	var A, xA, yA float64
	for _, r := range p {
		a := area(r)
		cx, cy := 0., 0.
		for i := 0; i < len(r)-1; i++ {
			cx += (r[i].X + r[i+1].X) *
				(r[i].X*r[i+1].Y - r[i+1].X*r[i].Y)
			cy += (r[i].Y + r[i+1].Y) *
				(r[i].X*r[i+1].Y - r[i+1].X*r[i].Y)
		}
		cx /= 6 * a
		cy /= 6 * a
		A += a
		xA += cx * a
		yA += cy * a
	}
	return Point{X: xA / A, Y: yA / A}
}

// Within calculates whether p is completely within p2. Edges that touch are
// considered to be within. It may not work correctly if p has holes.
func (p Polygon) Within(p2 Polygonal) bool {
	for _, r := range p {
		for _, pp := range r {
			if !pointInPolygonal(pp, p2) {
				return false
			}
		}
	}
	return true
}
