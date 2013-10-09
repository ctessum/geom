package wkt

import (
	"github.com/twpayne/gogeom/geom"
)

func appendPolygonWKT(dst []byte, polygon *geom.Polygon) []byte {
	dst = append(dst, []byte("POLYGON(")...)
	dst = appendPointssCoords(dst, polygon.Rings)
	dst = append(dst, ')')
	return dst
}

func appendPolygonZWKT(dst []byte, polygonZ *geom.PolygonZ) []byte {
	dst = append(dst, []byte("POLYGONZ(")...)
	dst = appendPointZssCoords(dst, polygonZ.Rings)
	dst = append(dst, ')')
	return dst
}

func appendPolygonMWKT(dst []byte, polygonM *geom.PolygonM) []byte {
	dst = append(dst, []byte("POLYGONM(")...)
	dst = appendPointMssCoords(dst, polygonM.Rings)
	dst = append(dst, ')')
	return dst
}

func appendPolygonZMWKT(dst []byte, polygonZM *geom.PolygonZM) []byte {
	dst = append(dst, []byte("POLYGONZM(")...)
	dst = appendPointZMssCoords(dst, polygonZM.Rings)
	dst = append(dst, ')')
	return dst
}
