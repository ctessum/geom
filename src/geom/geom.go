package geom

import (
	"math"
)

type T interface {
	Bounds() *Bounds
}

type Bounds struct {
	Min, Max Point
}

func NewBounds() *Bounds {
	return &Bounds{Point{math.Inf(1), math.Inf(1)}, Point{math.Inf(-1), math.Inf(-1)}}
}

func NewBoundsPoint(point Point) *Bounds {
	return &Bounds{Point{point.X, point.Y}, Point{point.X, point.Y}}
}

func NewBoundsPointZ(pointZ PointZ) *Bounds {
	return &Bounds{Point{pointZ.X, pointZ.Y}, Point{pointZ.X, pointZ.Y}}
}

func NewBoundsPointM(pointM PointM) *Bounds {
	return &Bounds{Point{pointM.X, pointM.Y}, Point{pointM.X, pointM.Y}}
}

func NewBoundsPointZM(pointZM PointZM) *Bounds {
	return &Bounds{Point{pointZM.X, pointZM.Y}, Point{pointZM.X, pointZM.Y}}
}

func (b *Bounds) Copy() *Bounds {
	return &Bounds{Point{b.Min.X, b.Min.Y}, Point{b.Max.X, b.Max.Y}}
}

func (b *Bounds) Empty() bool {
	return b.Max.X < b.Min.X || b.Max.Y < b.Min.Y
}

func (b *Bounds) ExtendPoint(point Point) *Bounds {
	b.Min.X = math.Min(b.Min.X, point.X)
	b.Min.Y = math.Min(b.Min.Y, point.Y)
	b.Max.X = math.Max(b.Max.X, point.X)
	b.Max.Y = math.Max(b.Max.Y, point.Y)
	return b
}

func (b *Bounds) ExtendPointZ(pointZ PointZ) *Bounds {
	b.Min.X = math.Min(b.Min.X, pointZ.X)
	b.Min.Y = math.Min(b.Min.Y, pointZ.Y)
	b.Max.X = math.Max(b.Max.X, pointZ.X)
	b.Max.Y = math.Max(b.Max.Y, pointZ.Y)
	return b
}

func (b *Bounds) ExtendPointM(pointM PointM) *Bounds {
	b.Min.X = math.Min(b.Min.X, pointM.X)
	b.Min.Y = math.Min(b.Min.Y, pointM.Y)
	b.Max.X = math.Max(b.Max.X, pointM.X)
	b.Max.Y = math.Max(b.Max.Y, pointM.Y)
	return b
}

func (b *Bounds) ExtendPointZM(pointZM PointZM) *Bounds {
	b.Min.X = math.Min(b.Min.X, pointZM.X)
	b.Min.Y = math.Min(b.Min.Y, pointZM.Y)
	b.Max.X = math.Max(b.Max.X, pointZM.X)
	b.Max.Y = math.Max(b.Max.Y, pointZM.Y)
	return b
}

func (b *Bounds) ExtendLinearRing(linearRing LinearRing) *Bounds {
	for _, point := range linearRing {
		b.ExtendPoint(point)
	}
	return b
}

func (b *Bounds) ExtendLinearRingZ(linearRingZ LinearRingZ) *Bounds {
	for _, pointZ := range linearRingZ {
		b.ExtendPointZ(pointZ)
	}
	return b
}

func (b *Bounds) ExtendLinearRingM(linearRingM LinearRingM) *Bounds {
	for _, pointM := range linearRingM {
		b.ExtendPointM(pointM)
	}
	return b
}

func (b *Bounds) ExtendLinearRingZM(linearRingZM LinearRingZM) *Bounds {
	for _, pointZM := range linearRingZM {
		b.ExtendPointZM(pointZM)
	}
	return b
}

func (b *Bounds) ExtendLinearRings(linearRings LinearRings) *Bounds {
	for _, linearRing := range linearRings {
		b.ExtendLinearRing(linearRing)
	}
	return b
}

func (b *Bounds) ExtendLinearRingZs(linearRingZs LinearRingZs) *Bounds {
	for _, linearRingZ := range linearRingZs {
		b.ExtendLinearRingZ(linearRingZ)
	}
	return b
}

func (b *Bounds) ExtendLinearRingMs(linearRingMs LinearRingMs) *Bounds {
	for _, linearRingM := range linearRingMs {
		b.ExtendLinearRingM(linearRingM)
	}
	return b
}

func (b *Bounds) ExtendLinearRingZMs(linearRingZMs LinearRingZMs) *Bounds {
	for _, linearRingZM := range linearRingZMs {
		b.ExtendLinearRingZM(linearRingZM)
	}
	return b
}

func (b1 *Bounds) Overlaps(b2 *Bounds) bool {
	return b1.Min.X <= b2.Max.X && b1.Min.Y <= b2.Max.Y && b1.Max.X >= b2.Min.X && b1.Max.Y >= b2.Min.Y
}

type Point struct {
	X float64
	Y float64
}

func (point Point) Bounds() *Bounds {
	return NewBoundsPoint(point)
}

type PointZ struct {
	X float64
	Y float64
	Z float64
}

func (pointZ PointZ) Bounds() *Bounds {
	return NewBoundsPointZ(pointZ)
}

type PointM struct {
	X float64
	Y float64
	M float64
}

func (pointM PointM) Bounds() *Bounds {
	return NewBoundsPointM(pointM)
}

type PointZM struct {
	X float64
	Y float64
	Z float64
	M float64
}

func (pointZM PointZM) Bounds() *Bounds {
	return NewBoundsPointZM(pointZM)
}

type LinearRing []Point

func (linearRing LinearRing) Bounds() *Bounds {
	return NewBounds().ExtendLinearRing(linearRing)
}

type LinearRingZ []PointZ

func (linearRingZ LinearRingZ) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZ(linearRingZ)
}

type LinearRingM []PointM

func (linearRingM LinearRingM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingM(linearRingM)
}

type LinearRingZM []PointZM

func (linearRingZM LinearRingZM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZM(linearRingZM)
}

type LinearRings []LinearRing
type LinearRingZs []LinearRingZ
type LinearRingMs []LinearRingM
type LinearRingZMs []LinearRingZM

type LineString struct {
	Points LinearRing
}

func (lineString LineString) Bounds() *Bounds {
	return lineString.Points.Bounds()
}

type LineStringZ struct {
	Points LinearRingZ
}

func (lineStringZ LineStringZ) Bounds() *Bounds {
	return lineStringZ.Points.Bounds()
}

type LineStringM struct {
	Points LinearRingM
}

func (lineStringM LineStringM) Bounds() *Bounds {
	return lineStringM.Points.Bounds()
}

type LineStringZM struct {
	Points LinearRingZM
}

func (lineStringZM LineStringZM) Bounds() *Bounds {
	return lineStringZM.Points.Bounds()
}

type Polygon struct {
	Rings LinearRings
}

func (polygon Polygon) Bounds() *Bounds {
	return NewBounds().ExtendLinearRings(polygon.Rings)
}

type PolygonZ struct {
	Rings LinearRingZs
}

func (polygonZ PolygonZ) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZs(polygonZ.Rings)
}

type PolygonM struct {
	Rings LinearRingMs
}

func (polygonM PolygonM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingMs(polygonM.Rings)
}

type PolygonZM struct {
	Rings LinearRingZMs
}

func (polygonZM PolygonZM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZMs(polygonZM.Rings)
}
