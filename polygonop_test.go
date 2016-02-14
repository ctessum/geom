package geom

import "testing"

func TestPolygonOp(t *testing.T) {
	tests := []struct {
		subject, clipping                    Polygonal
		intersection, union, difference, xor Polygonal
	}{
		{
			subject: Polygon{{
				Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
				Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
				Point{X: 0, Y: 0}}},
			clipping: Polygon{{
				Point{X: -1, Y: -1}, Point{X: 1, Y: -1},
				Point{X: 1, Y: 1}, Point{X: -1, Y: 1},
				Point{X: -1, Y: -1}}},
			intersection: Polygon{{
				Point{X: 0, Y: 0}, Point{X: 1, Y: 0},
				Point{X: 1, Y: 1}, Point{X: 0, Y: 1},
				Point{X: 0, Y: 0}}},
			union: Polygon{{
				Point{X: -1, Y: -1}, Point{X: 1, Y: -1},
				Point{X: 1, Y: 0}, Point{X: 2, Y: 0},
				Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
				Point{X: 0, Y: 1}, Point{X: -1, Y: 1},
				Point{X: -1, Y: -1}}},
			difference: Polygon{{
				Point{X: 1, Y: 0}, Point{X: 2, Y: 0},
				Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
				Point{X: 0, Y: 1}, Point{X: 1, Y: 1},
				Point{X: 1, Y: 0}}},
			xor: Polygon{
				{
					Point{X: 1, Y: 0}, Point{X: 2, Y: 0},
					Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
					Point{X: 0, Y: 1}, Point{X: 1, Y: 1},
					Point{X: 1, Y: 0},
				},
				{
					Point{X: -1, Y: -1}, Point{X: 1, Y: -1},
					Point{X: 1, Y: 0}, Point{X: 0, Y: 0},
					Point{X: 0, Y: 1}, Point{X: -1, Y: 1},
					Point{X: -1, Y: -1},
				},
			},
		},
		{
			subject: MultiPolygon{
				{{
					Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
					Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
					Point{X: 0, Y: 0}}},
				{{
					Point{X: 4, Y: 0}, Point{X: 6, Y: 0},
					Point{X: 6, Y: 2}, Point{X: 4, Y: 2},
					Point{X: 4, Y: 0}}},
			},
			clipping: Polygon{{
				Point{X: 1, Y: 1}, Point{X: 5, Y: 1},
				Point{X: 5, Y: 3}, Point{X: 1, Y: 3},
				Point{X: 1, Y: 1}}},
			intersection: Polygon{
				{
					Point{X: 1, Y: 1}, Point{X: 2, Y: 1},
					Point{X: 2, Y: 2}, Point{X: 1, Y: 2},
					Point{X: 1, Y: 1},
				},
				{
					Point{X: 4, Y: 1}, Point{X: 5, Y: 1},
					Point{X: 5, Y: 2}, Point{X: 4, Y: 2},
					Point{X: 4, Y: 1},
				},
			},
			union: Polygon{{
				Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
				Point{X: 2, Y: 1}, Point{X: 4, Y: 1},
				Point{X: 4, Y: 0}, Point{X: 6, Y: 0},
				Point{X: 6, Y: 2}, Point{X: 5, Y: 2},
				Point{X: 5, Y: 3}, Point{X: 1, Y: 3},
				Point{X: 1, Y: 2}, Point{X: 0, Y: 2},
				Point{X: 0, Y: 0}}},
			difference: Polygon{
				{
					Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
					Point{X: 2, Y: 1}, Point{X: 1, Y: 1},
					Point{X: 1, Y: 2}, Point{X: 0, Y: 2},
					Point{X: 0, Y: 0},
				},
				{
					Point{X: 4, Y: 0}, Point{X: 6, Y: 0},
					Point{X: 6, Y: 2}, Point{X: 5, Y: 2},
					Point{X: 5, Y: 1}, Point{X: 4, Y: 1},
					Point{X: 4, Y: 0},
				}},
			xor: Polygon{
				{
					Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
					Point{X: 2, Y: 1}, Point{X: 1, Y: 1},
					Point{X: 1, Y: 2}, Point{X: 0, Y: 2},
					Point{X: 0, Y: 0},
				},
				{
					Point{X: 4, Y: 0}, Point{X: 6, Y: 0},
					Point{X: 6, Y: 2}, Point{X: 5, Y: 2},
					Point{X: 5, Y: 1}, Point{X: 4, Y: 1},
					Point{X: 4, Y: 0},
				},
				{
					Point{X: 2, Y: 1}, Point{X: 4, Y: 1},
					Point{X: 4, Y: 2}, Point{X: 5, Y: 2},
					Point{X: 5, Y: 3}, Point{X: 1, Y: 3},
					Point{X: 1, Y: 2}, Point{X: 2, Y: 2},
					Point{X: 2, Y: 1},
				}},
		},
	}
	for i, test := range tests {
		intersection := test.subject.Intersection(test.clipping)
		if !intersection.Similar(test.intersection, 1.e-9) {
			t.Errorf("%d: intersection expected %g, got %g", i, test.intersection, intersection)
		}
		union := test.subject.Union(test.clipping)
		if !union.Similar(test.union, 1.e-9) {
			t.Errorf("%d: union expected %g, got %g", i, test.union, union)
		}
		difference := test.subject.Difference(test.clipping)
		if !difference.Similar(test.difference, 1.e-9) {
			t.Errorf("%d: difference expected %g, got %g", i, test.difference, difference)
		}
		xor := test.subject.XOr(test.clipping)
		if !xor.Similar(test.xor, 1.e-9) {
			t.Errorf("%d: xor expected %g, got %g", i, test.xor, xor)
		}
	}
}
