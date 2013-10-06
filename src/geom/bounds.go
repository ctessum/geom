package geom

import (
	"math"
)

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
