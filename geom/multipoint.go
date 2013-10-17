package geom

type MultiPoint struct {
	Points []Point
}

func (multiPoint MultiPoint) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, point := range multiPoint.Points {
		b = point.Bounds(b)
	}
	return b
}

type MultiPointZ struct {
	Points []PointZ
}

func (multiPointZ MultiPointZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, pointZ := range multiPointZ.Points {
		b = pointZ.Bounds(b)
	}
	return b
}

type MultiPointM struct {
	Points []PointM
}

func (multiPointM MultiPointM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, pointM := range multiPointM.Points {
		b = pointM.Bounds(b)
	}
	return b
}

type MultiPointZM struct {
	Points []PointZM
}

func (multiPointZM MultiPointZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, pointZM := range multiPointZM.Points {
		b = pointZM.Bounds(b)
	}
	return b
}
