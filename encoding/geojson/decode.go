package geojson

import (
	"encoding/json"

	"github.com/ctessum/geom"
)

func decodeCoordinates(jsonCoordinates interface{}) []float64 {
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

func decodeCoordinates2(jsonCoordinates interface{}) [][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}
	coordinates := make([][]float64, len(array))
	for i, element := range array {
		coordinates[i] = decodeCoordinates(element)
	}
	return coordinates
}

func decodeCoordinates3(jsonCoordinates interface{}) [][][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}
	coordinates := make([][][]float64, len(array))
	for i, element := range array {
		coordinates[i] = decodeCoordinates2(element)
	}
	return coordinates
}

func decodeCoordinates4(jsonCoordinates interface{}) [][][][]float64 {
	array, ok := jsonCoordinates.([]interface{})
	if !ok {
		panic(&InvalidGeometryError{})
	}
	coordinates := make([][][][]float64, len(array))
	for i, element := range array {
		coordinates[i] = decodeCoordinates3(element)
	}
	return coordinates
}

func makeLinearRing(coordinates [][]float64) geom.Path {
	points := make(geom.Path, len(coordinates))
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

func makeLinearRings(coordinates [][][]float64) []geom.Path {
	pointss := make([]geom.Path, len(coordinates))
	for i, element := range coordinates {
		pointss[i] = makeLinearRing(element)
	}
	return pointss
}

func doFromGeoJSON(g *Geometry) geom.Geom {
	switch g.Type {
	case "Point":
		coordinates := decodeCoordinates(g.Coordinates)
		switch len(coordinates) {
		case 2:
			return geom.Point{coordinates[0], coordinates[1]}
		default:
			panic(&InvalidGeometryError{})
		}
	case "MultiPoint":
		coordinates := decodeCoordinates2(g.Coordinates)
		if len(coordinates) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0]) {
		case 2:
			return geom.MultiPoint(makeLinearRing(coordinates))
		default:
			panic(&InvalidGeometryError{})
		}
	case "LineString":
		coordinates := decodeCoordinates2(g.Coordinates)
		if len(coordinates) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0]) {
		case 2:
			return geom.LineString(makeLinearRing(coordinates))
		default:
			panic(&InvalidGeometryError{})
		}
	case "MultiLineString":
		coordinates := decodeCoordinates3(g.Coordinates)
		if len(coordinates) == 0 || len(coordinates[0]) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0][0]) {
		case 2:
			multiLineString := make(geom.MultiLineString, len(coordinates))
			for i, coord := range coordinates {
				multiLineString[i] = geom.LineString(makeLinearRing(coord))
			}
			return multiLineString
		default:
			panic(&InvalidGeometryError{})
		}
	case "Polygon":
		coordinates := decodeCoordinates3(g.Coordinates)
		if len(coordinates) == 0 || len(coordinates[0]) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0][0]) {
		case 2:
			return geom.Polygon(makeLinearRings(coordinates))
		default:
			panic(&InvalidGeometryError{})
		}
	case "MultiPolygon":
		coordinates := decodeCoordinates4(g.Coordinates)
		if len(coordinates) == 0 || len(coordinates[0]) == 0 || len(coordinates[0][0]) == 0 {
			panic(&InvalidGeometryError{})
		}
		switch len(coordinates[0][0][0]) {
		case 2:
			multiPolygon := make(geom.MultiPolygon, len(coordinates))
			for i, coord := range coordinates {
				multiPolygon[i] = makeLinearRings(coord)
			}
			return multiPolygon
		default:
			panic(&InvalidGeometryError{})
		}
	default:
		panic(&UnsupportedGeometryError{g.Type})
	}
}

func FromGeoJSON(geom *Geometry) (g geom.Geom, err error) {
	defer func() {
		if e := recover(); e != nil {
			g = nil
			err = e.(error)
		}
	}()
	return doFromGeoJSON(geom), nil
}

func Decode(data []byte) (geom.Geom, error) {
	var geom Geometry
	if err := json.Unmarshal(data, &geom); err == nil {
		return FromGeoJSON(&geom)
	} else {
		return nil, err
	}
}
