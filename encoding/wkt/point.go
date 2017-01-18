package wkt

import (
	"github.com/ctessum/geom"
	"strconv"
)

func appendPointCoords(dst []byte, point *geom.Point) []byte {
	dst = strconv.AppendFloat(dst, point.X, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, point.Y, 'g', -1, 64)
	return dst
}

func appendPointsCoords(dst []byte, points []geom.Point) []byte {
	for i, point := range points {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = appendPointCoords(dst, &point)
	}
	return dst
}

func appendPointssCoords(dst []byte, pointss []geom.Path) []byte {
	for i, points := range pointss {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '(')
		dst = appendPointsCoords(dst, points)
		dst = append(dst, ')')
	}
	return dst
}

func appendPointWKT(dst []byte, point *geom.Point) []byte {
	dst = append(dst, []byte("POINT(")...)
	dst = appendPointCoords(dst, point)
	dst = append(dst, ')')
	return dst
}
