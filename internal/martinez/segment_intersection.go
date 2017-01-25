package martinez

import "math"

const epsilon = 1e-9

// krossProduct Finds the magnitude of the cross product of two vectors (if we pretend
// they're in three dimensions).
func krossProduct(a, b Point) float64 {
	return a.X()*b.Y() - a.Y()*b.X()
}

// dotProduct returns the dot product of two vectors.
func dotProduct(a, b Point) float64 {
	return a.X()*b.X() + a.Y()*b.Y()
}

/**
 * Finds the intersection (if any) between two line segments a and b, given the
 * line segments' end points a1, a2 and b1, b2.
 *
 * This algorithm is based on Schneider and Eberly.
 * http://www.cimec.org.ar/~ncalvo/Schneider_Eberly.pdf
 * Page 244.
 *
 * @param {Array.<Number>} a1 point of first line
 * @param {Array.<Number>} a2 point of first line
 * @param {Array.<Number>} b1 point of second line
 * @param {Array.<Number>} b2 point of second line
 * @param {Boolean=}       noEndpointTouch whether to skip single touchpoints
 *                                         (meaning connected segments) as
 *                                         intersections
 * @returns {Array.<Array.<Number>>|Null} If the lines intersect, the point of
 * intersection. If they overlap, the two end points of the overlapping segment.
 * Otherwise, null.
 */
func segmentIntersection(a1, a2, b1, b2 Point, noEndpointTouch bool) (i1, i2 Point, num int) {
	// The algorithm expects our lines in the form P + sd, where P is a point,
	// s is on the interval [0, 1], and d is a vector.
	// We are passed two points. P can be the first point of each pair. The
	// vector, then, could be thought of as the distance (in x and y components)
	// from the first point to the second point.
	// So first, let's make our vectors:
	var va = point{x: a2.X() - a1.X(), y: a2.Y() - a1.Y()}
	var vb = point{x: b2.X() - b1.X(), y: b2.Y() - b1.Y()}
	// We also define a function to convert back to regular point form:

	toPoint := func(p Point, s float64, d Point) Point {
		return point{
			x: p.X() + s*d.X(),
			y: p.Y() + s*d.Y(),
		}
	}

	// The rest is pretty much a straight port of the algorithm.
	var e = point{x: b1.X() - a1.X(), y: b1.Y() - a1.Y()}
	var kross = krossProduct(va, vb)
	var sqrKross = kross * kross
	var sqrLenA = dotProduct(va, va)
	var sqrLenB = dotProduct(vb, vb)

	// Check for line intersection. This works because of the properties of the
	// cross product -- specifically, two vectors are parallel if and only if the
	// cross product is the 0 vector. The full calculation involves relative error
	// to account for possible very small line segments. See Schneider & Eberly
	// for details.
	if sqrKross > epsilon*sqrLenA*sqrLenB {
		// If they're not parallel, then (because these are line segments) they
		// still might not actually intersect. This code checks that the
		// intersection point of the lines is actually on both line segments.
		var s = krossProduct(e, vb) / kross
		if s < 0 || s > 1 {
			// not on line segment a
			return point{}, point{}, 0
		}
		var t = krossProduct(e, va) / kross
		if t < 0 || t > 1 {
			// not on line segment b
			return point{}, point{}, 0
		}
		if noEndpointTouch {
			return point{}, point{}, 0
		}
		return toPoint(a1, s, va), point{}, 1
	}

	// If we've reached this point, then the lines are either parallel or the
	// same, but the segments could overlap partially or fully, or not at all.
	// So we need to find the overlap, if any. To do that, we can use e, which is
	// the (vector) difference between the two initial points. If this is parallel
	// with the line itself, then the two lines are the same line, and there will
	// be overlap.
	var sqrLenE = dotProduct(e, e)
	kross = krossProduct(e, va)
	sqrKross = kross * kross

	if sqrKross > epsilon*sqrLenA*sqrLenE {
		// Lines are just parallel, not the same. No overlap.
		return point{}, point{}, 0
	}

	var sa = dotProduct(va, e) / sqrLenA
	var sb = sa + dotProduct(va, vb)/sqrLenA
	var smin = math.Min(sa, sb)
	var smax = math.Max(sa, sb)

	// this is, essentially, the FindIntersection acting on floats from
	// Schneider & Eberly, just inlined into this function.
	if smin <= 1 && smax >= 0 {

		// overlap on an end point
		if smin == 1 {
			if noEndpointTouch {
				return point{}, point{}, 0
			}
			return toPoint(a1, math.Max(smin, 0), va), point{}, 1
		}

		if smax == 0 {
			if noEndpointTouch {
				return point{}, point{}, 0
			}
			return toPoint(a1, math.Min(smax, 1), va), point{}, 1
		}

		if noEndpointTouch && smin == 0 && smax == 1 {
			return point{}, point{}, 0
		}

		// There's overlap on a segment -- two points of intersection. Return both.
		return toPoint(a1, math.Max(smin, 0), va), toPoint(a1, math.Min(smax, 1), va), 2
	}
	return point{}, point{}, 0
}
