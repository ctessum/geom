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
