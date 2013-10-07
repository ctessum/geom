package geojson

import (
	"encoding/json"
	"github.com/twpayne/gogeom/geom"
	"reflect"
)

func pointCoordinates(point geom.Point) []float64 {
	return []float64{point.X, point.Y}
}

func pointZCoordinates(pointZ geom.PointZ) []float64 {
	return []float64{pointZ.X, pointZ.Y, pointZ.Z}
}

func pointsCoordinates(points []geom.Point) [][]float64 {
	coordinates := make([][]float64, len(points))
	for i, point := range points {
		coordinates[i] = pointCoordinates(point)
	}
	return coordinates
}

func pointZsCoordinates(pointZs []geom.PointZ) [][]float64 {
	coordinates := make([][]float64, len(pointZs))
	for i, pointZ := range pointZs {
		coordinates[i] = pointZCoordinates(pointZ)
	}
	return coordinates
}

func pointssCoordinates(pointss [][]geom.Point) [][][]float64 {
	coordinates := make([][][]float64, len(pointss))
	for i, points := range pointss {
		coordinates[i] = pointsCoordinates(points)
	}
	return coordinates
}

func pointZssCoordinates(pointZss [][]geom.PointZ) [][][]float64 {
	coordinates := make([][][]float64, len(pointZss))
	for i, pointZs := range pointZss {
		coordinates[i] = pointZsCoordinates(pointZs)
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
			Coordinates: pointsCoordinates(g.(geom.LineString).Points),
		}, nil
	case geom.LineStringZ:
		return &Geometry{
			Type:        "LineString",
			Coordinates: pointZsCoordinates(g.(geom.LineStringZ).Points),
		}, nil
	case geom.Polygon:
		return &Geometry{
			Type:        "Polygon",
			Coordinates: pointssCoordinates(g.(geom.Polygon).Rings),
		}, nil
	case geom.PolygonZ:
		return &Geometry{
			Type:        "Polygon",
			Coordinates: pointZssCoordinates(g.(geom.PolygonZ).Rings),
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
