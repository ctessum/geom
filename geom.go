package geom

type T interface {
	Bounds(*Bounds) *Bounds
}

type Geom interface {
	T
}

type GeomZ interface {
	T
}

type GeomM interface {
	T
}

type GeomZM interface {
	T
}
