package geom

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

type LinearRings []LinearRing
type LinearRingZs []LinearRingZ
type LinearRingMs []LinearRingM
type LinearRingZMs []LinearRingZM

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
	Rings LinearRings
}

type PolygonZ struct {
	Rings LinearRingZs
}

type PolygonM struct {
	Rings LinearRingMs
}

type PolygonZM struct {
	Rings LinearRingZMs
}
