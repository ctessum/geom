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
