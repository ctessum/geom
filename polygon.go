package geom

type Polygon [][]Point

func (polygon Polygon) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointss(polygon)
}

type PolygonZ [][]PointZ

func (polygonZ PolygonZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZss(polygonZ)
}

type PolygonM [][]PointM

func (polygonM PolygonM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointMss(polygonM)
}

type PolygonZM [][]PointZM

func (polygonZM PolygonZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	return b.ExtendPointZMss(polygonZM)
}
