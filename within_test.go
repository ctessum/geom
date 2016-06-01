package geom

import "testing"

func TestWithin1(t *testing.T) {
	p := Point{620858.7034230313, -1.3334340701764394e+06}
	b := Polygon{
		[]Point{
			{-2.758081092115788e+06, -2.1035219712004187e+06},
			{-2.7580810921157864e+06, 1.9603377468041454e+06},
			{2.6080741578387334e+06, 1.954523927465083e+06},
			{2.60226033849967e+06, -2.10352197120042e+06},
			{-2.758081092115788e+06, -2.1035219712004187e+06},
		},
	}
	if p.Within(b) != Inside {
		t.Errorf("Point %v should be within polygon %v", p, b)
	}
}

// Adapted from https://rosettacode.org/wiki/Ray-casting_algorithm#Go
func TestWithin2(t *testing.T) {
	var (
		p1  = Point{0, 0}
		p2  = Point{10, 0}
		p3  = Point{10, 10}
		p4  = Point{0, 10}
		p5  = Point{2.5, 2.5}
		p6  = Point{7.5, 2.5}
		p7  = Point{7.5, 7.5}
		p8  = Point{2.5, 7.5}
		p9  = Point{0, 5}
		p10 = Point{10, 5}
		p11 = Point{3, 0}
		p12 = Point{7, 0}
		p13 = Point{7, 10}
		p14 = Point{3, 10}
	)
	type poly struct {
		name    string
		sides   Polygon
		results []WithinStatus
	}

	var tpg = []poly{
		poly{
			name:    "square",
			sides:   Polygon{[]Point{p1, p2, p3, p4, p1}},
			results: []WithinStatus{Inside, Inside, Outside, OnEdge, OnEdge, Inside, OnEdge, Inside, Inside},
		},
		poly{
			name: "square hole",
			sides: Polygon{
				[]Point{p1, p2, p3, p4, p1},
				[]Point{p5, p6, p7, p8, p5},
			},
			results: []WithinStatus{Outside, Inside, Outside, OnEdge, OnEdge, Inside, OnEdge, Inside, Inside},
		},
		poly{
			name:    "strange",
			sides:   Polygon{[]Point{p1, p5, p4, p8, p7, p3, p2, p5}},
			results: []WithinStatus{Inside, Outside, Outside, Outside, OnEdge, Inside, OnEdge, Outside, Outside},
		},
		poly{
			name:    "exagon",
			sides:   Polygon{[]Point{p11, p12, p10, p13, p14, p9, p11}},
			results: []WithinStatus{Inside, Inside, Outside, OnEdge, OnEdge, Inside, Outside, Outside, Outside},
		},
	}

	var tpt = []Point{{5, 5}, {5, 8}, {-10, 5}, {0, 5}, {10, 5}, {8, 5}, {10, 10}, {1, 2}, {2, 1}}

	for _, pg := range tpg {
		for i, pt := range tpt {
			if pg.results[i] != pt.Within(pg.sides) {
				t.Errorf("point %v within polygon %v should be %v", pt, pg.name, pg.results[i])
			}
		}
	}
}

func TestWithin3(t *testing.T) {
	p := Point{X: 2, Y: 2}
	poly := Polygon{{
		Point{X: 1, Y: 0},
		Point{X: 2, Y: 2},
		Point{X: 0, Y: 2},
	}}
	if p.Within(poly) != OnEdge {
		t.Errorf("%v should be on edge of %v", p, poly)
	}
}

func TestWithin4(t *testing.T) {
	p := Point{X: 1, Y: 0.5}
	poly := Polygon{{
		Point{X: 0, Y: 1},
		Point{X: 1, Y: 0},
		Point{X: 2, Y: 1},
		Point{X: 1, Y: 2},
	}}
	if p.Within(poly) != Inside {
		t.Errorf("%v should be within %v", p, poly)
	}
}
