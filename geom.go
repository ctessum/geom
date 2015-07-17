/*
Package geom holds geometry objects that can be encoded, decoded,
and operated on by the other packages in this repository.
*/
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
