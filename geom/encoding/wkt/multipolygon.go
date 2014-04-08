package wkt

import (
	"github.com/twpayne/gogeom/geom"
)

func appendMultiPolygonWKT(dst []byte,
	multiPolygon *geom.MultiPolygon) []byte {
	dst = append(dst, []byte("MULTIPOLYGON((")...)
	for i, pg := range multiPolygon.Polygons {
		dst = appendPointssCoords(dst, pg.Rings)
		if i != len(multiPolygon.Polygons)-1 {
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
	multiPolygonZ *geom.MultiPolygonZ) []byte {
	dst = append(dst, []byte("MULTIPOLYGONZ((")...)
	for i, pg := range multiPolygonZ.Polygons {
		dst = appendPointZssCoords(dst, pg.Rings)
		if i != len(multiPolygonZ.Polygons)-1 {
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
	multiPolygonM *geom.MultiPolygonM) []byte {
	dst = append(dst, []byte("MULTIPOLYGONM((")...)
	for i, pg := range multiPolygonM.Polygons {
		dst = appendPointMssCoords(dst, pg.Rings)
		if i != len(multiPolygonM.Polygons)-1 {
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
	multiPolygonZM *geom.MultiPolygonZM) []byte {
	dst = append(dst, []byte("MULTIPOLYGONZM((")...)
	for i, pg := range multiPolygonZM.Polygons {
		dst = appendPointZMssCoords(dst, pg.Rings)
		if i != len(multiPolygonZM.Polygons)-1 {
			dst = append(dst, ')')
			dst = append(dst, ',')
			dst = append(dst, '(')
		}
	}
	dst = append(dst, ')')
	dst = append(dst, ')')
	return dst
}

