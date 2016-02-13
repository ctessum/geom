package geom

// MultiLineString is a holder for multiple related LineStrings.
type MultiLineString []LineString

// Bounds gives the rectangular extents of the MultiLineString.
func (multiLineString MultiLineString) Bounds() *Bounds {
	b := NewBounds()
	for _, lineString := range multiLineString {
		b.Extend(lineString.Bounds())
	}
	return b
}
