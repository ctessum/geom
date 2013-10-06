package geom

type T interface {
	Bounds() *Bounds
}

type LineString struct {
	Points LinearRing
}

func (lineString LineString) Bounds() *Bounds {
	return lineString.Points.Bounds()
}

type LineStringZ struct {
	Points LinearRingZ
}

func (lineStringZ LineStringZ) Bounds() *Bounds {
	return lineStringZ.Points.Bounds()
}

type LineStringM struct {
	Points LinearRingM
}

func (lineStringM LineStringM) Bounds() *Bounds {
	return lineStringM.Points.Bounds()
}

type LineStringZM struct {
	Points LinearRingZM
}

func (lineStringZM LineStringZM) Bounds() *Bounds {
	return lineStringZM.Points.Bounds()
}

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
