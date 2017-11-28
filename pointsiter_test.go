package geom

import (
	"testing"
)

func TestGeom_points(t *testing.T) {
	tests := []Geom{
		Point{0, 1},
		MultiPoint{{0, 1}, {0, 2}},
		LineString{{0, 1}, {0, 2}},
		MultiLineString{{{0, 1}, {0, 2}}, {{0, 3}, {0, 4}}},
		Polygon{
			{{0, 1}, {0, 2}}, {{0, 3}, {0, 4}},
			{{0, 5}, {0, 6}}, {{0, 7}, {0, 8}},
		},
		MultiPolygon{
			{
				{{0, 1}, {0, 2}}, {{0, 3}, {0, 4}},
				{{0, 5}, {0, 6}}, {{0, 7}, {0, 8}},
			},
			{
				{{0, 9}, {0, 10}},
				{{0, 11}, {0, 12}},
			},
		},
		GeometryCollection{
			Point{0, 1},
			MultiPoint{{0, 2}, {0, 3}},
		},
	}

	for j, test := range tests {
		pf := test.Points()
		for i := 0; i < test.Len(); i++ {
			want := Point{0, float64(i + 1)}
			p := pf()
			if p != want {
				t.Errorf("test %d index %d: %+v != %+v", j, i, p, want)
			}
		}
	}
}
