package geojson

import (
	"encoding/json"
	"geom"
	"reflect"
)

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type InvalidGeometryError struct{}

func (e InvalidGeometryError) Error() string {
	return "geojson: invalid geometry"
}

type UnsupportedGeometryError struct {
	Type string
}

func (e UnsupportedGeometryError) Error() string {
	return "geojson: unsupported geometry type " + e.Type
}

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

func unmarshalCoordinates(jsonCoordinates interface{}) []float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}
	coordinates := make([]float64, len(array))
	for i, element := range array {
		var ok bool
		if coordinates[i], ok = element.(float64); !ok {
			panic(&InvalidGeometryError{})
		}
	}
	return coordinates
}

func unmarshalCoordinates2(jsonCoordinates interface{}) [][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}
	coordinates := make([][]float64, len(array))
	for i, element := range array {
		coordinates[i] = unmarshalCoordinates(element)
	}
	return coordinates
}

func unmarshalCoordinates3(jsonCoordinates interface{}) [][][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}
	coordinates := make([][][]float64, len(array))
	for i, element := range array {
		coordinates[i] = unmarshalCoordinates2(element)
	}
	return coordinates
}

func makeLinearRing(coordinates [][]float64) geom.LinearRing {
	linearRing := make(geom.LinearRing, len(coordinates))
	for i, element := range coordinates {
		if len(element) == 2 {
			linearRing[i].X = element[0]
			linearRing[i].Y = element[1]
		} else {
			panic(&InvalidGeometryError{})
		}
	}
	return linearRing
}

func makeLinearRingZ(coordinates [][]float64) geom.LinearRingZ {
	linearRingZ := make(geom.LinearRingZ, len(coordinates))
	for i, element := range coordinates {
		if len(element) == 3 {
			linearRingZ[i].X = element[0]
			linearRingZ[i].Y = element[1]
			linearRingZ[i].Z = element[2]
		} else {
			panic(&InvalidGeometryError{})
		}
	}
	return linearRingZ
}

func makeLinearRings(coordinates [][][]float64) geom.LinearRings {
	linearRings := make(geom.LinearRings, len(coordinates))
	for i, element := range coordinates {
		linearRings[i] = makeLinearRing(element)
	}
	return linearRings
}

func makeLinearRingZs(coordinates [][][]float64) geom.LinearRingZs {
	linearRingZs := make(geom.LinearRingZs, len(coordinates))
	for i, element := range coordinates {
		linearRingZs[i] = makeLinearRingZ(element)
	}
	return linearRingZs
}

func doFromGeoJSON(g *Geometry) geom.T {
	switch g.Type {
	case "Point":
		coordinates := unmarshalCoordinates(g.Coordinates)
		switch len(coordinates) {
		case 2:
			return geom.Point{coordinates[0], coordinates[1]}
		case 3:
			return geom.PointZ{coordinates[0], coordinates[1], coordinates[2]}
		default:
			panic(&InvalidGeometryError{})
		}
	case "LineString":
		coordinates := unmarshalCoordinates2(g.Coordinates)
		if len(coordinates) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0]) {
		case 2:
			return geom.LineString{makeLinearRing(coordinates)}
		case 3:
			return geom.LineStringZ{makeLinearRingZ(coordinates)}
		default:
			panic(&InvalidGeometryError{})
		}
	case "Polygon":
		coordinates := unmarshalCoordinates3(g.Coordinates)
		if len(coordinates) == 0 || len(coordinates[0]) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0][0]) {
		case 2:
			return geom.Polygon{makeLinearRings(coordinates)}
		case 3:
			return geom.PolygonZ{makeLinearRingZs(coordinates)}
		default:
			panic(&InvalidGeometryError{})
		}
	default:
		panic(&UnsupportedGeometryError{g.Type})
	}
}

func FromGeoJSON(geom *Geometry) (g geom.T, err error) {
	defer func() {
		if e := recover(); e != nil {
			g = nil
			err = e.(error)
		}
	}()
	return doFromGeoJSON(geom), nil
}

func Unmarshal(data []byte) (geom.T, error) {
	var geom Geometry
	if err := json.Unmarshal(data, &geom); err == nil {
		return FromGeoJSON(&geom)
	} else {
		return nil, err
	}
}
