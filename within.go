package geom

import "math"

// pointInPolygonal determines whether "pt" is
// within any of the polygons in "pg".
// adapted from https://rosettacode.org/wiki/Ray-casting_algorithm#Go.
// In this version of the algorithm, points that lie on the edge of the polygon
// are considered inside.
func pointInPolygonal(pt Point, pg Polygonal) (in bool) {
	for _, poly := range pg.Polygons() {
		for _, ring := range poly {
			if len(ring) < 3 {
				return false
			}
			// check segment between beginning and ending points
			if !ring[len(ring)-1].Equals(ring[0]) {
				if distPointToSegment(pt, ring[len(ring)-1], ring[0]) == 0 {
					return true
				}
				if rayIntersectsSegment(pt, ring[len(ring)-1], ring[0]) {
					in = !in
				}
			}
			// check the rest of the segments.
			for i := 1; i < len(ring); i++ {
				if distPointToSegment(pt, ring[i-1], ring[i]) == 0 {
					return true
				}
				if rayIntersectsSegment(pt, ring[i-1], ring[i]) {
					in = !in
				}
			}
		}
	}
	return in
}

func rayIntersectsSegment(p, a, b Point) bool {
	if a.Y > b.Y {
		a, b = b, a
	}
	for p.Y == a.Y || p.Y == b.Y {
		p.Y = math.Nextafter(p.Y, math.Inf(1))
	}
	if p.Y < a.Y || p.Y > b.Y {
		return false
	}
	if a.X > b.X {
		if p.X >= a.X {
			return false
		}
		if p.X < b.X {
			return true
		}
	} else {
		if p.X > b.X {
			return false
		}
		if p.X < a.X {
			return true
		}
	}
	return (p.Y-a.Y)/(p.X-a.X) >= (b.Y-a.Y)/(b.X-a.X)
}
