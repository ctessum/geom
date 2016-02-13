package geom

import "math"

// MultiPolygon is a holder for multiple related polygons.
type MultiPolygon []Polygon

// Bounds gives the rectangular extents of the MultiPolygon.
func (mp MultiPolygon) Bounds() *Bounds {
	b := NewBounds()
	for _, polygon := range mp {
		b.Extend(polygon.Bounds())
	}
	return b
}

// Area returns the combined area of the polygons in p,
//  assuming that none of the polygons in the p
// overlap and that nested polygons have alternating winding directions.
func (mp MultiPolygon) Area() float64 {
	a := 0.
	for _, pp := range mp {
		a += pp.Area()
	}
	return math.Abs(a)
}

func (mp MultiPolygon) Intersection(p2 Polygonal) Polygonal {
	var o MultiPolygon
	for _, pp := range p2.Polygons() {
		o = append(o, pp.Intersection(mp).Polygons()...)
	}
	if len(o) == 0 {
		return nil
	} else if len(o) == 1 {
		return o[0]
	}
	return o
}

// Polygons returns the polygons that make up mp.
func (mp MultiPolygon) Polygons() []Polygon {
	return mp
}
