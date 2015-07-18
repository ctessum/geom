package geom

type LineString []Point

func (lineString LineString) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPoints(lineString)
}
