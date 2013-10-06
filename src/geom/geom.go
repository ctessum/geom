package geom

type T interface {
	Bounds() *Bounds
}

type LinearRing []Point

func (linearRing LinearRing) Bounds() *Bounds {
	return NewBounds().ExtendLinearRing(linearRing)
}

type LinearRingZ []PointZ

func (linearRingZ LinearRingZ) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZ(linearRingZ)
}

type LinearRingM []PointM

func (linearRingM LinearRingM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingM(linearRingM)
}

type LinearRingZM []PointZM

func (linearRingZM LinearRingZM) Bounds() *Bounds {
	return NewBounds().ExtendLinearRingZM(linearRingZM)
}

type LinearRings []LinearRing
type LinearRingZs []LinearRingZ
type LinearRingMs []LinearRingM
type LinearRingZMs []LinearRingZM

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
