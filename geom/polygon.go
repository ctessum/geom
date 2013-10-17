package geom

type Polygon struct {
	Rings [][]Point
}

func (polygon Polygon) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointss(polygon.Rings)
}

type PolygonZ struct {
	Rings [][]PointZ
}

func (polygonZ PolygonZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZss(polygonZ.Rings)
}

type PolygonM struct {
	Rings [][]PointM
}

func (polygonM PolygonM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointMss(polygonM.Rings)
}

type PolygonZM struct {
	Rings [][]PointZM
}

func (polygonZM PolygonZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZMss(polygonZM.Rings)
}
