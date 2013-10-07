package geom

type Polygon struct {
	Rings LinearRings
}

func (polygon Polygon) Bounds() *Bounds {
	return NewBounds().ExtendLinearRings(polygon.Rings)
}

type PolygonZ struct {
	Rings LinearRingZs
}

func (polygonZ PolygonZ) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZs(polygonZ.Rings)
}

type PolygonM struct {
	Rings LinearRingMs
}

func (polygonM PolygonM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingMs(polygonM.Rings)
}

type PolygonZM struct {
	Rings LinearRingZMs
}

func (polygonZM PolygonZM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZMs(polygonZM.Rings)
}
