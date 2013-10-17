package geom

type MultiPolygon struct {
	Polygons []Polygon
}

func (multiPolygon MultiPolygon) Bounds() *Bounds {
	bounds := NewBounds()
	for _, polygon := range multiPolygon.Polygons {
		bounds.ExtendPointss(polygon.Rings)
	}
	return bounds
}

type MultiPolygonZ struct {
	Polygons []PolygonZ
}

func (multiPolygonZ MultiPolygonZ) Bounds() *Bounds {
	bounds := NewBounds()
	for _, polygonZ := range multiPolygonZ.Polygons {
		bounds.ExtendPointZss(polygonZ.Rings)
	}
	return bounds
}

type MultiPolygonM struct {
	Polygons []PolygonM
}

func (multiPolygonM MultiPolygonM) Bounds() *Bounds {
	bounds := NewBounds()
	for _, polygonM := range multiPolygonM.Polygons {
		bounds.ExtendPointMss(polygonM.Rings)
	}
	return bounds
}

type MultiPolygonZM struct {
	Polygons []PolygonZM
}

func (multiPolygonZM MultiPolygonZM) Bounds() *Bounds {
	bounds := NewBounds()
	for _, polygonZM := range multiPolygonZM.Polygons {
		bounds.ExtendPointZMss(polygonZM.Rings)
	}
	return bounds
}
