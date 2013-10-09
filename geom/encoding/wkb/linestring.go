package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)

func lineStringReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	points, err := readPoints(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineString{points}, nil
}

func writeLineString(w io.Writer, byteOrder binary.ByteOrder, lineString geom.LineString) error {
	return writePoints(w, byteOrder, lineString.Points)
}

func lineStringZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZs, err := readPointZs(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringZ{pointZs}, nil
}

func writeLineStringZ(w io.Writer, byteOrder binary.ByteOrder, lineStringZ geom.LineStringZ) error {
	return writePointZs(w, byteOrder, lineStringZ.Points)
}

func lineStringMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointMs, err := readPointMs(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringM{pointMs}, nil
}

func writeLineStringM(w io.Writer, byteOrder binary.ByteOrder, lineStringM geom.LineStringM) error {
	return writePointMs(w, byteOrder, lineStringM.Points)
}

func lineStringZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZMs, err := readPointZMs(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringZM{pointZMs}, nil
}

func writeLineStringZM(w io.Writer, byteOrder binary.ByteOrder, lineStringZM geom.LineStringZM) error {
	return writePointZMs(w, byteOrder, lineStringZM.Points)
}
