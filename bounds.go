package geom

import (
	"math"
)

type Bounds struct {
	Min, Max Point
}

func (b1 *Bounds) Bounds(b2 *Bounds) *Bounds {
	if b2 == nil {
		return b1
	} else {
		return b1.ExtendPoint(b2.Min).ExtendPoint(b2.Max)
	}
}

func NewBounds() *Bounds {
	return &Bounds{Point{X: math.Inf(1), Y: math.Inf(1)}, Point{X: math.Inf(-1), Y: math.Inf(-1)}}
}

func NewBoundsPoint(point Point) *Bounds {
	return &Bounds{Point{X: point.X, Y: point.Y}, Point{X: point.X, Y: point.Y}}
}

func (b *Bounds) Copy() *Bounds {
	return &Bounds{Point{X: b.Min.X, Y: b.Min.Y}, Point{X: b.Max.X, Y: b.Max.Y}}
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

func (b *Bounds) ExtendPoints(points []Point) *Bounds {
	for _, point := range points {
		b.ExtendPoint(point)
	}
	return b
}

func (b *Bounds) ExtendPointss(pointss [][]Point) *Bounds {
	for _, points := range pointss {
		b.ExtendPoints(points)
	}
	return b
}

func (b1 *Bounds) Overlaps(b2 *Bounds) bool {
	return b1.Min.X <= b2.Max.X && b1.Min.Y <= b2.Max.Y && b1.Max.X >= b2.Min.X && b1.Max.Y >= b2.Min.Y
}
