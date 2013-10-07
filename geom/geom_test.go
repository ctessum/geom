package geom

import (
	"reflect"
	"testing"
)

func TestBounds(t *testing.T) {

	var testCases = []struct {
		g      T
		bounds *Bounds
	}{
		{
			Point{1, 2},
			&Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			PointZ{1, 2, 3},
			&Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			PointM{1, 2, 3},
			&Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			PointZM{1, 2, 3, 4},
			&Bounds{Point{1, 2}, Point{1, 2}},
		},
		{
			LineString{[]Point{{1, 2}, {3, 4}}},
			&Bounds{Point{1, 2}, Point{3, 4}},
		},
		{
			LineStringZ{[]PointZ{{1, 2, 3}, {4, 5, 6}}},
			&Bounds{Point{1, 2}, Point{4, 5}},
		},
		{
			LineStringM{[]PointM{{1, 2, 3}, {4, 5, 6}}},
			&Bounds{Point{1, 2}, Point{4, 5}},
		},
		{
			LineStringZM{[]PointZM{{1, 2, 3, 4}, {5, 6, 7, 8}}},
			&Bounds{Point{1, 2}, Point{5, 6}},
		},
		{
			Polygon{[][]Point{{{1, 2}, {3, 4}, {5, 6}}}},
			&Bounds{Point{1, 2}, Point{5, 6}},
		},
	}

	for _, tc := range testCases {
		if got := tc.g.Bounds(); !reflect.DeepEqual(got, tc.bounds) {
			t.Errorf("%q.Bounds() == %q, want %q", tc.g, got, tc.bounds)
		}
	}

}

func TestBoundsEmpty(t *testing.T) {
	if got := NewBounds().Empty(); got != true {
		t.Errorf("NewBounds.Empty() == %q, want true", got)
	}
}
