package wkt

import (
	"github.com/ctessum/geom"
)

func appendMultiLineStringWKT(dst []byte,
	multiLineString geom.MultiLineString) []byte {
	dst = append(dst, []byte("MULTILINESTRING((")...)
	for i, ls := range multiLineString {
		dst = appendPointsCoords(dst, ls)
		if i != len(multiLineString)-1 {
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
	multiLineStringZ geom.MultiLineStringZ) []byte {
	dst = append(dst, []byte("MULTILINESTRINGZ((")...)
	for i, ls := range multiLineStringZ {
		dst = appendPointZsCoords(dst, ls)
		if i != len(multiLineStringZ)-1 {
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
	multiLineStringM geom.MultiLineStringM) []byte {
	dst = append(dst, []byte("MULTILINESTRINGM((")...)
	for i, ls := range multiLineStringM {
		dst = appendPointMsCoords(dst, ls)
		if i != len(multiLineStringM)-1 {
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
	multiLineStringZM geom.MultiLineStringZM) []byte {
	dst = append(dst, []byte("MULTILINESTRINGZM((")...)
	for i, ls := range multiLineStringZM {
		dst = appendPointZMsCoords(dst, ls)
		if i != len(multiLineStringZM)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}
