package wkt

import (
	"geom"
	"testing"
)

func TestWKT(t *testing.T) {
	var testCases = []struct {
		g   geom.T
		wkt string
	}{
		{
			geom.Point{1, 2},
			"POINT(1 2)",
		},
		{
			geom.PointZ{1, 2, 3},
			"POINTZ(1 2 3)",
		},
		{
			geom.PointM{1, 2, 3},
			"POINTM(1 2 3)",
		},
		{
			geom.PointZM{1, 2, 3, 4},
			"POINTZM(1 2 3 4)",
		},
		{
			geom.LineString{geom.LinearRing{{1, 2}, {3, 4}}},
			"LINESTRING(1 2,3 4)",
		},
		{
			geom.LineStringZ{geom.LinearRingZ{{1, 2, 3}, {4, 5, 6}}},
			"LINESTRINGZ(1 2 3,4 5 6)",
		},
		{
			geom.LineStringM{geom.LinearRingM{{1, 2, 3}, {4, 5, 6}}},
			"LINESTRINGM(1 2 3,4 5 6)",
		},
		{
			geom.LineStringZM{geom.LinearRingZM{{1, 2, 3, 4}, {5, 6, 7, 8}}},
			"LINESTRINGZM(1 2 3 4,5 6 7 8)",
		},
		{
			geom.Polygon{geom.LinearRings{geom.LinearRing{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}},
			"POLYGON((1 2,3 4,5 6,1 2))",
		},
		{
			geom.PolygonZ{geom.LinearRingZs{geom.LinearRingZ{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}},
			"POLYGONZ((1 2 3,4 5 6,7 8 9,1 2 3))",
		},
		{
			geom.PolygonM{geom.LinearRingMs{geom.LinearRingM{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}},
			"POLYGONM((1 2 3,4 5 6,7 8 9,1 2 3))",
		},
		{
			geom.PolygonZM{geom.LinearRingZMs{geom.LinearRingZM{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {1, 2, 3, 4}}}},
			"POLYGONZM((1 2 3 4,5 6 7 8,9 10 11 12,1 2 3 4))",
		},
	}
	for _, tc := range testCases {
		if got, err := WKT(tc.g); err != nil || got != tc.wkt {
			t.Errorf("WKT(%q) == %q, %q, want %q, nil", tc.g, got, err, tc.wkt)
		}
	}
}
