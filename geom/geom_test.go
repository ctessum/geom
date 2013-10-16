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
		{
			MultiPoint{[]Point{{1, 2}, {3, 4}}},
			&Bounds{Point{1, 2}, Point{3, 4}},
		},
		{
			MultiPointZ{[]PointZ{{1, 2, 3}, {4, 5, 6}}},
			&Bounds{Point{1, 2}, Point{4, 5}},
		},
		{
			MultiPointM{[]PointM{{1, 2, 3}, {4, 5, 6}}},
			&Bounds{Point{1, 2}, Point{4, 5}},
		},
		{
			MultiPointZM{[]PointZM{{1, 2, 3, 4}, {5, 6, 7, 8}}},
			&Bounds{Point{1, 2}, Point{5, 6}},
		},
		{
			MultiLineString{[]LineString{{[]Point{{1, 2}, {3, 4}}}, {[]Point{{5, 6}, {7, 8}}}}},
			&Bounds{Point{1, 2}, Point{7, 8}},
		},
		{
			MultiLineStringZ{[]LineStringZ{{[]PointZ{{1, 2, 3}, {4, 5, 6}}}, {[]PointZ{{7, 8, 9}, {10, 11, 12}}}}},
			&Bounds{Point{1, 2}, Point{10, 11}},
		},
		{
			MultiLineStringM{[]LineStringM{{[]PointM{{1, 2, 3}, {4, 5, 6}}}, {[]PointM{{7, 8, 9}, {10, 11, 12}}}}},
			&Bounds{Point{1, 2}, Point{10, 11}},
		},
		{
			MultiLineStringZM{[]LineStringZM{{[]PointZM{{1, 2, 3, 4}, {5, 6, 7, 8}}}, {[]PointZM{{9, 10, 11, 12}, {13, 14, 15, 16}}}}},
			&Bounds{Point{1, 2}, Point{13, 14}},
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
