package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func lineStringReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	points, err := readPoints(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineString(points), nil
}

func writeLineString(w io.Writer, byteOrder binary.ByteOrder, lineString geom.LineString) error {
	return writePoints(w, byteOrder, lineString)
}
