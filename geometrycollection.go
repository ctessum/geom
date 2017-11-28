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

// Len returns the number of points in the receiver.
func (gc GeometryCollection) Len() int {
	var i int
	for _, g := range gc {
		i += g.Len()
	}
	return i
}

// Points returns an iterator for the points in the receiver.
func (gc GeometryCollection) Points() func() Point {
	var i, j int
	p := gc[0].Points()
	return func() Point {
		if i == gc[j].Len() {
			j++
			i = 0
			p = gc[j].Points()
		}
		i++
		return p()
	}
}
