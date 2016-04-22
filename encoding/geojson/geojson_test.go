package geojson

import (
	"github.com/ctessum/geom"
	"reflect"
	"testing"
)

func TestGeoJSON(t *testing.T) {
	testCases := []struct {
		g       geom.Geom
		geoJSON []byte
	}{
		{
			geom.Point{1, 2},
			[]byte(`{"type":"Point","coordinates":[1,2]}`),
		},
		{
			geom.LineString([]geom.Point{{1, 2}, {3, 4}}),
			[]byte(`{"type":"LineString","coordinates":[[1,2],[3,4]]}`),
		},
		{
			geom.Polygon([][]geom.Point{{{1, 2}, {3, 4}, {5, 6}}}),
			[]byte(`{"type":"Polygon","coordinates":[[[1,2],[3,4],[5,6]]]}`),
		},
	}
	for _, tc := range testCases {
		if got, err := Encode(tc.g); err != nil || !reflect.DeepEqual(got, tc.geoJSON) {
			t.Errorf("Encode(%#v) == %#v, %#v, want %#v, nil", tc.g, got, err, tc.geoJSON)
		}
		if got, err := Decode(tc.geoJSON); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("Decode(%#v) == %#v, %#v, want %#v, nil", tc.geoJSON, got, err, tc.g)
		}
	}
}

func TestGeoJSONDecode(t *testing.T) {
	testCases := [][]byte{
		[]byte(`{}`),
		[]byte(`{"type":""}`),
		[]byte(`{"type":"Point"}`),
		[]byte(`{"coordinates":[],"type":"Point"}`),
		[]byte(`{"coordinates":[1],"type":"Point"}`),
		[]byte(`{"coordinates":[1,2,3,4],"type":"Point"}`),
		[]byte(`{"coordinates":[""],"type":"Point"}`),
		[]byte(`{"type":"LineString"}`),
		[]byte(`{"coordinates":[],"type":"LineString"}`),
		[]byte(`{"coordinates":[[]],"type":"LineString"}`),
		[]byte(`{"coordinates":[1],"type":"LineString"}`),
		[]byte(`{"coordinates":[[1,2],[3,4,5]],"type":"LineString"}`),
		[]byte(`{"coordinates":[""],"type":"LineString"}`),
		[]byte(`{"coordinates":[[1,2,3,4],[5,6,7,8]],"type":"LineString"}`),
		[]byte(`{"type":"Polygon"}`),
		[]byte(`{"coordinates":[],"type":"Polygon"}`),
		[]byte(`{"coordinates":[[]],"type":"Polygon"}`),
		[]byte(`{"coordinates":[[[]]],"type":"Polygon"}`),
		[]byte(`{"coordinates":[[[1,2],[3,4,5]]],"type":"Polygon"}`),
	}
	for _, tc := range testCases {
		if got, err := Decode(tc); err == nil {
			t.Errorf("Decode(%#v) == %#v, nil, want err != nil", tc, got)
		}
	}
}
