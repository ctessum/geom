package geom

// MultiPoint is a holder for multiple related points.
type MultiPoint []Point

// Bounds gives the rectangular extents of the MultiPoint.
func (mp MultiPoint) Bounds() *Bounds {
	b := NewBounds()
	for _, p := range mp {
		b.extendPoint(p)
	}
	return b
}

// Within calculates whether all of the points in mp are within poly or touching
// its edge.
func (mp MultiPoint) Within(poly Polygonal) WithinStatus {
	for _, p := range mp {
		if pointInPolygonal(p, poly) == Outside {
			return Outside
		}
	}
	return Inside
}

// Len returns the number of points in the receiver.
func (mp MultiPoint) Len() int { return len(mp) }

// Points returns an iterator for the points in the receiver.
func (mp MultiPoint) Points() func() Point {
	var i int
	return func() Point {
		i++
		return mp[i-1]
	}
}
