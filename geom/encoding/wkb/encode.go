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

func writeMultiPoint(w io.Writer, byteOrder binary.ByteOrder, multiPoint geom.MultiPoint) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPoint.Points))); err != nil {
		return err
	}
	for _, point := range multiPoint.Points {
		if err := Write(w, byteOrder, point); err != nil {
			return err
		}
	}
	return nil
}

func writeMultiPointZ(w io.Writer, byteOrder binary.ByteOrder, multiPointZ geom.MultiPointZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPointZ.Points))); err != nil {
		return err
	}
	for _, pointZ := range multiPointZ.Points {
		if err := Write(w, byteOrder, pointZ); err != nil {
			return err
		}
	}
	return nil
}

func writeMultiPointM(w io.Writer, byteOrder binary.ByteOrder, multiPointM geom.MultiPointM) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPointM.Points))); err != nil {
		return err
	}
	for _, pointM := range multiPointM.Points {
		if err := Write(w, byteOrder, pointM); err != nil {
			return err
		}
	}
	return nil
}

func writeMultiPointZM(w io.Writer, byteOrder binary.ByteOrder, multiPointZM geom.MultiPointZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPointZM.Points))); err != nil {
		return err
	}
	for _, pointZM := range multiPointZM.Points {
		if err := Write(w, byteOrder, pointZM); err != nil {
			return err
		}
	}
	return nil
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
	case geom.MultiPoint:
		wkbGeometryType = wkbMultiPoint
	case geom.MultiPointZ:
		wkbGeometryType = wkbMultiPointZ
	case geom.MultiPointM:
		wkbGeometryType = wkbMultiPointM
	case geom.MultiPointZM:
		wkbGeometryType = wkbMultiPointZM
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
	case geom.MultiPoint:
		return writeMultiPoint(w, byteOrder, g.(geom.MultiPoint))
	case geom.MultiPointZ:
		return writeMultiPointZ(w, byteOrder, g.(geom.MultiPointZ))
	case geom.MultiPointM:
		return writeMultiPointM(w, byteOrder, g.(geom.MultiPointM))
	case geom.MultiPointZM:
		return writeMultiPointZM(w, byteOrder, g.(geom.MultiPointZM))
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
