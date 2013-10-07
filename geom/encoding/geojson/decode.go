package geojson

import (
	"encoding/json"
	"github.com/twpayne/gogeom/geom"
)

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

func makeLinearRing(coordinates [][]float64) []geom.Point {
	points := make([]geom.Point, len(coordinates))
	for i, element := range coordinates {
		if len(element) == 2 {
			points[i].X = element[0]
			points[i].Y = element[1]
		} else {
			panic(&InvalidGeometryError{})
		}
	}
	return points
}

func makeLinearRingZ(coordinates [][]float64) []geom.PointZ {
	pointZs := make([]geom.PointZ, len(coordinates))
	for i, element := range coordinates {
		if len(element) == 3 {
			pointZs[i].X = element[0]
			pointZs[i].Y = element[1]
			pointZs[i].Z = element[2]
		} else {
			panic(&InvalidGeometryError{})
		}
	}
	return pointZs
}

func makeLinearRings(coordinates [][][]float64) [][]geom.Point {
	pointss := make([][]geom.Point, len(coordinates))
	for i, element := range coordinates {
		pointss[i] = makeLinearRing(element)
	}
	return pointss
}

func makeLinearRingZs(coordinates [][][]float64) [][]geom.PointZ {
	pointZss := make([][]geom.PointZ, len(coordinates))
	for i, element := range coordinates {
		pointZss[i] = makeLinearRingZ(element)
	}
	return pointZss
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
