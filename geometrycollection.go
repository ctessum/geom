package geom

type GeometryCollection []T

func (geometryCollection GeometryCollection) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geom := range geometryCollection {
		b = geom.Bounds(b)
	}
	return b
}
