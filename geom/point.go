package geom

type Point struct {
	X float64
	Y float64
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

func (pointZM PointZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		return NewBoundsPointZM(pointZM)
	} else {
		return b.ExtendPointZM(pointZM)
	}
}
