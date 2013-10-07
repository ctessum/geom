package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"io"
	"reflect"
)

func writeMany(w io.Writer, byteOrder binary.ByteOrder, data ...interface{}) error {
	for _, datum := range data {
		if err := binary.Write(w, byteOrder, datum); err != nil {
			return err
		}
	}
	return nil
}

func writePoint(w io.Writer, byteOrder binary.ByteOrder, point geom.Point) error {
	return binary.Write(w, byteOrder, &point)
}

func writePointZ(w io.Writer, byteOrder binary.ByteOrder, pointZ geom.PointZ) error {
	return binary.Write(w, byteOrder, &pointZ)
}

func writePointM(w io.Writer, byteOrder binary.ByteOrder, pointM geom.PointM) error {
	return binary.Write(w, byteOrder, &pointM)
}

func writePointZM(w io.Writer, byteOrder binary.ByteOrder, pointZM geom.PointZM) error {
	return binary.Write(w, byteOrder, &pointZM)
}

func writeLinearRing(w io.Writer, byteOrder binary.ByteOrder, linearRing geom.LinearRing) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRing))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRing)
}

func writeLinearRingZ(w io.Writer, byteOrder binary.ByteOrder, linearRingZ geom.LinearRingZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingZ))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRingZ)
}

func writeLinearRingM(w io.Writer, byteOrder binary.ByteOrder, linearRingM geom.LinearRingM) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingM))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRingM)
}

func writeLinearRingZM(w io.Writer, byteOrder binary.ByteOrder, linearRingZM geom.LinearRingZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingZM))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRingZM)
}

func writeLinearRings(w io.Writer, byteOrder binary.ByteOrder, linearRings geom.LinearRings) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRings))); err != nil {
		return err
	}
	for _, linearRing := range linearRings {
		if err := writeLinearRing(w, byteOrder, linearRing); err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRingZs(w io.Writer, byteOrder binary.ByteOrder, linearRingZs geom.LinearRingZs) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingZs))); err != nil {
		return err
	}
	for _, linearRingZ := range linearRingZs {
		if err := writeLinearRingZ(w, byteOrder, linearRingZ); err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRingMs(w io.Writer, byteOrder binary.ByteOrder, linearRingMs geom.LinearRingMs) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingMs))); err != nil {
		return err
	}
	for _, linearRingM := range linearRingMs {
		if err := writeLinearRingM(w, byteOrder, linearRingM); err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRingZMs(w io.Writer, byteOrder binary.ByteOrder, linearRingZMs geom.LinearRingZMs) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingZMs))); err != nil {
		return err
	}
	for _, linearRingZM := range linearRingZMs {
		if err := writeLinearRingZM(w, byteOrder, linearRingZM); err != nil {
			return err
		}
	}
	return nil
}

func writeLineString(w io.Writer, byteOrder binary.ByteOrder, lineString geom.LineString) error {
	return writeLinearRing(w, byteOrder, lineString.Points)
}

func writeLineStringZ(w io.Writer, byteOrder binary.ByteOrder, lineStringZ geom.LineStringZ) error {
	return writeLinearRingZ(w, byteOrder, lineStringZ.Points)
}

func writeLineStringM(w io.Writer, byteOrder binary.ByteOrder, lineStringM geom.LineStringM) error {
	return writeLinearRingM(w, byteOrder, lineStringM.Points)
}

func writeLineStringZM(w io.Writer, byteOrder binary.ByteOrder, lineStringZM geom.LineStringZM) error {
	return writeLinearRingZM(w, byteOrder, lineStringZM.Points)
}

func writePolygon(w io.Writer, byteOrder binary.ByteOrder, polygon geom.Polygon) error {
	return writeLinearRings(w, byteOrder, polygon.Rings)
}

func writePolygonZ(w io.Writer, byteOrder binary.ByteOrder, polygonZ geom.PolygonZ) error {
	return writeLinearRingZs(w, byteOrder, polygonZ.Rings)
}

func writePolygonM(w io.Writer, byteOrder binary.ByteOrder, polygonM geom.PolygonM) error {
	return writeLinearRingMs(w, byteOrder, polygonM.Rings)
}

func writePolygonZM(w io.Writer, byteOrder binary.ByteOrder, polygonZM geom.PolygonZM) error {
	return writeLinearRingZMs(w, byteOrder, polygonZM.Rings)
}

func Write(w io.Writer, byteOrder binary.ByteOrder, g geom.T) error {
	var wkbByteOrder uint8
	switch byteOrder {
	case XDR:
		wkbByteOrder = wkbXDR
	case NDR:
		wkbByteOrder = wkbNDR
	default:
		return fmt.Errorf("unsupported byte order %v", byteOrder)
	}
	if err := binary.Write(w, byteOrder, wkbByteOrder); err != nil {
		return err
	}
	var wkbGeometryType uint32
	switch g.(type) {
	case geom.Point:
		wkbGeometryType = wkbPoint
	case geom.PointZ:
		wkbGeometryType = wkbPointZ
	case geom.PointM:
		wkbGeometryType = wkbPointM
	case geom.PointZM:
		wkbGeometryType = wkbPointZM
	case geom.LineString:
		wkbGeometryType = wkbLineString
	case geom.LineStringZ:
		wkbGeometryType = wkbLineStringZ
	case geom.LineStringM:
		wkbGeometryType = wkbLineStringM
	case geom.LineStringZM:
		wkbGeometryType = wkbLineStringZM
	case geom.Polygon:
		wkbGeometryType = wkbPolygon
	case geom.PolygonZ:
		wkbGeometryType = wkbPolygonZ
	case geom.PolygonM:
		wkbGeometryType = wkbPolygonM
	case geom.PolygonZM:
		wkbGeometryType = wkbPolygonZM
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
	if err := binary.Write(w, byteOrder, wkbGeometryType); err != nil {
		return err
	}
	switch g.(type) {
	case geom.Point:
		return writePoint(w, byteOrder, g.(geom.Point))
	case geom.PointZ:
		return writePointZ(w, byteOrder, g.(geom.PointZ))
	case geom.PointM:
		return writePointM(w, byteOrder, g.(geom.PointM))
	case geom.PointZM:
		return writePointZM(w, byteOrder, g.(geom.PointZM))
	case geom.LineString:
		return writeLineString(w, byteOrder, g.(geom.LineString))
	case geom.LineStringZ:
		return writeLineStringZ(w, byteOrder, g.(geom.LineStringZ))
	case geom.LineStringM:
		return writeLineStringM(w, byteOrder, g.(geom.LineStringM))
	case geom.LineStringZM:
		return writeLineStringZM(w, byteOrder, g.(geom.LineStringZM))
	case geom.Polygon:
		return writePolygon(w, byteOrder, g.(geom.Polygon))
	case geom.PolygonZ:
		return writePolygonZ(w, byteOrder, g.(geom.PolygonZ))
	case geom.PolygonM:
		return writePolygonM(w, byteOrder, g.(geom.PolygonM))
	case geom.PolygonZM:
		return writePolygonZM(w, byteOrder, g.(geom.PolygonZM))
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}

func Marshal(g geom.T, byteOrder binary.ByteOrder) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
