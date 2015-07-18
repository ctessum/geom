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
			geom.LineString([]geom.Point{{1, 2}, {3, 4}}),
			[]byte(`LINESTRING(1 2,3 4)`),
		},
		{
			geom.Polygon([][]geom.Point{{{1, 2}, {3, 4}, {5, 6}, {1, 2}}}),
			[]byte(`POLYGON((1 2,3 4,5 6,1 2))`),
		},
	}
	for _, tc := range testCases {
		if got, err := Encode(tc.g); err != nil || !reflect.DeepEqual(got, tc.wkt) {
			t.Errorf("Encode(%#v) == %#v, %#v, want %#v, nil", tc.g, got, err, tc.wkt)
		}
	}
}
