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

func appendPointssCoords(dst []byte, pointss [][]geom.Point) []byte {
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

func appendPointZCoords(dst []byte, pointZ *geom.PointZ) []byte {
	dst = strconv.AppendFloat(dst, pointZ.X, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointZ.Y, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointZ.Z, 'g', -1, 64)
	return dst
}

func appendPointZsCoords(dst []byte, pointZs []geom.PointZ) []byte {
	for i, pointZ := range pointZs {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = appendPointZCoords(dst, &pointZ)
	}
	return dst
}

func appendPointZssCoords(dst []byte, pointZss [][]geom.PointZ) []byte {
	for i, pointZs := range pointZss {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '(')
		dst = appendPointZsCoords(dst, pointZs)
		dst = append(dst, ')')
	}
	return dst
}

func appendPointZWKT(dst []byte, pointZ *geom.PointZ) []byte {
	dst = append(dst, []byte("POINTZ(")...)
	dst = appendPointZCoords(dst, pointZ)
	dst = append(dst, ')')
	return dst
}

func appendPointMCoords(dst []byte, pointM *geom.PointM) []byte {
	dst = strconv.AppendFloat(dst, pointM.X, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointM.Y, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointM.M, 'g', -1, 64)
	return dst
}

func appendPointMsCoords(dst []byte, pointMs []geom.PointM) []byte {
	for i, pointM := range pointMs {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = appendPointMCoords(dst, &pointM)
	}
	return dst
}

func appendPointMssCoords(dst []byte, pointMss [][]geom.PointM) []byte {
	for i, pointMs := range pointMss {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '(')
		dst = appendPointMsCoords(dst, pointMs)
		dst = append(dst, ')')
	}
	return dst
}

func appendPointMWKT(dst []byte, pointM *geom.PointM) []byte {
	dst = append(dst, []byte("POINTM(")...)
	dst = appendPointMCoords(dst, pointM)
	dst = append(dst, ')')
	return dst
}

func appendPointZMCoords(dst []byte, pointZM *geom.PointZM) []byte {
	dst = strconv.AppendFloat(dst, pointZM.X, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointZM.Y, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointZM.Z, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, pointZM.M, 'g', -1, 64)
	return dst
}

func appendPointZMsCoords(dst []byte, pointZMs []geom.PointZM) []byte {
	for i, pointZM := range pointZMs {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = appendPointZMCoords(dst, &pointZM)
	}
	return dst
}

func appendPointZMssCoords(dst []byte, pointZMss [][]geom.PointZM) []byte {
	for i, pointZMs := range pointZMss {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '(')
		dst = appendPointZMsCoords(dst, pointZMs)
		dst = append(dst, ')')
	}
	return dst
}

func appendPointZMWKT(dst []byte, pointZM *geom.PointZM) []byte {
	dst = append(dst, []byte("POINTZM(")...)
	dst = appendPointZMCoords(dst, pointZM)
	dst = append(dst, ')')
	return dst
}
