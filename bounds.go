package geom

import (
	"math"
)

type Bounds struct {
	Min, Max Point
}

func (b1 *Bounds) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return b1
	} else {
		return b1.ExtendPoint(b.Min).ExtendPoint(b.Max)
	}
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

func (b *Bounds) ExtendPoints(points []Point) *Bounds {
	for _, point := range points {
		b.ExtendPoint(point)
	}
	return b
}

func (b *Bounds) ExtendPointZs(pointZs []PointZ) *Bounds {
	for _, pointZ := range pointZs {
		b.ExtendPointZ(pointZ)
	}
	return b
}

func (b *Bounds) ExtendPointMs(pointMs []PointM) *Bounds {
	for _, pointM := range pointMs {
		b.ExtendPointM(pointM)
	}
	return b
}

func (b *Bounds) ExtendPointZMs(pointZMs []PointZM) *Bounds {
	for _, pointZM := range pointZMs {
		b.ExtendPointZM(pointZM)
	}
	return b
}

func (b *Bounds) ExtendPointss(pointss [][]Point) *Bounds {
	for _, points := range pointss {
		b.ExtendPoints(points)
	}
	return b
}

func (b *Bounds) ExtendPointZss(pointZss [][]PointZ) *Bounds {
	for _, pointZs := range pointZss {
		b.ExtendPointZs(pointZs)
	}
	return b
}

func (b *Bounds) ExtendPointMss(pointMss [][]PointM) *Bounds {
	for _, pointMs := range pointMss {
		b.ExtendPointMs(pointMs)
	}
	return b
}

func (b *Bounds) ExtendPointZMss(pointZMss [][]PointZM) *Bounds {
	for _, pointZMs := range pointZMss {
		b.ExtendPointZMs(pointZMs)
	}
	return b
}

func (b1 *Bounds) Overlaps(b2 *Bounds) bool {
	return b1.Min.X <= b2.Max.X && b1.Min.Y <= b2.Max.Y && b1.Max.X >= b2.Min.X && b1.Max.Y >= b2.Min.Y
}
