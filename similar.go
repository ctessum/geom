package geom

import (
	"math"
	"reflect"
)

func similar(a, b, e float64) bool {
	return math.Abs(a-b) < e
}

func pointSimilar(p1, p2 Point, e float64) bool {
	return similar(p1.X, p2.X, e) && similar(p1.Y, p2.Y, e)
}

func pointsSimilar(p1s, p2s []Point, e float64) bool {
	if len(p1s) != len(p2s) {
		return false
	}
	for i, n := 0, len(p1s); i < n; i++ {
		if !pointSimilar(p1s[i], p2s[i], e) {
			return false
		}
	}
	return true
}

func pointssSimilar(p1ss, p2ss [][]Point, e float64) bool {
	if len(p1ss) != len(p2ss) {
		return false
	}
	for i, n := 0, len(p1ss); i < n; i++ {
		if !pointsSimilar(p1ss[i], p2ss[i], e) {
			return false
		}
	}
	return true
}

func Similar(g1, g2 T, e float64) bool {
	if reflect.TypeOf(g1) != reflect.TypeOf(g2) {
		return false
	}
	switch g1.(type) {
	case Point:
		return pointSimilar(g1.(Point), g2.(Point), e)
	case LineString:
		return pointsSimilar(g1.(LineString), g2.(LineString), e)
	case Polygon:
		return pointssSimilar(g1.(Polygon), g2.(Polygon), e)
	case MultiPoint:
		return pointsSimilar(g1.(MultiPoint), g2.(MultiPoint), e)
	default:
		return false
	}
}
