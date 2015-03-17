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

type MultiLineStringZ []LineStringZ

func (multiLineStringZ MultiLineStringZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineStringZ := range multiLineStringZ {
		b = lineStringZ.Bounds(b)
	}
	return b
}

type MultiLineStringM []LineStringM

func (multiLineStringM MultiLineStringM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineStringM := range multiLineStringM {
		b = lineStringM.Bounds(b)
	}
	return b
}

type MultiLineStringZM []LineStringZM

func (multiLineStringZM MultiLineStringZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineStringZM := range multiLineStringZM {
		b = lineStringZM.Bounds(b)
	}
	return b
}
