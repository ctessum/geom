package wkt

import (
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"strings"
)

func pointWKTCoordinates(point geom.Point) string {
        return fmt.Sprintf("%g %g", point.X, point.Y)
}

func pointsWKCoordinates(points []geom.Point) string {
        wktCoordinates := make([]string, len(points))
        for i, point := range points {
                wktCoordinates[i] = pointWKTCoordinates(point)
        }
        return strings.Join(wktCoordinates, ",")
}

func pointssWKTCoordinates(pointss [][]geom.Point) string {
        wktCoordinates := make([]string, len(pointss))
        for i, points := range pointss {
                wktCoordinates[i] = "(" + pointsWKCoordinates(points) + ")"
        }
        return strings.Join(wktCoordinates, ",")
}

func pointWKT(point geom.Point) string {
        return "POINT(" + pointWKTCoordinates(point) + ")"
}

func pointZWKTCoordinates(pointZ geom.PointZ) string {
        return fmt.Sprintf("%g %g %g", pointZ.X, pointZ.Y, pointZ.Z)
}

func pointZsWKCoordinates(pointZs []geom.PointZ) string {
        wktCoordinates := make([]string, len(pointZs))
        for i, pointZ := range pointZs {
                wktCoordinates[i] = pointZWKTCoordinates(pointZ)
        }
        return strings.Join(wktCoordinates, ",")
}

func pointZssWKTCoordinates(pointZss [][]geom.PointZ) string {
        wktCoordinates := make([]string, len(pointZss))
        for i, pointZs := range pointZss {
                wktCoordinates[i] = "(" + pointZsWKCoordinates(pointZs) + ")"
        }
        return strings.Join(wktCoordinates, ",")
}

func pointZWKT(pointZ geom.PointZ) string {
        return "POINTZ(" + pointZWKTCoordinates(pointZ) + ")"
}

func pointMWKTCoordinates(pointM geom.PointM) string {
        return fmt.Sprintf("%g %g %g", pointM.X, pointM.Y, pointM.M)
}

func pointMsWKCoordinates(pointMs []geom.PointM) string {
        wktCoordinates := make([]string, len(pointMs))
        for i, pointM := range pointMs {
                wktCoordinates[i] = pointMWKTCoordinates(pointM)
        }
        return strings.Join(wktCoordinates, ",")
}

func pointMssWKTCoordinates(pointMss [][]geom.PointM) string {
        wktCoordinates := make([]string, len(pointMss))
        for i, pointMs := range pointMss {
                wktCoordinates[i] = "(" + pointMsWKCoordinates(pointMs) + ")"
        }
        return strings.Join(wktCoordinates, ",")
}

func pointMWKT(pointM geom.PointM) string {
        return "POINTM(" + pointMWKTCoordinates(pointM) + ")"
}

func pointZMWKTCoordinates(pointZM geom.PointZM) string {
        return fmt.Sprintf("%g %g %g %g", pointZM.X, pointZM.Y, pointZM.Z, pointZM.M)
}

func pointZMsWKCoordinates(pointZMs []geom.PointZM) string {
        wktCoordinates := make([]string, len(pointZMs))
        for i, pointZM := range pointZMs {
                wktCoordinates[i] = pointZMWKTCoordinates(pointZM)
        }
        return strings.Join(wktCoordinates, ",")
}

func pointZMssWKTCoordinates(pointZMss [][]geom.PointZM) string {
        wktCoordinates := make([]string, len(pointZMss))
        for i, pointZMs := range pointZMss {
                wktCoordinates[i] = "(" + pointZMsWKCoordinates(pointZMs) + ")"
        }
        return strings.Join(wktCoordinates, ",")
}

func pointZMWKT(pointZM geom.PointZM) string {
        return "POINTZM(" + pointZMWKTCoordinates(pointZM) + ")"
}
