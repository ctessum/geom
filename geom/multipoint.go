package geom

type MultiPoint struct {
	Points []Point
}

func (multiPoint MultiPoint) Bounds() *Bounds {
	return NewBounds().ExtendPoints(multiPoint.Points)
}

type MultiPointZ struct {
	Points []PointZ
}

func (multiPointZ MultiPointZ) Bounds() *Bounds {
	return NewBounds().ExtendPointZs(multiPointZ.Points)
}

type MultiPointM struct {
	Points []PointM
}

func (multiPointM MultiPointM) Bounds() *Bounds {
	return NewBounds().ExtendPointMs(multiPointM.Points)
}

type MultiPointZM struct {
	Points []PointZM
}

func (multiPointZM MultiPointZM) Bounds() *Bounds {
	return NewBounds().ExtendPointZMs(multiPointZM.Points)
}
