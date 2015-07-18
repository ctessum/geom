package geom

type Polygon [][]Point

func (polygon Polygon) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointss(polygon)
}
