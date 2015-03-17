package geom

type MultiPoint []Point

func (multiPoint MultiPoint) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, point := range multiPoint {
		b = point.Bounds(b)
	}
	return b
}

type MultiPointZ []PointZ

func (multiPointZ MultiPointZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, pointZ := range multiPointZ {
		b = pointZ.Bounds(b)
	}
	return b
}

type MultiPointM []PointM

func (multiPointM MultiPointM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, pointM := range multiPointM {
		b = pointM.Bounds(b)
	}
	return b
}

type MultiPointZM []PointZM

func (multiPointZM MultiPointZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, pointZM := range multiPointZM {
		b = pointZM.Bounds(b)
	}
	return b
}
