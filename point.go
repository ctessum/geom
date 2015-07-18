package geom

type Point struct {
	X, Y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{X: x, Y: y}
}

func (point Point) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return NewBoundsPoint(point)
	} else {
		return b.ExtendPoint(point)
	}
}
