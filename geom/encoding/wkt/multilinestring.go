package wkt

import (
	"github.com/twpayne/gogeom/geom"
)

func appendMultiLineStringWKT(dst []byte,
	multiLineString *geom.MultiLineString) []byte {
	dst = append(dst, []byte("MULTILINESTRING((")...)
	for i, ls := range multiLineString.LineStrings {
		dst = appendPointsCoords(dst, ls.Points)
		if i != len(multiLineString.LineStrings)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}

func appendMultiLineStringZWKT(dst []byte,
	multiLineStringZ *geom.MultiLineStringZ) []byte {
	dst = append(dst, []byte("MULTILINESTRINGZ((")...)
	for i, ls := range multiLineStringZ.LineStrings {
		dst = appendPointZsCoords(dst, ls.Points)
		if i != len(multiLineStringZ.LineStrings)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}

func appendMultiLineStringMWKT(dst []byte,
	multiLineStringM *geom.MultiLineStringM) []byte {
	dst = append(dst, []byte("MULTILINESTRINGM((")...)
	for i, ls := range multiLineStringM.LineStrings {
		dst = appendPointMsCoords(dst, ls.Points)
		if i != len(multiLineStringM.LineStrings)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}

func appendMultiLineStringZMWKT(dst []byte,
	multiLineStringZM *geom.MultiLineStringZM) []byte {
	dst = append(dst, []byte("MULTILINESTRINGZM((")...)
	for i, ls := range multiLineStringZM.LineStrings {
		dst = appendPointZMsCoords(dst, ls.Points)
		if i != len(multiLineStringZM.LineStrings)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}
