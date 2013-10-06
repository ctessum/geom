package geom

import (
	"encoding/json"
	"fmt"
)

type GeoJSONGeom interface {
	Geom
	GeoJSON() map[string]interface{}
}

func (point Point) geoJSONCoordinates() []float64 {
	return []float64{point.X, point.Y}
}

func (pointZ PointZ) geoJSONCoordinates() []float64 {
	return []float64{pointZ.X, pointZ.Y, pointZ.Z}
}

func (linearRing LinearRing) geoJSONCoordinates() [][]float64 {
	result := make([][]float64, len(linearRing))
	for i, point := range linearRing {
		result[i] = point.geoJSONCoordinates()
	}
	return result
}

func (linearRingZ LinearRingZ) geoJSONCoordinates() [][]float64 {
	result := make([][]float64, len(linearRingZ))
	for i, pointZ := range linearRingZ {
		result[i] = pointZ.geoJSONCoordinates()
	}
	return result
}

func (linearRings LinearRings) geoJSONCoordinates() [][][]float64 {
	result := make([][][]float64, len(linearRings))
	for i, linearRing := range linearRings {
		result[i] = linearRing.geoJSONCoordinates()
	}
	return result
}

func (linearRingZs LinearRingZs) geoJSONCoordinates() [][][]float64 {
	result := make([][][]float64, len(linearRingZs))
	for i, linearRingZ := range linearRingZs {
		result[i] = linearRingZ.geoJSONCoordinates()
	}
	return result
}

func (point Point) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Point",
		"coordinates": point.geoJSONCoordinates(),
	}
}

func (pointZ PointZ) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Point",
		"coordinates": pointZ.geoJSONCoordinates(),
	}
}

func (lineString LineString) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "LineString",
		"coordinates": lineString.Points.geoJSONCoordinates(),
	}
}

func (lineStringZ LineStringZ) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "LineString",
		"coordinates": lineStringZ.Points.geoJSONCoordinates(),
	}
}

func (polygon Polygon) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Polygon",
		"coordinates": polygon.Rings.geoJSONCoordinates(),
	}
}

func (polygonZ PolygonZ) GeoJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":        "Polygon",
		"coordinates": polygonZ.Rings.geoJSONCoordinates(),
	}
}

type geoJSONObject struct {
	Type        string
	Coordinates interface{}
}

func unmarshalCoordinates(jsonCoordinates interface{}) ([]float64, error) {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected an array")
	}
	result := make([]float64, len(array))
	for i, element := range array {
		var ok bool
		if result[i], ok = element.(float64); !ok {
			return nil, fmt.Errorf("expected a float64 in position %d", i)
		}
	}
	return result, nil
}

func unmarshalCoordinates2(jsonCoordinates interface{}) ([][]float64, error) {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected an array")
	}
	result := make([][]float64, len(array))
	for i, element := range array {
		var err error
		if result[i], err = unmarshalCoordinates(element); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func unmarshallCoordinates3(jsonCoordinates interface{}) ([][][]float64, error) {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected an array")
	}
	result := make([][][]float64, len(array))
	for i, element := range array {
		var err error
		if result[i], err = unmarshalCoordinates2(element); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func makeLinearRing(coordinates [][]float64) (LinearRing, error) {
	linearRing := make(LinearRing, len(coordinates))
	for i, element := range coordinates {
		if len(element) == 2 {
			linearRing[i].X = element[0]
			linearRing[i].Y = element[1]
		} else {
			return nil, fmt.Errorf("invalid length")
		}
	}
	return linearRing, nil
}

func makeLinearRingZ(coordinates [][]float64) (LinearRingZ, error) {
	linearRingZ := make(LinearRingZ, len(coordinates))
	for i, element := range coordinates {
		if len(element) == 3 {
			linearRingZ[i].X = element[0]
			linearRingZ[i].Y = element[1]
			linearRingZ[i].Z = element[2]
		} else {
			return nil, fmt.Errorf("invalid length")
		}
	}
	return linearRingZ, nil
}

func unmarshalPoint(jsonCoordinates interface{}) (Geom, error) {
	coordinates, err := unmarshalCoordinates(jsonCoordinates)
	if err != nil {
		return nil, err
	}
	switch len(coordinates) {
	case 2:
		return Point{coordinates[0], coordinates[1]}, nil
	case 3:
		return PointZ{coordinates[0], coordinates[1], coordinates[2]}, nil
	default:
		return nil, fmt.Errorf("invalid coordinates %q", coordinates)
	}
}

func unmarshalLineString(jsonCoordinates interface{}) (Geom, error) {
	coordinates, err := unmarshalCoordinates2(jsonCoordinates)
	if err != nil {
		return nil, err
	}
	if len(coordinates) == 0 {
		return nil, fmt.Errorf("empty coordinates array")
	}
	switch dim := len(coordinates[0]); dim {
	case 2:
		if linearRing, err := makeLinearRing(coordinates); err != nil {
			return nil, err
		} else {
			return LineString{linearRing}, nil
		}
	case 3:
		if linearRingZ, err := makeLinearRingZ(coordinates); err != nil {
			return nil, err
		} else {
			return LineStringZ{linearRingZ}, nil
		}
	default:
		return nil, fmt.Errorf("invalid dimension %d", dim)
	}
}

func makePolygon(coordinates [][][]float64) (Polygon, error) {
	rings := make([]LinearRing, len(coordinates))
	for i, element := range coordinates {
		var err error
		if rings[i], err = makeLinearRing(element); err != nil {
			return Polygon{}, err
		}
	}
	return Polygon{rings}, nil
}

func makePolygonZ(coordinates [][][]float64) (PolygonZ, error) {
	ringZs := make([]LinearRingZ, len(coordinates))
	for i, element := range coordinates {
		var err error
		if ringZs[i], err = makeLinearRingZ(element); err != nil {
			return PolygonZ{}, err
		}
	}
	return PolygonZ{ringZs}, nil
}

func unmarshalPolygon(jsonCoordinates interface{}) (Geom, error) {
	coordinates, err := unmarshallCoordinates3(jsonCoordinates)
	if err != nil {
		return nil, err
	}
	if len(coordinates) == 0 || len(coordinates[0]) == 0 {
		return nil, fmt.Errorf("empty coordinates array")
	}
	switch dim := len(coordinates[0][0]); dim {
	case 2:
		return makePolygon(coordinates)
	case 3:
		return makePolygonZ(coordinates)
	default:
		return nil, fmt.Errorf("invalid dimension %d", dim)
	}
}

func GeoJSONUnmarshal(b []byte) (Geom, error) {
	var o geoJSONObject
	if err := json.Unmarshal(b, &o); err != nil {
		return nil, err
	}
	if unmarshaler, ok := geoJSONUnmarshalers[o.Type]; !ok {
		return nil, fmt.Errorf("unsupported type %s", o.Type)
	} else {
		return unmarshaler(o.Coordinates)
	}
}

type geoJSONUnmarshaler func(interface{}) (Geom, error)

var geoJSONUnmarshalers = map[string]geoJSONUnmarshaler{
	"Point":      unmarshalPoint,
	"LineString": unmarshalLineString,
	"Polygon":    unmarshalPolygon,
}
