package geom

type LineString []Point

func (lineString LineString) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPoints(lineString)
}

type LineStringZ []PointZ

func (lineStringZ LineStringZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZs(lineStringZ)
}

type LineStringM []PointM

func (lineStringM LineStringM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointMs(lineStringM)
}

type LineStringZM []PointZM

func (lineStringZM LineStringZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZMs(lineStringZM)
}
