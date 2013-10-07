package geom

type LineString struct {
	Points []Point
}

func (lineString LineString) Bounds() *Bounds {
	return NewBounds().ExtendPoints(lineString.Points)
}

type LineStringZ struct {
	Points []PointZ
}

func (lineStringZ LineStringZ) Bounds() *Bounds {
	return NewBounds().ExtendPointZs(lineStringZ.Points)
}

type LineStringM struct {
	Points []PointM
}

func (lineStringM LineStringM) Bounds() *Bounds {
	return NewBounds().ExtendPointMs(lineStringM.Points)
}

type LineStringZM struct {
	Points []PointZM
}

func (lineStringZM LineStringZM) Bounds() *Bounds {
	return NewBounds().ExtendPointZMs(lineStringZM.Points)
}
