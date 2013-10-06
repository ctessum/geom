package geom

import (
	"fmt"
	"strings"
)

func (point Point) wktCoordinates() string {
	return fmt.Sprintf("%g %g", point.X, point.Y)
}

func (pointZ PointZ) wktCoordinates() string {
	return fmt.Sprintf("%g %g %g", pointZ.X, pointZ.Y, pointZ.Z)
}

func (pointM PointM) wktCoordinates() string {
	return fmt.Sprintf("%g %g %g", pointM.X, pointM.Y, pointM.M)
}

func (pointZM PointZM) wktCoordinates() string {
	return fmt.Sprintf("%g %g %g %g", pointZM.X, pointZM.Y, pointZM.Z, pointZM.M)
}

func (linearRing LinearRing) wktCoordinates() string {
	coordinates := make([]string, len(linearRing))
	for i, point := range linearRing {
		coordinates[i] = point.wktCoordinates()
	}
	return strings.Join(coordinates, ",")
}

func (linearRingZ LinearRingZ) wktCoordinates() string {
	coordinates := make([]string, len(linearRingZ))
	for i, pointZ := range linearRingZ {
		coordinates[i] = pointZ.wktCoordinates()
	}
	return strings.Join(coordinates, ",")
}

func (linearRingM LinearRingM) wktCoordinates() string {
	coordinates := make([]string, len(linearRingM))
	for i, pointM := range linearRingM {
		coordinates[i] = pointM.wktCoordinates()
	}
	return strings.Join(coordinates, ",")
}

func (linearRingZM LinearRingZM) wktCoordinates() string {
	coordinates := make([]string, len(linearRingZM))
	for i, pointZM := range linearRingZM {
		coordinates[i] = pointZM.wktCoordinates()
	}
	return strings.Join(coordinates, ",")
}

func (linearRings LinearRings) wktCoordinates() string {
	coordinates := make([]string, len(linearRings))
	for i, linearRing := range linearRings {
		coordinates[i] = "(" + linearRing.wktCoordinates() + ")"
	}
	return strings.Join(coordinates, ",")
}

func (linearRingZs LinearRingZs) wktCoordinates() string {
	coordinates := make([]string, len(linearRingZs))
	for i, linearRingZ := range linearRingZs {
		coordinates[i] = "(" + linearRingZ.wktCoordinates() + ")"
	}
	return strings.Join(coordinates, ",")
}

func (linearRingMs LinearRingMs) wktCoordinates() string {
	coordinates := make([]string, len(linearRingMs))
	for i, linearRingM := range linearRingMs {
		coordinates[i] = "(" + linearRingM.wktCoordinates() + ")"
	}
	return strings.Join(coordinates, ",")
}

func (linearRingZMs LinearRingZMs) wktCoordinates() string {
	coordinates := make([]string, len(linearRingZMs))
	for i, linearRingZM := range linearRingZMs {
		coordinates[i] = "(" + linearRingZM.wktCoordinates() + ")"
	}
	return strings.Join(coordinates, ",")
}

func (point Point) WKT() string {
	return point.wktGeometryType() + "(" + point.wktCoordinates() + ")"
}

func (pointZ PointZ) WKT() string {
	return pointZ.wktGeometryType() + "(" + pointZ.wktCoordinates() + ")"
}

func (pointM PointM) WKT() string {
	return pointM.wktGeometryType() + "(" + pointM.wktCoordinates() + ")"
}

func (pointZM PointZM) WKT() string {
	return pointZM.wktGeometryType() + "(" + pointZM.wktCoordinates() + ")"
}

func (lineString LineString) WKT() string {
	return lineString.wktGeometryType() + "(" + lineString.Points.wktCoordinates() + ")"
}

func (lineStringZ LineStringZ) WKT() string {
	return lineStringZ.wktGeometryType() + "(" + lineStringZ.Points.wktCoordinates() + ")"
}

func (lineStringM LineStringM) WKT() string {
	return lineStringM.wktGeometryType() + "(" + lineStringM.Points.wktCoordinates() + ")"
}

func (lineStringZM LineStringZM) WKT() string {
	return lineStringZM.wktGeometryType() + "(" + lineStringZM.Points.wktCoordinates() + ")"
}

func (polygon Polygon) WKT() string {
	ringWKTs := make([]string, len(polygon.Rings))
	for i, ring := range polygon.Rings {
		ringWKTs[i] = "(" + ring.wktCoordinates() + ")"
	}
	return polygon.wktGeometryType() + "(" + strings.Join(ringWKTs, ",") + ")"
}

func (polygonZ PolygonZ) WKT() string {
	ringWKTs := make([]string, len(polygonZ.Rings))
	for i, ringZ := range polygonZ.Rings {
		ringWKTs[i] = "(" + ringZ.wktCoordinates() + ")"
	}
	return polygonZ.wktGeometryType() + "(" + strings.Join(ringWKTs, ",") + ")"
}

func (polygonM PolygonM) WKT() string {
	ringWKTs := make([]string, len(polygonM.Rings))
	for i, ringM := range polygonM.Rings {
		ringWKTs[i] = "(" + ringM.wktCoordinates() + ")"
	}
	return polygonM.wktGeometryType() + "(" + strings.Join(ringWKTs, ",") + ")"
}

func (polygonZM PolygonZM) WKT() string {
	ringWKTs := make([]string, len(polygonZM.Rings))
	for i, ringZM := range polygonZM.Rings {
		ringWKTs[i] = "(" + ringZM.wktCoordinates() + ")"
	}
	return polygonZM.wktGeometryType() + "(" + strings.Join(ringWKTs, ",") + ")"
}
