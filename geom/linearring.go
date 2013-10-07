package geom

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
