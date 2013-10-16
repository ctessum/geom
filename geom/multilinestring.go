package geom

type MultiLineString struct {
	LineStrings []LineString
}

func (multiLineString MultiLineString) Bounds() *Bounds {
	bounds := NewBounds()
	for _, lineString := range multiLineString.LineStrings {
		bounds.ExtendPoints(lineString.Points)
	}
	return bounds
}

type MultiLineStringZ struct {
	LineStrings []LineStringZ
}

func (multiLineStringZ MultiLineStringZ) Bounds() *Bounds {
	bounds := NewBounds()
	for _, lineStringZ := range multiLineStringZ.LineStrings {
		bounds.ExtendPointZs(lineStringZ.Points)
	}
	return bounds
}

type MultiLineStringM struct {
	LineStrings []LineStringM
}

func (multiLineStringM MultiLineStringM) Bounds() *Bounds {
	bounds := NewBounds()
	for _, lineStringM := range multiLineStringM.LineStrings {
		bounds.ExtendPointMs(lineStringM.Points)
	}
	return bounds
}

type MultiLineStringZM struct {
	LineStrings []LineStringZM
}

func (multiLineStringZM MultiLineStringZM) Bounds() *Bounds {
	bounds := NewBounds()
	for _, lineStringZM := range multiLineStringZM.LineStrings {
		bounds.ExtendPointZMs(lineStringZM.Points)
	}
	return bounds
}
