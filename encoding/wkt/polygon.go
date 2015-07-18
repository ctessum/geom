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
