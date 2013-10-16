package wkt

import (
	"github.com/twpayne/gogeom/geom"
	"reflect"
)

func Encode(g geom.T) ([]byte, error) {
	switch g.(type) {
	case geom.Point:
		point := g.(geom.Point)
		return appendPointWKT(nil, &point), nil
	case geom.PointZ:
		pointZ := g.(geom.PointZ)
		return appendPointZWKT(nil, &pointZ), nil
	case geom.PointM:
		pointM := g.(geom.PointM)
		return appendPointMWKT(nil, &pointM), nil
	case geom.PointZM:
		pointZM := g.(geom.PointZM)
		return appendPointZMWKT(nil, &pointZM), nil
	case geom.LineString:
		lineString := g.(geom.LineString)
		return appendLineStringWKT(nil, &lineString), nil
	case geom.LineStringZ:
		lineStringZ := g.(geom.LineStringZ)
		return appendLineStringZWKT(nil, &lineStringZ), nil
	case geom.LineStringM:
		lineStringM := g.(geom.LineStringM)
		return appendLineStringMWKT(nil, &lineStringM), nil
	case geom.LineStringZM:
		lineStringZM := g.(geom.LineStringZM)
		return appendLineStringZMWKT(nil, &lineStringZM), nil
	case geom.Polygon:
		polygon := g.(geom.Polygon)
		return appendPolygonWKT(nil, &polygon), nil
	case geom.PolygonZ:
		polygonZ := g.(geom.PolygonZ)
		return appendPolygonZWKT(nil, &polygonZ), nil
	case geom.PolygonM:
		polygonM := g.(geom.PolygonM)
		return appendPolygonMWKT(nil, &polygonM), nil
	case geom.PolygonZM:
		polygonZM := g.(geom.PolygonZM)
		return appendPolygonZMWKT(nil, &polygonZM), nil
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}
