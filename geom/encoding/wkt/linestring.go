package wkt

import (
	"github.com/twpayne/gogeom/geom"
)

func appendLineStringWKT(dst []byte, lineString *geom.LineString) []byte {
	dst = append(dst, []byte("LINESTRING(")...)
	dst = appendPointsCoords(dst, lineString.Points)
	dst = append(dst, ')')
	return dst
}

func appendLineStringZWKT(dst []byte, lineStringZ *geom.LineStringZ) []byte {
	dst = append(dst, []byte("LINESTRINGZ(")...)
	dst = appendPointZsCoords(dst, lineStringZ.Points)
	dst = append(dst, ')')
	return dst
}

func appendLineStringMWKT(dst []byte, lineStringM *geom.LineStringM) []byte {
	dst = append(dst, []byte("LINESTRINGM(")...)
	dst = appendPointMsCoords(dst, lineStringM.Points)
	dst = append(dst, ')')
	return dst
}

func appendLineStringZMWKT(dst []byte, lineStringZM *geom.LineStringZM) []byte {
	dst = append(dst, []byte("LINESTRINGZM(")...)
	dst = appendPointZMsCoords(dst, lineStringZM.Points)
	dst = append(dst, ')')
	return dst
}
