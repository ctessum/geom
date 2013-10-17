package geom

type GeometryCollection struct {
	Geoms []Geom
}

func (geometryCollection GeometryCollection) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geom := range geometryCollection.Geoms {
		b = geom.Bounds(b)
	}
	return b
}

type GeometryCollectionZ struct {
	Geoms []GeomZ
}

func (geometryCollectionZ GeometryCollectionZ) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geomZ := range geometryCollectionZ.Geoms {
		b = geomZ.Bounds(b)
	}
	return b
}

type GeometryCollectionM struct {
	Geoms []GeomM
}

func (geometryCollectionM GeometryCollectionM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geomM := range geometryCollectionM.Geoms {
		b = geomM.Bounds(b)
	}
	return b
}

type GeometryCollectionZM struct {
	Geoms []GeomZM
}

func (geometryCollectionZM GeometryCollectionZM) Bounds(b *Bounds) *Bounds {
	if b == nil {
		b = NewBounds()
	}
	for _, geomZM := range geometryCollectionZM.Geoms {
		b = geomZM.Bounds(b)
	}
	return b
}
