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

func appendMultiPolygonZWKT(dst []byte,
	multiPolygonZ geom.MultiPolygonZ) []byte {
	dst = append(dst, []byte("MULTIPOLYGONZ((")...)
	for i, pg := range multiPolygonZ {
		dst = appendPointZssCoords(dst, pg)
		if i != len(multiPolygonZ)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}

func appendMultiPolygonMWKT(dst []byte,
	multiPolygonM geom.MultiPolygonM) []byte {
	dst = append(dst, []byte("MULTIPOLYGONM((")...)
	for i, pg := range multiPolygonM {
		dst = appendPointMssCoords(dst, pg)
		if i != len(multiPolygonM)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}

func appendMultiPolygonZMWKT(dst []byte,
	multiPolygonZM geom.MultiPolygonZM) []byte {
	dst = append(dst, []byte("MULTIPOLYGONZM((")...)
	for i, pg := range multiPolygonZM {
		dst = appendPointZMssCoords(dst, pg)
		if i != len(multiPolygonZM)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}
