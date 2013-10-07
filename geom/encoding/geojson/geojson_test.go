package geojson

import (
	"github.com/twpayne/gogeom/geom"
	"reflect"
	"testing"
)

func TestGeoJSON(t *testing.T) {
	testCases := []struct {
		g       geom.T
		geoJSON []byte
	}{
		{
			geom.Point{1, 2},
			[]byte(`{"type":"Point","coordinates":[1,2]}`),
		},
		{
			geom.PointZ{1, 2, 3},
			[]byte(`{"type":"Point","coordinates":[1,2,3]}`),
		},
		{
			geom.LineString{[]geom.Point{{1, 2}, {3, 4}}},
			[]byte(`{"type":"LineString","coordinates":[[1,2],[3,4]]}`),
		},
		{
			geom.LineStringZ{[]geom.PointZ{{1, 2, 3}, {3, 4, 5}}},
			[]byte(`{"type":"LineString","coordinates":[[1,2,3],[3,4,5]]}`),
		},
		{
			geom.Polygon{[][]geom.Point{{{1, 2}, {3, 4}, {5, 6}}}},
			[]byte(`{"type":"Polygon","coordinates":[[[1,2],[3,4],[5,6]]]}`),
		},
		{
			geom.PolygonZ{[][]geom.PointZ{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}}},
			[]byte(`{"type":"Polygon","coordinates":[[[1,2,3],[4,5,6],[7,8,9]]]}`),
		},
	}
	for _, tc := range testCases {
		if got, err := Marshal(tc.g); err != nil || !reflect.DeepEqual(got, tc.geoJSON) {
			t.Errorf("Marshal(%q) == %q, %q, want %q, nil", tc.g, got, err, tc.geoJSON)
		}
		if got, err := Unmarshal(tc.geoJSON); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("Unmarshal(%q) == %q, %q, want %q, nil", tc.geoJSON, got, err, tc.g)
		}
	}
}

func TestGeoJSONUnmarshallErrors(t *testing.T) {
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
		if got, err := Unmarshal(tc); err == nil {
			t.Errorf("Unmarshal(%q) == %q, nil, want err != nil", tc, got)
		}
	}
}
