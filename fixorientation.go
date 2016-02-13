package geom

// orientation2D_Polygon(): test the orientation of a simple 2D polygon
//  Input:  Point* V = an array of n+1 vertex points with V[n]=V[0]
//  Return: >0 for counterclockwise
//          =0 for none (degenerate)
//          <0 for clockwise
//  Note: this algorithm is faster than computing the signed area.
//  From http://geomalgorithms.com/a01-_area.html#orientation2D_Polygon()
func orientation(V Polygon) []float64 {
	// first find rightmost lowest vertex of the polygon
	out := make([]float64, len(V))
	for j, r := range V {
		rmin := 0
		xmin := r[0].X
		ymin := r[0].Y
		for i, p := range r {
			if p.Y > ymin {
				continue
			} else if p.Y == ymin { // just as low
				if p.X < xmin { // and to left
					continue
				}
			}
			rmin = i // a new rightmost lowest vertex
			xmin = p.X
			ymin = p.Y
		}

		// test orientation at the rmin vertex
		// ccw <=> the edge leaving V[rmin] is left of the entering edge
		if rmin == 0 || rmin == len(r)-1 {
			out[j] = isLeft(r[len(r)-2], r[0], r[1])
		} else {
			out[j] = isLeft(r[rmin-1], r[rmin], r[rmin+1])
		}
	}
	return out
}

// FixOrientation changes the winding direction of the outer and inner
// rings of p so the outer ring is counter-clockwise and
// nesting rings alternate directions.
func (p Polygon) FixOrientation() {
	o := orientation(p)
	for i, inner := range p {
		numInside := 0
		for j, outer := range p {
			if i != j {
				if polyInPoly(outer, inner) {
					numInside++
				}
			}
		}
		if numInside%2 == 1 && o[i] > 0. {
			reversePolygon(inner)
		} else if numInside%2 == 0 && o[i] < 0. {
			reversePolygon(inner)
		}
	}
}

// FixOrientation fixes the orientation of the polygons that comprise mp.
func (mp MultiPolygon) FixOrientation() {
	for _, p := range mp {
		p.FixOrientation()
	}
}

func reversePolygon(s []Point) []Point {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
