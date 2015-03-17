package geom

type MultiPolygon []Polygon

func (multiPolygon MultiPolygon) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygon := range multiPolygon {
		b = polygon.Bounds(b)
	}
	return b
}

type MultiPolygonZ []PolygonZ

func (multiPolygonZ MultiPolygonZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygonZ := range multiPolygonZ {
		b = polygonZ.Bounds(b)
	}
	return b
}

type MultiPolygonM []PolygonM

func (multiPolygonM MultiPolygonM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygonM := range multiPolygonM {
		b = polygonM.Bounds(b)
	}
	return b
}

type MultiPolygonZM []PolygonZM

func (multiPolygonZM MultiPolygonZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygonZM := range multiPolygonZM {
		b = polygonZM.Bounds(b)
	}
	return b
}
