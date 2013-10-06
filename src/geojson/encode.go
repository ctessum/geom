package geojson

import (
	"encoding/json"
	"geom"
	"reflect"
)

func pointCoordinates(point geom.Point) []float64 {
	return []float64{point.X, point.Y}
}

func pointZCoordinates(pointZ geom.PointZ) []float64 {
	return []float64{pointZ.X, pointZ.Y, pointZ.Z}
}

func linearRingCoordinates(linearRing geom.LinearRing) [][]float64 {
	coordinates := make([][]float64, len(linearRing))
	for i, point := range linearRing {
		coordinates[i] = pointCoordinates(point)
	}
	return coordinates
}

func linearRingZCoordinates(linearRingZ geom.LinearRingZ) [][]float64 {
	coordinates := make([][]float64, len(linearRingZ))
	for i, pointZ := range linearRingZ {
		coordinates[i] = pointZCoordinates(pointZ)
	}
	return coordinates
}

func linearRingsCoordinates(linearRings geom.LinearRings) [][][]float64 {
	coordinates := make([][][]float64, len(linearRings))
	for i, linearRing := range linearRings {
		coordinates[i] = linearRingCoordinates(linearRing)
	}
	return coordinates
}

func linearRingZsCoordinates(linearRingZs geom.LinearRingZs) [][][]float64 {
	coordinates := make([][][]float64, len(linearRingZs))
	for i, linearRingZ := range linearRingZs {
		coordinates[i] = linearRingZCoordinates(linearRingZ)
	}
	return coordinates
}

func ToGeoJSON(g geom.T) (*Geometry, error) {
	switch g.(type) {
	case geom.Point:
		return &Geometry{
			Type:        "Point",
			Coordinates: pointCoordinates(g.(geom.Point)),
		}, nil
	case geom.PointZ:
		return &Geometry{
			Type:        "Point",
			Coordinates: pointZCoordinates(g.(geom.PointZ)),
		}, nil
	case geom.LineString:
		return &Geometry{
			Type:        "LineString",
			Coordinates: linearRingCoordinates(g.(geom.LineString).Points),
		}, nil
	case geom.LineStringZ:
		return &Geometry{
			Type:        "LineString",
			Coordinates: linearRingZCoordinates(g.(geom.LineStringZ).Points),
		}, nil
	case geom.Polygon:
		return &Geometry{
			Type:        "Polygon",
			Coordinates: linearRingsCoordinates(g.(geom.Polygon).Rings),
		}, nil
	case geom.PolygonZ:
		return &Geometry{
			Type:        "Polygon",
			Coordinates: linearRingZsCoordinates(g.(geom.PolygonZ).Rings),
		}, nil
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g).String()}
	}
}

func Marshal(g geom.T) ([]byte, error) {
	if object, err := ToGeoJSON(g); err == nil {
		return json.Marshal(object)
	} else {
		return nil, err
	}
}
