package geom

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestGeoJSON(t *testing.T) {
	testCases := []struct {
		g       GeoJSONGeom
		geoJSON []byte
	}{
		{
			Point{1, 2},
			[]byte(`{"coordinates":[1,2],"type":"Point"}`),
		},
		{
			PointZ{1, 2, 3},
			[]byte(`{"coordinates":[1,2,3],"type":"Point"}`),
		},
		{
			LineString{LinearRing{{1, 2}, {3, 4}}},
			[]byte(`{"coordinates":[[1,2],[3,4]],"type":"LineString"}`),
		},
		{
			LineStringZ{LinearRingZ{{1, 2, 3}, {3, 4, 5}}},
			[]byte(`{"coordinates":[[1,2,3],[3,4,5]],"type":"LineString"}`),
		},
		{
			Polygon{LinearRings{{{1, 2}, {3, 4}, {5, 6}}}},
			[]byte(`{"coordinates":[[[1,2],[3,4],[5,6]]],"type":"Polygon"}`),
		},
		{
			PolygonZ{LinearRingZs{{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}}},
			[]byte(`{"coordinates":[[[1,2,3],[4,5,6],[7,8,9]]],"type":"Polygon"}`),
		},
	}
	for _, tc := range testCases {
		if got, err := json.Marshal(tc.g.GeoJSON()); err != nil || !reflect.DeepEqual(got, tc.geoJSON) {
			t.Errorf("GeoJSON(%q) == %q, %q, want %q, nil", tc.g, got, err, tc.geoJSON)
		}
		if got, err := GeoJSONUnmarshal(tc.geoJSON); err != nil || !reflect.DeepEqual(got, tc.g) {
			t.Errorf("GeoJSONUnmarshal(%q) == %q, %q, want %q, nil", tc.geoJSON, got, err, tc.g)
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
		if got, err := GeoJSONUnmarshal(tc); err == nil {
			t.Errorf("GeoJSONUnmarshal(%q) == %q, nil, want err != nil", tc, got)
		}
	}
}
