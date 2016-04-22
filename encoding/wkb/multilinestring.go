package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func multiLineStringReader(r io.Reader, byteOrder binary.ByteOrder) (geom.Geom, error) {
	var numLineStrings uint32
	if err := binary.Read(r, byteOrder, &numLineStrings); err != nil {
		return nil, err
	}
	lineStrings := make([]geom.LineString, numLineStrings)
	for i := uint32(0); i < numLineStrings; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			lineStrings[i], ok = g.(geom.LineString)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiLineString(lineStrings), nil
}

func writeMultiLineString(w io.Writer, byteOrder binary.ByteOrder, multiLineString geom.MultiLineString) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiLineString))); err != nil {
		return err
	}
	for _, lineString := range multiLineString {
		if err := Write(w, byteOrder, lineString); err != nil {
			return err
		}
	}
	return nil
}
