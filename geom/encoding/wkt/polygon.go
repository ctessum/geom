package wkt

import (
	"github.com/twpayne/gogeom/geom"
)

func polygonWKT(polygon geom.Polygon) string {
	return "POLYGON(" + pointssWKTCoordinates(polygon.Rings) + ")"
}

func polygonZWKT(polygonZ geom.PolygonZ) string {
	return "POLYGONZ(" + pointZssWKTCoordinates(polygonZ.Rings) + ")"
}

func polygonMWKT(polygonM geom.PolygonM) string {
	return "POLYGONM(" + pointMssWKTCoordinates(polygonM.Rings) + ")"
}

func polygonZMWKT(polygonZM geom.PolygonZM) string {
	return "POLYGONZM(" + pointZMssWKTCoordinates(polygonZM.Rings) + ")"
}
