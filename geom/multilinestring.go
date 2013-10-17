package geom

type MultiLineString struct {
	LineStrings []LineString
}

func (multiLineString MultiLineString) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineString := range multiLineString.LineStrings {
		b = lineString.Bounds(b)
	}
	return b
}

type MultiLineStringZ struct {
	LineStrings []LineStringZ
}

func (multiLineStringZ MultiLineStringZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineStringZ := range multiLineStringZ.LineStrings {
		b = lineStringZ.Bounds(b)
	}
	return b
}

type MultiLineStringM struct {
	LineStrings []LineStringM
}

func (multiLineStringM MultiLineStringM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineStringM := range multiLineStringM.LineStrings {
		b = lineStringM.Bounds(b)
	}
	return b
}

type MultiLineStringZM struct {
	LineStrings []LineStringZM
}

func (multiLineStringZM MultiLineStringZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, lineStringZM := range multiLineStringZM.LineStrings {
		b = lineStringZM.Bounds(b)
	}
	return b
}
