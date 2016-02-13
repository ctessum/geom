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
func (p Polygon) Intersection(p2 Polygonal) Polygonal {
	return p.op(p2, polyclip.INTERSECTION)
}

func (p Polygon) op(p2 Polygonal, op polyclip.Op) Polygonal {
	var o MultiPolygon
	pp := p.toPolyClip()
	for _, pp2 := range p2.Polygons() {
		oTemp := pp.Construct(op, pp2.toPolyClip())
		if len(oTemp) > 0 {
			o = append(o, polyClipToPolygon(oTemp))
		}
	}
	if len(o) == 1 {
		return o[0]
	}
	return o
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
		pp[i] = make([]Point, len(r))
		for j, ppp := range r {
			pp[i][j] = Point(ppp)
		}
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
