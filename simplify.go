package geom

// Simplify simplifies p
// by removing points according to the tolerance parameter,
// while ensuring that the resulting shape is not self intersecting
// (but only if the input shape is not self intersecting). Self-intersecting
// polygons may cause the algorithm to fall into an infinite loop.
//
// It is based on the algorithm:
// J. L. G. Pallero, Robust line simplification on the plane.
// Comput. Geosci. 61, 152–159 (2013).
func (p Polygon) Simplify(tolerance float64) Polygonal {
	var out Polygon = make([][]Point, len(p))
	for i, r := range p {
		out[i] = simplifyCurve(r, p, tolerance)
	}
	return out
}

// Simplify simplifies mp
// by removing points according to the tolerance parameter,
// while ensuring that the resulting shape is not self intersecting
// (but only if the input shape is not self intersecting). Self-intersecting
// polygons may cause the algorithm to fall into an infinite loop.
//
// It is based on the algorithm:
// J. L. G. Pallero, Robust line simplification on the plane.
// Comput. Geosci. 61, 152–159 (2013).
func (mp MultiPolygon) Simplify(tolerance float64) Polygonal {
	out := make(MultiPolygon, len(mp))
	for i, p := range mp {
		out[i] = p.Simplify(tolerance).(Polygon)
	}
	return out
}

// Simplify simplifies l
// by removing points according to the tolerance parameter,
// while ensuring that the resulting shape is not self intersecting
// (but only if the input shape is not self intersecting).
//
// It is based on the algorithm:
// J. L. G. Pallero, Robust line simplification on the plane.
// Comput. Geosci. 61, 152–159 (2013).
func (l LineString) Simplify(tolerance float64) Linear {
	return LineString(simplifyCurve(l, [][]Point{}, tolerance))
}

// Simplify simplifies ml
// by removing points according to the tolerance parameter,
// while ensuring that the resulting shape is not self intersecting
// (but only if the input shape is not self intersecting).
//
// It is based on the algorithm:
// J. L. G. Pallero, Robust line simplification on the plane.
// Comput. Geosci. 61, 152–159 (2013).
func (ml MultiLineString) Simplify(tolerance float64) Linear {
	out := make(MultiLineString, len(ml))
	for i, l := range ml {
		out[i] = l.Simplify(tolerance).(LineString)
	}
	return out
}

func simplifyCurve(curve []Point,
	otherCurves [][]Point, tol float64) []Point {
	out := make([]Point, 0, len(curve))

	i := 0
	for {
		out = append(out, curve[i])
		breakTime := false
		for j := i + 2; j < len(curve); j++ {
			breakTime2 := false
			for k := i + 1; k < j; k++ {
				d := distPointToSegment(curve[k], curve[i], curve[j])
				if d > tol {
					// we have found a candidate point to keep
					for {
						// Make sure this simplifcation doesn't cause any self
						// intersections.
						if j > i+2 &&
							(segMakesNotSimple(curve[i], curve[j-1], [][]Point{out[0:i]}) ||
								segMakesNotSimple(curve[i], curve[j-1], [][]Point{curve[j:]}) ||
								segMakesNotSimple(curve[i], curve[j-1], otherCurves)) {
							j--
						} else {
							i = j - 1
							out = append(out, curve[i])
							breakTime2 = true
							break
						}
					}
				}
				if breakTime2 {
					break
				}
			}
			if j == len(curve)-1 {
				out = append(out, curve[j])
				breakTime = true
			}
		}
		if breakTime {
			break
		}
	}
	return out
}

func segMakesNotSimple(segStart, segEnd Point, paths [][]Point) bool {
	seg1 := segment{segStart, segEnd}
	for _, p := range paths {
		for i := 0; i < len(p)-1; i++ {
			seg2 := segment{p[i], p[i+1]}
			if seg1.start == seg2.start || seg1.end == seg2.end ||
				seg1.start == seg2.end || seg1.end == seg2.start {
				// colocated endpoints are not a problem here
				return false
			}
			numIntersections, _, _ := findIntersection(seg1, seg2)
			if numIntersections > 0 {
				return true
			}
		}
	}
	return false
}
