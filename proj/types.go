package projgeom

import (
	"github.com/pebbe/go-proj-4/proj"
	"github.com/twpayne/gogeom/geom"
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

func projectLineString(lineString *geom.LineString, src, dst *proj.Proj,
	inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newPoint geom.T
	var newLineString geom.LineString
	newLineString.Points = make([]geom.Point, len(lineString.Points))
	for i, point := range lineString.Points {
		newPoint, err = projectPoint(&point, src, dst, inputDegrees,
			outputDegrees)
		if err != nil {
			return nil, err
		}
		newLineString.Points[i] = newPoint.(geom.Point)
	}
	return newLineString, err
}

func projectMultiLineString(multiLineString *geom.MultiLineString,
	src, dst *proj.Proj, inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newLineString geom.T
	var newMultiLineString geom.MultiLineString
	newMultiLineString.LineStrings =
		make([]geom.LineString, len(multiLineString.LineStrings))
	for i, lineString := range multiLineString.LineStrings {
		newLineString, err = projectLineString(&lineString,
			src, dst, inputDegrees, outputDegrees)
		if err != nil {
			return nil, err
		}
		newMultiLineString.LineStrings[i] = newLineString.(geom.LineString)
	}
	return newMultiLineString, err
}

func projectPolygon(polygon *geom.Polygon, src, dst *proj.Proj,
	inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newPoint geom.T
	var newPolygon geom.Polygon
	newPolygon.Rings = make([][]geom.Point, len(polygon.Rings))
	for i := 0; i < len(polygon.Rings); i++ {
		newPolygon.Rings[i] = make([]geom.Point, len(polygon.Rings[i]))
		for j := 0; j < len(polygon.Rings[i]); j++ {
			newPoint, err = projectPoint(&polygon.Rings[i][j],
				src, dst, inputDegrees, outputDegrees)
			if err != nil {
				return nil, err
			}
			newPolygon.Rings[i][j] = newPoint.(geom.Point)
		}
	}
	return newPolygon, err
}

func projectMultiPolygon(multiPolygon *geom.MultiPolygon,
	src, dst *proj.Proj, inputDegrees, outputDegrees bool) (geom.T, error) {
	var err error
	var newPolygon geom.T
	var newMultiPolygon geom.MultiPolygon
	newMultiPolygon.Polygons = make([]geom.Polygon, len(multiPolygon.Polygons))
	for i, polygon := range multiPolygon.Polygons {
		newPolygon, err = projectPolygon(&polygon,
			src, dst, inputDegrees, outputDegrees)
		if err != nil {
			return nil, err
		}
		multiPolygon.Polygons[i] = newPolygon.(geom.Polygon)
	}
	return newMultiPolygon, err
}
