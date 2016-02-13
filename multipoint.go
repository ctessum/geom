package geom

// MultiPoint is a holder for multiple related points.
type MultiPoint []Point

// Bounds gives the rectangular extents of the MultiPoint.
func (multiPoint MultiPoint) Bounds() *Bounds {
	b := NewBounds()
	for _, point := range multiPoint {
		b.extendPoint(point)
	}
	return b
}
