package projgeom

import (
	"github.com/ctessum/geom"
	"github.com/pebbe/go-proj-4/proj"
)

func projectPoint(point *geom.Point, src, dst *proj.Proj, inputDegrees,
	outputDegrees bool) (geom.T, error) {
	var newPoint geom.Point
	var err error
	if inputDegrees {
		newPoint.X, newPoint.Y, err = proj.Transform2(src, dst,
			proj.DegToRad(point.X), proj.DegToRad(point.Y))
	} else {
		newPoint.X, newPoint.Y, err = proj.Transform2(src, dst,
			point.X, point.Y)
	}
	if outputDegrees {
		newPoint.X, newPoint.Y =
			proj.RadToDeg(newPoint.X), proj.RadToDeg(newPoint.Y)
	}
	return newPoint, err
}

func projectLineString(lineString geom.LineString, src, dst *proj.Proj,
	inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newPoint geom.T
	var newLineString geom.LineString
	newLineString = make([]geom.Point, len(lineString))
	for i, point := range lineString {
		newPoint, err = projectPoint(&point, src, dst, inputDegrees,
			outputDegrees)
		if err != nil {
			return nil, err
		}
		newLineString[i] = newPoint.(geom.Point)
	}
	return newLineString, err
}

func projectMultiLineString(multiLineString geom.MultiLineString,
	src, dst *proj.Proj, inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newLineString geom.T
	var newMultiLineString geom.MultiLineString
	newMultiLineString =
		make([]geom.LineString, len(multiLineString))
	for i, lineString := range multiLineString {
		newLineString, err = projectLineString(lineString,
			src, dst, inputDegrees, outputDegrees)
		if err != nil {
			return nil, err
		}
		newMultiLineString[i] = newLineString.(geom.LineString)
	}
	return newMultiLineString, err
}

func projectPolygon(polygon geom.Polygon, src, dst *proj.Proj,
	inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newPoint geom.T
	var newPolygon geom.Polygon
	newPolygon = make([][]geom.Point, len(polygon))
	for i := 0; i < len(polygon); i++ {
		newPolygon[i] = make([]geom.Point, len(polygon[i]))
		for j := 0; j < len(polygon[i]); j++ {
			newPoint, err = projectPoint(&polygon[i][j],
				src, dst, inputDegrees, outputDegrees)
			if err != nil {
				return nil, err
			}
			newPolygon[i][j] = newPoint.(geom.Point)
		}
	}
	return newPolygon, err
}

func projectMultiPolygon(multiPolygon geom.MultiPolygon,
	src, dst *proj.Proj, inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newPolygon geom.T
	var newMultiPolygon geom.MultiPolygon
	newMultiPolygon = make([]geom.Polygon, len(multiPolygon))
	for i, polygon := range multiPolygon {
		newPolygon, err = projectPolygon(polygon,
			src, dst, inputDegrees, outputDegrees)
		if err != nil {
			return nil, err
		}
		multiPolygon[i] = newPolygon.(geom.Polygon)
	}
	return newMultiPolygon, err
}
