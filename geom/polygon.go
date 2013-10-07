package geom

type Polygon struct {
	Rings [][]Point
}

func (polygon Polygon) Bounds() *Bounds {
	return NewBounds().ExtendPointss(polygon.Rings)
}

type PolygonZ struct {
	Rings [][]PointZ
}

func (polygonZ PolygonZ) Bounds() *Bounds {
	return NewBounds().ExtendPointZss(polygonZ.Rings)
}

type PolygonM struct {
	Rings [][]PointM
}

func (polygonM PolygonM) Bounds() *Bounds {
	return NewBounds().ExtendPointMss(polygonM.Rings)
}

type PolygonZM struct {
	Rings [][]PointZM
}

func (polygonZM PolygonZM) Bounds() *Bounds {
	return NewBounds().ExtendPointZMss(polygonZM.Rings)
}
