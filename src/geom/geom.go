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

func (bounds *Bounds) Copy() *Bounds {
	return &Bounds{Point{bounds.Min.X, bounds.Min.Y}, Point{bounds.Max.X, bounds.Max.Y}}
}

func (bounds *Bounds) Empty() bool {
	return bounds.Max.X < bounds.Min.X || bounds.Max.Y < bounds.Min.Y
}

func (bounds *Bounds) ExtendPoint(point Point) *Bounds {
	bounds.Min.X = math.Min(bounds.Min.X, point.X)
	bounds.Min.Y = math.Min(bounds.Min.Y, point.Y)
	bounds.Max.X = math.Max(bounds.Max.X, point.X)
	bounds.Max.Y = math.Max(bounds.Max.Y, point.Y)
	return bounds
}

func (bounds *Bounds) ExtendPointZ(pointZ PointZ) *Bounds {
	bounds.Min.X = math.Min(bounds.Min.X, pointZ.X)
	bounds.Min.Y = math.Min(bounds.Min.Y, pointZ.Y)
	bounds.Max.X = math.Max(bounds.Max.X, pointZ.X)
	bounds.Max.Y = math.Max(bounds.Max.Y, pointZ.Y)
	return bounds
}

func (bounds *Bounds) ExtendPointM(pointM PointM) *Bounds {
	bounds.Min.X = math.Min(bounds.Min.X, pointM.X)
	bounds.Min.Y = math.Min(bounds.Min.Y, pointM.Y)
	bounds.Max.X = math.Max(bounds.Max.X, pointM.X)
	bounds.Max.Y = math.Max(bounds.Max.Y, pointM.Y)
	return bounds
}

func (bounds *Bounds) ExtendPointZM(pointZM PointZM) *Bounds {
	bounds.Min.X = math.Min(bounds.Min.X, pointZM.X)
	bounds.Min.Y = math.Min(bounds.Min.Y, pointZM.Y)
	bounds.Max.X = math.Max(bounds.Max.X, pointZM.X)
	bounds.Max.Y = math.Max(bounds.Max.Y, pointZM.Y)
	return bounds
}

func (bounds *Bounds) ExtendLinearRing(linearRing LinearRing) *Bounds {
	for _, point := range linearRing {
		bounds.ExtendPoint(point)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRingZ(linearRingZ LinearRingZ) *Bounds {
	for _, pointZ := range linearRingZ {
		bounds.ExtendPointZ(pointZ)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRingM(linearRingM LinearRingM) *Bounds {
	for _, pointM := range linearRingM {
		bounds.ExtendPointM(pointM)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRingZM(linearRingZM LinearRingZM) *Bounds {
	for _, pointZM := range linearRingZM {
		bounds.ExtendPointZM(pointZM)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRings(linearRings LinearRings) *Bounds {
	for _, linearRing := range linearRings {
		bounds.ExtendLinearRing(linearRing)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRingZs(linearRingZs LinearRingZs) *Bounds {
	for _, linearRingZ := range linearRingZs {
		bounds.ExtendLinearRingZ(linearRingZ)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRingMs(linearRingMs LinearRingMs) *Bounds {
	for _, linearRingM := range linearRingMs {
		bounds.ExtendLinearRingM(linearRingM)
	}
	return bounds
}

func (bounds *Bounds) ExtendLinearRingZMs(linearRingZMs LinearRingZMs) *Bounds {
	for _, linearRingZM := range linearRingZMs {
		bounds.ExtendLinearRingZM(linearRingZM)
	}
	return bounds
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
