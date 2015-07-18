package wkt

import (
	"github.com/ctessum/geom"
)

func appendMultiPolygonWKT(dst []byte,
	multiPolygon geom.MultiPolygon) []byte {
	dst = append(dst, []byte("MULTIPOLYGON((")...)
	for i, pg := range multiPolygon {
		dst = appendPointssCoords(dst, pg)
		if i != len(multiPolygon)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}
