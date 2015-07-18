package geom

type MultiLineString []LineString

func (multiLineString MultiLineString) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineString := range multiLineString {
		b = lineString.Bounds(b)
	}
	return b
}
