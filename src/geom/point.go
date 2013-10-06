package geom

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
