package geom

// LineString is a number of points that make up a path or line.
type LineString []Point

// Bounds gives the rectangular extents of the LineString.
func (lineString LineString) Bounds() *Bounds {
	b := NewBounds()
	b.extendPoints(lineString)
	return b
}
