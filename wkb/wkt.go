package geom

type WKTGeom interface {
	Geom
	WKT() string
	wktGeometryType() string
}

func (Point) wktGeometryType() string {
	return "POINT"
}

func (PointZ) wktGeometryType() string {
	return "POINTZ"
}

func (PointM) wktGeometryType() string {
	return "POINTM"
}

func (PointZM) wktGeometryType() string {
	return "POINTZM"
}

func (LineString) wktGeometryType() string {
	return "LINESTRING"
}

func (LineStringZ) wktGeometryType() string {
	return "LINESTRINGZ"
}

func (LineStringM) wktGeometryType() string {
	return "LINESTRINGM"
}

func (LineStringZM) wktGeometryType() string {
	return "LINESTRINGZM"
}

func (Polygon) wktGeometryType() string {
	return "POLYGON"
}

func (PolygonZ) wktGeometryType() string {
	return "POLYGONZ"
}

func (PolygonM) wktGeometryType() string {
	return "POLYGONM"
}

func (PolygonZM) wktGeometryType() string {
	return "POLYGONZM"
}
