package wkt

import (
	"github.com/ctessum/geom"
	"reflect"
	"testing"
)

func TestWKT(t *testing.T) {
	var testCases = []struct {
		g   geom.T
		wkt []byte
	}{
		{
			geom.Point{1, 2},
			[]byte(`POINT(1 2)`),
		},
		{
			geom.PointZ{1, 2, 3},
			[]byte(`POINTZ(1 2 3)`),
		},
		{
			geom.PointM{1, 2, 3},
			[]byte(`POINTM(1 2 3)`),
		},
		{
			geom.PointZM{1, 2, 3, 4},
			[]byte(`POINTZM(1 2 3 4)`),
		},
		{
			geom.LineString([]geom.Point{{1, 2}, {3, 4}}),
			[]byte(`LINESTRING(1 2,3 4)`),
		},
		{
			geom.LineStringZ([]geom.PointZ{{1, 2, 3}, {4, 5, 6}}),
			[]byte(`LINESTRINGZ(1 2 3,4 5 6)`),
		},
		{
			geom.LineStringM([]geom.PointM{{1, 2, 3}, {4, 5, 6}}),
			[]byte(`LINESTRINGM(1 2 3,4 5 6)`),
		},
		{
			geom.LineStringZM([]geom.PointZM{{1, 2, 3, 4}, {5, 6, 7, 8}}),
			[]byte(`LINESTRINGZM(1 2 3 4,5 6 7 8)`),
		},
		{
			geom.Polygon([][]geom.Point{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}),
			[]byte(`POLYGON((1 2,3 4,5 6,1 2))`),
		},
		{
			geom.PolygonZ([][]geom.PointZ{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}),
			[]byte(`POLYGONZ((1 2 3,4 5 6,7 8 9,1 2 3))`),
		},
		{
			geom.PolygonM([][]geom.PointM{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}),
			[]byte(`POLYGONM((1 2 3,4 5 6,7 8 9,1 2 3))`),
		},
		{
			geom.PolygonZM([][]geom.PointZM{{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {1, 2, 3, 4}}}),
			[]byte(`POLYGONZM((1 2 3 4,5 6 7 8,9 10 11 12,1 2 3 4))`),
		},
	}
	for _, tc := range testCases {
		if got, err := Encode(tc.g); err != nil || !reflect.DeepEqual(got, tc.wkt) {
			t.Errorf("Encode(%#v) == %#v, %#v, want %#v, nil", tc.g, got, err, tc.wkt)
		}
	}
}
