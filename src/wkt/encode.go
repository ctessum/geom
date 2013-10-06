package wkt

import (
	"fmt"
	"geom"
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

func linearRingWKCoordinates(linearRing geom.LinearRing) string {
	wktCoordinates := make([]string, len(linearRing))
	for i, point := range linearRing {
		wktCoordinates[i] = pointWKTCoordinates(point)
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingZWKCoordinates(linearRingZ geom.LinearRingZ) string {
	wktCoordinates := make([]string, len(linearRingZ))
	for i, pointZ := range linearRingZ {
		wktCoordinates[i] = pointZWKTCoordinates(pointZ)
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingMWKCoordinates(linearRingM geom.LinearRingM) string {
	wktCoordinates := make([]string, len(linearRingM))
	for i, pointM := range linearRingM {
		wktCoordinates[i] = pointMWKTCoordinates(pointM)
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingZMWKCoordinates(linearRingZM geom.LinearRingZM) string {
	wktCoordinates := make([]string, len(linearRingZM))
	for i, pointZM := range linearRingZM {
		wktCoordinates[i] = pointZMWKTCoordinates(pointZM)
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingsWKTCoordinates(linearRings geom.LinearRings) string {
	wktCoordinates := make([]string, len(linearRings))
	for i, linearRing := range linearRings {
		wktCoordinates[i] = "(" + linearRingWKCoordinates(linearRing) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingZsWKTCoordinates(linearRingZs geom.LinearRingZs) string {
	wktCoordinates := make([]string, len(linearRingZs))
	for i, linearRingZ := range linearRingZs {
		wktCoordinates[i] = "(" + linearRingZWKCoordinates(linearRingZ) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingMsWKTCoordinates(linearRingMs geom.LinearRingMs) string {
	wktCoordinates := make([]string, len(linearRingMs))
	for i, linearRingM := range linearRingMs {
		wktCoordinates[i] = "(" + linearRingMWKCoordinates(linearRingM) + ")"
	}
	return strings.Join(wktCoordinates, ",")
}

func linearRingZMsWKTCoordinates(linearRingZMs geom.LinearRingZMs) string {
	wktCoordinates := make([]string, len(linearRingZMs))
	for i, linearRingZM := range linearRingZMs {
		wktCoordinates[i] = "(" + linearRingZMWKCoordinates(linearRingZM) + ")"
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
	return "LINESTRING(" + linearRingWKCoordinates(lineString.Points) + ")"
}

func lineStringZWKT(lineStringZ geom.LineStringZ) string {
	return "LINESTRINGZ(" + linearRingZWKCoordinates(lineStringZ.Points) + ")"
}

func lineStringMWKT(lineStringM geom.LineStringM) string {
	return "LINESTRINGM(" + linearRingMWKCoordinates(lineStringM.Points) + ")"
}

func lineStringZMWKT(lineStringZM geom.LineStringZM) string {
	return "LINESTRINGZM(" + linearRingZMWKCoordinates(lineStringZM.Points) + ")"
}

func polygonWKT(polygon geom.Polygon) string {
	return "POLYGON(" + linearRingsWKTCoordinates(polygon.Rings) + ")"
}

func polygonZWKT(polygonZ geom.PolygonZ) string {
	return "POLYGONZ(" + linearRingZsWKTCoordinates(polygonZ.Rings) + ")"
}

func polygonMWKT(polygonM geom.PolygonM) string {
	return "POLYGONM(" + linearRingMsWKTCoordinates(polygonM.Rings) + ")"
}

func polygonZMWKT(polygonZM geom.PolygonZM) string {
	return "POLYGONZM(" + linearRingZMsWKTCoordinates(polygonZM.Rings) + ")"
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
