package geojson

import (
	"reflect"
	"testing"

	"github.com/ctessum/geom"
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
			geom.MultiPoint{
				geom.Point{1, 2},
				geom.Point{3, 4},
			},
			[]byte(`{"type":"MultiPoint","coordinates":[[1,2],[3,4]]}`),
		},
		{
			geom.LineString(geom.Path{{1, 2}, {3, 4}}),
			[]byte(`{"type":"LineString","coordinates":[[1,2],[3,4]]}`),
		},
		{
			geom.MultiLineString{
				geom.LineString(geom.Path{{1, 2}, {3, 4}}),
				geom.LineString(geom.Path{{5, 6}, {7, 8}}),
			},
			[]byte(`{"type":"MultiLineString","coordinates":[[[1,2],[3,4]],[[5,6],[7,8]]]}`),
		},
		{
			geom.Polygon([]geom.Path{{{1, 2}, {3, 4}, {5, 6}}}),
			[]byte(`{"type":"Polygon","coordinates":[[[1,2],[3,4],[5,6]]]}`),
		},
		{
			geom.MultiPolygon{
				geom.Polygon([]geom.Path{{{1, 2}, {3, 4}, {5, 6}}}),
				geom.Polygon([]geom.Path{{{7, 8}, {9, 10}, {11, 12}}}),
			},
			[]byte(`{"type":"MultiPolygon","coordinates":[[[[1,2],[3,4],[5,6]]],[[[7,8],[9,10],[11,12]]]]}`),
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
		[]byte(`{"type":"MultiPoint"}`),
		[]byte(`{"coordinates":[1,2],type":"MultiPoint"}`),
		[]byte(`{"type":"LineString"}`),
		[]byte(`{"coordinates":[],"type":"LineString"}`),
		[]byte(`{"coordinates":[[]],"type":"LineString"}`),
		[]byte(`{"coordinates":[1],"type":"LineString"}`),
		[]byte(`{"coordinates":[[1,2],[3,4,5]],"type":"LineString"}`),
		[]byte(`{"coordinates":[""],"type":"LineString"}`),
		[]byte(`{"coordinates":[[1,2,3,4],[5,6,7,8]],"type":"LineString"}`),
		[]byte(`{"type":"MultiLineString"}`),
		[]byte(`{"coordinates":[[1,2,3,4],[5,6,7,8]],"type":"MultiLineString"}`),
		[]byte(`{"type":"Polygon"}`),
		[]byte(`{"coordinates":[],"type":"Polygon"}`),
		[]byte(`{"coordinates":[[]],"type":"Polygon"}`),
		[]byte(`{"coordinates":[[[]]],"type":"Polygon"}`),
		[]byte(`{"coordinates":[[[1,2],[3,4,5]]],"type":"Polygon"}`),
		[]byte(`{"type":"MultiPolygon"}`),
		[]byte(`{"coordinates":[[[1,2],[3,4,5]]],"type":"MultiPolygon"}`),
	}
	for _, tc := range testCases {
		if got, err := Decode(tc); err == nil {
			t.Errorf("Decode(%#v) == %#v, nil, want err != nil", tc, got)
		}
	}
}
