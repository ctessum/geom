package geom

type MultiPolygon struct {
	Polygons []Polygon
}

func (multiPolygon MultiPolygon) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygon := range multiPolygon.Polygons {
		b = polygon.Bounds(b)
	}
	return b
}

type MultiPolygonZ struct {
	Polygons []PolygonZ
}

func (multiPolygonZ MultiPolygonZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygonZ := range multiPolygonZ.Polygons {
		b = polygonZ.Bounds(b)
	}
	return b
}

type MultiPolygonM struct {
	Polygons []PolygonM
}

func (multiPolygonM MultiPolygonM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygonM := range multiPolygonM.Polygons {
		b = polygonM.Bounds(b)
	}
	return b
}

type MultiPolygonZM struct {
	Polygons []PolygonZM
}

func (multiPolygonZM MultiPolygonZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, polygonZM := range multiPolygonZM.Polygons {
		b = polygonZM.Bounds(b)
	}
	return b
}
