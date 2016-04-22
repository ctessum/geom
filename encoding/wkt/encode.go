package wkt

import (
	"github.com/ctessum/geom"
	"reflect"
)

func Encode(g geom.Geom) ([]byte, error) {
	switch g.(type) {
	case geom.Point:
		point := g.(geom.Point)
		return appendPointWKT(nil, &point), nil
	case geom.LineString:
		lineString := g.(geom.LineString)
		return appendLineStringWKT(nil, lineString), nil
	case geom.MultiLineString:
		multiLineString := g.(geom.MultiLineString)
		return appendMultiLineStringWKT(nil, multiLineString), nil
	case geom.Polygon:
		polygon := g.(geom.Polygon)
		return appendPolygonWKT(nil, polygon), nil
	case geom.MultiPolygon:
		multiPolygon := g.(geom.MultiPolygon)
		return appendMultiPolygonWKT(nil, multiPolygon), nil
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}
