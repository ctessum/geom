package wkt

import (
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"reflect"
	"strings"
)

func pointWKTCoordinates(point geom.Point) string {
	return fmt.Sprintf("%g %g", point.X, point.Y)
}

func pointZWKTCoordinates(pointZ geom.PointZ) string {
	return fmt.Sprintf("%g %g %g", pointZ.X, pointZ.Y, pointZ.Z)
}

func pointMWKTCoordinates(pointM geom.PointM) string {
	return fmt.Sprintf("%g %g %g", pointM.X, pointM.Y, pointM.M)
}

func pointZMWKTCoordinates(pointZM geom.PointZM) string {
	return fmt.Sprintf("%g %g %g %g", pointZM.X, pointZM.Y, pointZM.Z, pointZM.M)
}

func pointsWKCoordinates(points []geom.Point) string {
	wktCoordinates := make([]string, len(points))
	for i, point := range points {
		wktCoordinates[i] = pointWKTCoordinates(point)
	}
	return strings.Join(wktCoordinates, ",")
}

func pointZsWKCoordinates(pointZs []geom.PointZ) string {
	wktCoordinates := make([]string, len(pointZs))
	for i, pointZ := range pointZs {
		wktCoordinates[i] = pointZWKTCoordinates(pointZ)
	}
	return strings.Join(wktCoordinates, ",")
}

func pointMsWKCoordinates(pointMs []geom.PointM) string {
	wktCoordinates := make([]string, len(pointMs))
	for i, pointM := range pointMs {
		wktCoordinates[i] = pointMWKTCoordinates(pointM)
	}
	return strings.Join(wktCoordinates, ",")
}

func pointZMsWKCoordinates(pointZMs []geom.PointZM) string {
	wktCoordinates := make([]string, len(pointZMs))
	for i, pointZM := range pointZMs {
		wktCoordinates[i] = pointZMWKTCoordinates(pointZM)
	}
	return strings.Join(wktCoordinates, ",")
}

func pointssWKTCoordinates(pointss [][]geom.Point) string {
	wktCoordinates := make([]string, len(pointss))
	for i, points := range pointss {
		wktCoordinates[i] = "(" + pointsWKCoordinates(points) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func pointZssWKTCoordinates(pointZss [][]geom.PointZ) string {
	wktCoordinates := make([]string, len(pointZss))
	for i, pointZs := range pointZss {
		wktCoordinates[i] = "(" + pointZsWKCoordinates(pointZs) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func pointMssWKTCoordinates(pointMss [][]geom.PointM) string {
	wktCoordinates := make([]string, len(pointMss))
	for i, pointMs := range pointMss {
		wktCoordinates[i] = "(" + pointMsWKCoordinates(pointMs) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func pointZMssWKTCoordinates(pointZMss [][]geom.PointZM) string {
	wktCoordinates := make([]string, len(pointZMss))
	for i, pointZMs := range pointZMss {
		wktCoordinates[i] = "(" + pointZMsWKCoordinates(pointZMs) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func pointWKT(point geom.Point) string {
	return "POINT(" + pointWKTCoordinates(point) + ")"
}

func pointZWKT(pointZ geom.PointZ) string {
	return "POINTZ(" + pointZWKTCoordinates(pointZ) + ")"
}

func pointMWKT(pointM geom.PointM) string {
	return "POINTM(" + pointMWKTCoordinates(pointM) + ")"
}

func pointZMWKT(pointZM geom.PointZM) string {
	return "POINTZM(" + pointZMWKTCoordinates(pointZM) + ")"
}

func lineStringWKT(lineString geom.LineString) string {
	return "LINESTRING(" + pointsWKCoordinates(lineString.Points) + ")"
}

func lineStringZWKT(lineStringZ geom.LineStringZ) string {
	return "LINESTRINGZ(" + pointZsWKCoordinates(lineStringZ.Points) + ")"
}

func lineStringMWKT(lineStringM geom.LineStringM) string {
	return "LINESTRINGM(" + pointMsWKCoordinates(lineStringM.Points) + ")"
}

func lineStringZMWKT(lineStringZM geom.LineStringZM) string {
	return "LINESTRINGZM(" + pointZMsWKCoordinates(lineStringZM.Points) + ")"
}

func polygonWKT(polygon geom.Polygon) string {
	return "POLYGON(" + pointssWKTCoordinates(polygon.Rings) + ")"
}

func polygonZWKT(polygonZ geom.PolygonZ) string {
	return "POLYGONZ(" + pointZssWKTCoordinates(polygonZ.Rings) + ")"
}

func polygonMWKT(polygonM geom.PolygonM) string {
	return "POLYGONM(" + pointMssWKTCoordinates(polygonM.Rings) + ")"
}

func polygonZMWKT(polygonZM geom.PolygonZM) string {
	return "POLYGONZM(" + pointZMssWKTCoordinates(polygonZM.Rings) + ")"
}

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
