package geom

type LineString struct {
	Points []Point
}

func (lineString LineString) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPoints(lineString.Points)
}

type LineStringZ struct {
	Points []PointZ
}

func (lineStringZ LineStringZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZs(lineStringZ.Points)
}

type LineStringM struct {
	Points []PointM
}

func (lineStringM LineStringM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointMs(lineStringM.Points)
}

type LineStringZM struct {
	Points []PointZM
}

func (lineStringZM LineStringZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZMs(lineStringZM.Points)
}
