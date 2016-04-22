package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func geometryCollectionReader(r io.Reader, byteOrder binary.ByteOrder) (geom.Geom, error) {
	var numGeometries uint32
	if err := binary.Read(r, byteOrder, &numGeometries); err != nil {
		return nil, err
	}
	geoms := make([]geom.Geom, numGeometries)
	for i := uint32(0); i < numGeometries; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			geoms[i], ok = g.(geom.Geom)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.GeometryCollection(geoms), nil
}

func writeGeometryCollection(w io.Writer, byteOrder binary.ByteOrder, geometryCollection geom.GeometryCollection) error {
	if err := binary.Write(w, byteOrder, uint32(len(geometryCollection))); err != nil {
		return err
	}
	for _, geom := range geometryCollection {
		if err := Write(w, byteOrder, geom); err != nil {
			return err
		}
	}
	return nil
}
