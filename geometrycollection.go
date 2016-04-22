package geom

// GeometryCollection is a holder for multiple related geometry objects of
// arbitrary type.
type GeometryCollection []Geom

// Bounds gives the rectangular extents of the GeometryCollection.
func (gc GeometryCollection) Bounds() *Bounds {
	b := NewBounds()
	for _, geom := range gc {
		b.Extend(geom.Bounds())
	}
	return b
}

// Within calculates whether gc is completely within p. Points that touch
// the edge of p are considered within.
func (gc GeometryCollection) Within(p Polygonal) bool {
	for _, g := range gc {
		if !g.Within(p) {
			return false
		}
	}
	return true
}
