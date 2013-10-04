package wkb

type Geom interface {
}

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

type LineString struct {
	Points []Point
}

type LineStringZ struct {
	Points []PointZ
}

type LineStringM struct {
	Points []PointM
}

type LineStringZM struct {
	Points []PointZM
}
