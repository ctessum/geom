package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func multiPolygonReader(r io.Reader, byteOrder binary.ByteOrder) (geom.Geom, error) {
	var numPolygons uint32
	if err := binary.Read(r, byteOrder, &numPolygons); err != nil {
		return nil, err
	}
	polygons := make([]geom.Polygon, numPolygons)
	for i := uint32(0); i < numPolygons; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			polygons[i], ok = g.(geom.Polygon)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPolygon(polygons), nil
}

func writeMultiPolygon(w io.Writer, byteOrder binary.ByteOrder, multiPolygon geom.MultiPolygon) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPolygon))); err != nil {
		return err
	}
	for _, polygon := range multiPolygon {
		if err := Write(w, byteOrder, polygon); err != nil {
			return err
		}
	}
	return nil
}
