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
			LineString([]Point{{1, 2}, {3, 4}}),
			&Bounds{Point{1, 2}, Point{3, 4}},
		},
		{
			Polygon([][]Point{{{1, 2}, {3, 4}, {5, 6}}}),
			&Bounds{Point{1, 2}, Point{5, 6}},
		},
		{
			MultiPoint([]Point{{1, 2}, {3, 4}}),
			&Bounds{Point{1, 2}, Point{3, 4}},
		},
		{
			MultiLineString([]LineString{[]Point{{1, 2}, {3, 4}}, []Point{{5, 6}, {7, 8}}}),
			&Bounds{Point{1, 2}, Point{7, 8}},
		},
	}

	for _, tc := range testCases {
		if got := tc.g.Bounds(); !reflect.DeepEqual(got, tc.bounds) {
			t.Errorf("%#v.Bounds() == %#v, want %#v", tc.g, got, tc.bounds)
		}
	}

}

func TestBoundsEmpty(t *testing.T) {
	if got := NewBounds().Empty(); got != true {
		t.Errorf("NewBounds.Empty() == %#v, want true", got)
	}
}
