package geom

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
