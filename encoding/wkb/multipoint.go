package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func multiPointReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make([]geom.Point, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			points[i], ok = g.(geom.Point)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPoint(points), nil
}

func writeMultiPoint(w io.Writer, byteOrder binary.ByteOrder, multiPoint geom.MultiPoint) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPoint))); err != nil {
		return err
	}
	for _, point := range multiPoint {
		if err := Write(w, byteOrder, point); err != nil {
			return err
		}
	}
	return nil
}

func multiPointZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZs := make([]geom.PointZ, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			pointZs[i], ok = g.(geom.PointZ)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPointZ(pointZs), nil
}

func writeMultiPointZ(w io.Writer, byteOrder binary.ByteOrder, multiPointZ geom.MultiPointZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPointZ))); err != nil {
		return err
	}
	for _, pointZ := range multiPointZ {
		if err := Write(w, byteOrder, pointZ); err != nil {
			return err
		}
	}
	return nil
}

func multiPointMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointMs := make([]geom.PointM, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			pointMs[i], ok = g.(geom.PointM)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPointM(pointMs), nil
}

func writeMultiPointM(w io.Writer, byteOrder binary.ByteOrder, multiPointM geom.MultiPointM) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPointM))); err != nil {
		return err
	}
	for _, pointM := range multiPointM {
		if err := Write(w, byteOrder, pointM); err != nil {
			return err
		}
	}
	return nil
}

func multiPointZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZMs := make([]geom.PointZM, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			pointZMs[i], ok = g.(geom.PointZM)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPointZM(pointZMs), nil
}

func writeMultiPointZM(w io.Writer, byteOrder binary.ByteOrder, multiPointZM geom.MultiPointZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPointZM))); err != nil {
		return err
	}
	for _, pointZM := range multiPointZM {
		if err := Write(w, byteOrder, pointZM); err != nil {
			return err
		}
	}
	return nil
}
