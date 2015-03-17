package geom

type Point struct {
	X float64
	Y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func (point Point) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return NewBoundsPoint(point)
	} else {
		return b.ExtendPoint(point)
	}
}

type PointZ struct {
	X float64
	Y float64
	Z float64
}

func NewPointZ(x, y, z float64) *PointZ {
	return &PointZ{x, y, z}
}

func (pointZ PointZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return NewBoundsPointZ(pointZ)
	} else {
		return b.ExtendPointZ(pointZ)
	}
}

type PointM struct {
	X float64
	Y float64
	M float64
}

func NewPointM(x, y, m float64) *PointM {
	return &PointM{x, y, m}
}

func (pointM PointM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return NewBoundsPointM(pointM)
	} else {
		return b.ExtendPointM(pointM)
	}
}

type PointZM struct {
	X float64
	Y float64
	Z float64
	M float64
}

func NewPointZM(x, y, z, m float64) *PointZM {
	return &PointZM{x, y, z, m}
}

func (pointZM PointZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return NewBoundsPointZM(pointZM)
	} else {
		return b.ExtendPointZM(pointZM)
	}
}
