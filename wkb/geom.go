package wkb

type Geom interface{}

type Point struct {
	X float64
	Y float64
}

type PointZ struct {
	X float64
	Y float64
	Z float64
}

type PointM struct {
	X float64
	Y float64
	M float64
}

type PointZM struct {
	X float64
	Y float64
	Z float64
	M float64
}

type LinearRing []Point
type LinearRingZ []PointZ
type LinearRingM []PointM
type LinearRingZM []PointZM

type LineString struct {
	Points LinearRing
}

type LineStringZ struct {
	Points LinearRingZ
}

type LineStringM struct {
	Points LinearRingM
}

type LineStringZM struct {
	Points LinearRingZM
}

type Polygon struct {
	Rings []LinearRing
}

type PolygonZ struct {
	Rings []LinearRingZ
}

type PolygonM struct {
	Rings []LinearRingM
}

type PolygonZM struct {
	Rings []LinearRingZM
}
