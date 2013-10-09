package wkt

import (
	"github.com/twpayne/gogeom/geom"
)

func lineStringWKT(lineString geom.LineString) string {
	return "LINESTRING(" + pointsWKCoordinates(lineString.Points) + ")"
}

func lineStringZWKT(lineStringZ geom.LineStringZ) string {
	return "LINESTRINGZ(" + pointZsWKCoordinates(lineStringZ.Points) + ")"
}

func lineStringMWKT(lineStringM geom.LineStringM) string {
	return "LINESTRINGM(" + pointMsWKCoordinates(lineStringM.Points) + ")"
}

func lineStringZMWKT(lineStringZM geom.LineStringZM) string {
	return "LINESTRINGZM(" + pointZMsWKCoordinates(lineStringZM.Points) + ")"
}
