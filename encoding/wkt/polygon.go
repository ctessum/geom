package wkt

import (
	"github.com/ctessum/geom"
)

func appendPolygonWKT(dst []byte, polygon geom.Polygon) []byte {
	dst = append(dst, []byte("POLYGON(")...)
	dst = appendPointssCoords(dst, polygon)
	dst = append(dst, ')')
	return dst
}

func appendPolygonZWKT(dst []byte, polygonZ geom.PolygonZ) []byte {
	dst = append(dst, []byte("POLYGONZ(")...)
	dst = appendPointZssCoords(dst, polygonZ)
	dst = append(dst, ')')
	return dst
}

func appendPolygonMWKT(dst []byte, polygonM geom.PolygonM) []byte {
	dst = append(dst, []byte("POLYGONM(")...)
	dst = appendPointMssCoords(dst, polygonM)
	dst = append(dst, ')')
	return dst
}

func appendPolygonZMWKT(dst []byte, polygonZM geom.PolygonZM) []byte {
	dst = append(dst, []byte("POLYGONZM(")...)
	dst = appendPointZMssCoords(dst, polygonZM)
	dst = append(dst, ')')
	return dst
}
