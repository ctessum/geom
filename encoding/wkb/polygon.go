package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func polygonReader(r io.Reader, byteOrder binary.ByteOrder) (geom.Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make([]geom.Path, numRings)
	for i := uint32(0); i < numRings; i++ {
		if points, err := readPoints(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = points
		}
	}
	return geom.Polygon(rings), nil
}

func writePolygon(w io.Writer, byteOrder binary.ByteOrder, polygon geom.Polygon) error {
	return writePointss(w, byteOrder, polygon)
}
