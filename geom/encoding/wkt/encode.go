package wkt

import (
	"github.com/twpayne/gogeom/geom"
	"reflect"
)

func WKT(g geom.T) (string, error) {
	switch g.(type) {
	case geom.Point:
		return pointWKT(g.(geom.Point)), nil
	case geom.PointZ:
		return pointZWKT(g.(geom.PointZ)), nil
	case geom.PointM:
		return pointMWKT(g.(geom.PointM)), nil
	case geom.PointZM:
		return pointZMWKT(g.(geom.PointZM)), nil
	case geom.LineString:
		return lineStringWKT(g.(geom.LineString)), nil
	case geom.LineStringZ:
		return lineStringZWKT(g.(geom.LineStringZ)), nil
	case geom.LineStringM:
		return lineStringMWKT(g.(geom.LineStringM)), nil
	case geom.LineStringZM:
		return lineStringZMWKT(g.(geom.LineStringZM)), nil
	case geom.Polygon:
		return polygonWKT(g.(geom.Polygon)), nil
	case geom.PolygonZ:
		return polygonZWKT(g.(geom.PolygonZ)), nil
	case geom.PolygonM:
		return polygonMWKT(g.(geom.PolygonM)), nil
	case geom.PolygonZM:
		return polygonZMWKT(g.(geom.PolygonZM)), nil
	default:
		return "", &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}
