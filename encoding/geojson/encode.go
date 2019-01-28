package geojson

import (
	"encoding/json"
	"reflect"

	"github.com/ctessum/geom"
)

func pointCoordinates(point geom.Point) []float64 {
	return []float64{point.X, point.Y}
}

func pointsCoordinates(points []geom.Point) [][]float64 {
	coordinates := make([][]float64, len(points))
	for i, point := range points {
		coordinates[i] = pointCoordinates(point)
	}
	return coordinates
}

func pointssCoordinates(pointss []geom.Path) [][][]float64 {
	coordinates := make([][][]float64, len(pointss))
	for i, points := range pointss {
		coordinates[i] = pointsCoordinates(points)
	}
	return coordinates
}

func pointsssCoordinates(pointsss [][]geom.Path) [][][][]float64 {
	coordinates := make([][][][]float64, len(pointsss))
	for i, points := range pointsss {
		coordinates[i] = pointssCoordinates(points)
	}
	return coordinates
}

func ToGeoJSON(g geom.Geom) (*Geometry, error) {
	switch g.(type) {
	case geom.Point:
		return &Geometry{
			Type:        "Point",
			Coordinates: pointCoordinates(g.(geom.Point)),
		}, nil
	case geom.MultiPoint:
		return &Geometry{
			Type:        "MultiPoint",
			Coordinates: pointsCoordinates(g.(geom.MultiPoint)),
		}, nil
	case geom.LineString:
		return &Geometry{
			Type:        "LineString",
			Coordinates: pointsCoordinates(g.(geom.LineString)),
		}, nil
	case geom.MultiLineString:
		lines := []geom.LineString(g.(geom.MultiLineString))
		paths := make([]geom.Path, len(lines))
		for i, line := range lines {
			paths[i] = geom.Path(line)
		}
		return &Geometry{
			Type:        "MultiLineString",
			Coordinates: pointssCoordinates(paths),
		}, nil
	case geom.Polygon:
		return &Geometry{
			Type:        "Polygon",
			Coordinates: pointssCoordinates(g.(geom.Polygon)),
		}, nil
	case geom.MultiPolygon:
		polys := []geom.Polygon(g.(geom.MultiPolygon))
		pathsList := make([][]geom.Path, len(polys))
		for i, poly := range polys {
			pathsList[i] = []geom.Path(poly)
		}
		return &Geometry{
			Type:        "MultiPolygon",
			Coordinates: pointsssCoordinates(pathsList),
		}, nil
	default:
		return nil, &UnsupportedGeometryError{reflect.TypeOf(g).String()}
	}
}

func Encode(g geom.Geom) ([]byte, error) {
	if object, err := ToGeoJSON(g); err == nil {
		return json.Marshal(object)
	} else {
		return nil, err
	}
}
