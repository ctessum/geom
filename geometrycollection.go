package geom

type GeometryCollection []Geom

func (geometryCollection GeometryCollection) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geom := range geometryCollection {
		b = geom.Bounds(b)
	}
	return b
}

type GeometryCollectionZ []GeomZ

func (geometryCollectionZ GeometryCollectionZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geomZ := range geometryCollectionZ {
		b = geomZ.Bounds(b)
	}
	return b
}

type GeometryCollectionM []GeomM

func (geometryCollectionM GeometryCollectionM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geomM := range geometryCollectionM {
		b = geomM.Bounds(b)
	}
	return b
}

type GeometryCollectionZM []GeomZM

func (geometryCollectionZM GeometryCollectionZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geomZM := range geometryCollectionZM {
		b = geomZM.Bounds(b)
	}
	return b
}
