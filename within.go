package geom

import "math"

func polyInPoly(outer, inner []Point) bool {
	for _, p := range inner {
		if pointInPoly(p, outer) == 0 {
			return false
		}
	}
	return true
}

func pointInPolygon(point Point, polygon Polygon) bool {
	inCount := 0
	o := orientation(polygon)
	for i, r := range polygon {
		if pointInPoly(point, r) != 0 {
			if o[i] > 0. {
				inCount++
			} else if o[i] < 0. {
				inCount--
			}
		}
	}
	return inCount > 0
}

// Function pointInPolygonal determines whether "point" is
// within any of the polygons in "polygonal". Also returns true if "point" is on the
// edge of "polygon".
func pointInPolygonal(point Point, polygonal Polygonal) bool {
	for _, poly := range polygonal.Polygons() {
		if pointInPolygon(point, poly) {
			return true
		}
	}
	return false
}

//returns 0 if false, +1 if true, -1 if pt ON polygon boundary
//See "The Point in Polygon Problem for Arbitrary Polygons" by Hormann & Agathos
//http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.88.5498&rep=rep1&type=pdf
func pointInPoly(pt Point, path []Point) int {
	result := 0
	cnt := len(path)
	if cnt < 3 {
		return 0
	}
	ip := path[0]
	for i := 1; i <= cnt; i++ {
		var ipNext Point
		if i == cnt {
			ipNext = path[0]
		} else {
			ipNext = path[i]
		}
		if ipNext.Y == pt.Y {
			if ipNext.X == pt.X || (ip.Y == pt.Y &&
				((ipNext.X-pt.X > 0) == (ip.X-pt.X < 0))) {
				return -1
			}
		}
		if (ip.Y-pt.Y < 0) != (ipNext.Y-pt.Y < 0) {
			if ip.X-pt.X >= 0 {
				if ipNext.X-pt.X > 0 {
					result = 1 - result
				} else {
					d := (ip.X-pt.X)*(ipNext.Y-pt.Y) -
						(ipNext.X-pt.X)*(ip.Y-pt.Y)
					if d == 0 {
						return -1
					} else if (d > 0) == (ipNext.Y-ip.Y > 0) {
						result = 1 - result
					}
				}
			} else {
				if ipNext.X-pt.X > 0 {
					d := (ip.X-pt.X)*(ipNext.Y-pt.Y) -
						(ipNext.X-pt.X)*(ip.Y-pt.Y)
					if d == 0 {
						return -1
					} else if (d > 0) == (ipNext.Y-ip.Y > 0) {
						result = 1 - result
					}
				}
			}
		}
		ip = ipNext
	}
	return result
}

// dot product
func dot(u, v Point) float64 { return u.X*v.X + u.Y*v.Y }

// norm = length of  vector
func norm(v Point) float64 { return math.Sqrt(dot(v, v)) }

// distance = norm of difference
func d(u, v Point) float64 { return norm(pointSubtract(u, v)) }

// dist_Point_to_Segment(): get the distance of a point to a segment
//     Input:  a Point P and a Segment S (in any dimension)
//     Return: the shortest distance from P to S
// from http://geomalgorithms.com/a02-_lines.html
func distPointToSegment(p, segStart, segEnd Point) float64 {
	v := pointSubtract(segEnd, segStart)
	w := pointSubtract(p, segStart)

	c1 := dot(w, v)
	if c1 <= 0. {
		return d(p, segStart)
	}

	c2 := dot(v, v)
	if c2 <= c1 {
		return d(p, segEnd)
	}

	b := c1 / c2
	pb := Point{segStart.X + b*v.X, segStart.Y + b*v.Y}
	return d(p, pb)
}

func pointSubtract(p1, p2 Point) Point {
	return Point{X: p1.X - p2.X, Y: p1.Y - p2.Y}
}

func pointOnSegment(p, segStart, segEnd Point, tolerance float64) bool {
	return distPointToSegment(p, segStart, segEnd) < tolerance
}

// isLeft(): test if a point is Left|On|Right of an infinite 2D line.
//    Input:  three points P0, P1, and P2
//    Return: >0 for P2 left of the line through P0 to P1
//          =0 for P2 on the line
//          <0 for P2 right of the line
//    From http://geomalgorithms.com/a01-_area.html#isLeft()
func isLeft(P0, P1, P2 Point) float64 {
	return ((P1.X-P0.X)*(P2.Y-P0.Y) -
		(P2.X-P0.X)*(P1.Y-P0.Y))
}
