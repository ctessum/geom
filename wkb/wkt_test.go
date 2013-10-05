package geom

import (
	"testing"
)

func TestWKT(t *testing.T) {
	var testCases = []struct {
		g   WKTGeom
		wkt string
	}{
		{
			Point{1, 2},
			"POINT(1 2)",
		},
		{
			PointZ{1, 2, 3},
			"POINTZ(1 2 3)",
		},
		{
			PointM{1, 2, 3},
			"POINTM(1 2 3)",
		},
		{
			PointZM{1, 2, 3, 4},
			"POINTZM(1 2 3 4)",
		},
		{
			LineString{LinearRing{{1, 2}, {3, 4}}},
			"LINESTRING(1 2,3 4)",
		},
		{
			LineStringZ{LinearRingZ{{1, 2, 3}, {4, 5, 6}}},
			"LINESTRINGZ(1 2 3,4 5 6)",
		},
		{
			LineStringM{LinearRingM{{1, 2, 3}, {4, 5, 6}}},
			"LINESTRINGM(1 2 3,4 5 6)",
		},
		{
			LineStringZM{LinearRingZM{{1, 2, 3, 4}, {5, 6, 7, 8}}},
			"LINESTRINGZM(1 2 3 4,5 6 7 8)",
		},
		{
			Polygon{LinearRings{LinearRing{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}},
			"POLYGON((1 2,3 4,5 6,1 2))",
		},
		{
			PolygonZ{LinearRingZs{LinearRingZ{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}},
			"POLYGONZ((1 2 3,4 5 6,7 8 9,1 2 3))",
		},
		{
			PolygonM{LinearRingMs{LinearRingM{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {1, 2, 3}}}},
			"POLYGONM((1 2 3,4 5 6,7 8 9,1 2 3))",
		},
		{
			PolygonZM{LinearRingZMs{LinearRingZM{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {1, 2, 3, 4}}}},
			"POLYGONZM((1 2 3 4,5 6 7 8,9 10 11 12,1 2 3 4))",
		},
	}
	for _, tc := range testCases {
		if got := tc.g.WKT(); got != tc.wkt {
			t.Errorf("%q.WKT() == %q, want %q", tc.g, got, tc.wkt)
		}
	}
}
