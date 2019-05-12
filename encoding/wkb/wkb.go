package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/ctessum/geom"
)

const (
	wkbXDR = 0
	wkbNDR = 1
)

const (
	wkbPoint              = 1
	wkbLineString         = 2
	wkbPolygon            = 3
	wkbMultiPoint         = 4
	wkbMultiLineString    = 5
	wkbMultiPolygon       = 6
	wkbGeometryCollection = 7
	wkbPolyhedralSurface  = 15
	wkbTIN                = 16
	wkbTriangle           = 17
)

var (
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

type UnexpectedGeometryError struct {
	Geom geom.Geom
}

func (e UnexpectedGeometryError) Error() string {
	return fmt.Sprintf("wkb: unexpected geometry %v", e.Geom)
}

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "wkb: unsupported type: " + e.Type.String()
}

type wkbReader func(io.Reader, binary.ByteOrder) (geom.Geom, error)

var wkbReaders map[uint32]wkbReader

func init() {
	wkbReaders = make(map[uint32]wkbReader)
	wkbReaders[wkbPoint] = pointReader
	wkbReaders[wkbLineString] = lineStringReader
	wkbReaders[wkbPolygon] = polygonReader
	wkbReaders[wkbMultiPoint] = multiPointReader
	wkbReaders[wkbMultiLineString] = multiLineStringReader
	wkbReaders[wkbMultiPolygon] = multiPolygonReader
	wkbReaders[wkbGeometryCollection] = geometryCollectionReader
}

func Read(r io.Reader) (geom.Geom, error) {

	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return nil, err
	}
	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case wkbXDR:
		byteOrder = binary.BigEndian
	case wkbNDR:
		byteOrder = binary.LittleEndian
	default:
		return nil, fmt.Errorf("invalid byte order %v", wkbByteOrder)
	}

	var wkbGeometryType uint32
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return nil, err
	}

	if reader, ok := wkbReaders[wkbGeometryType]; ok {
		return reader(r, byteOrder)
	} else {
		return nil, fmt.Errorf("unsupported geometry type %v", wkbGeometryType)
	}

}

func Decode(buf []byte) (geom.Geom, error) {
	return Read(bytes.NewBuffer(buf))
}

func writeMany(w io.Writer, byteOrder binary.ByteOrder, data ...interface{}) error {
	for _, datum := range data {
		if err := binary.Write(w, byteOrder, datum); err != nil {
			return err
		}
	}
	return nil
}

func Write(w io.Writer, byteOrder binary.ByteOrder, g geom.Geom) error {
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
	case geom.LineString:
		wkbGeometryType = wkbLineString
	case geom.Polygon:
		wkbGeometryType = wkbPolygon
	case geom.MultiPoint:
		wkbGeometryType = wkbMultiPoint
	case geom.MultiLineString:
		wkbGeometryType = wkbMultiLineString
	case geom.MultiPolygon:
		wkbGeometryType = wkbMultiPolygon
	case geom.GeometryCollection:
		wkbGeometryType = wkbGeometryCollection
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
	if err := binary.Write(w, byteOrder, wkbGeometryType); err != nil {
		return err
	}
	switch g.(type) {
	case geom.Point:
		return writePoint(w, byteOrder, g.(geom.Point))
	case geom.LineString:
		return writeLineString(w, byteOrder, g.(geom.LineString))
	case geom.Polygon:
		return writePolygon(w, byteOrder, g.(geom.Polygon))
	case geom.MultiPoint:
		return writeMultiPoint(w, byteOrder, g.(geom.MultiPoint))
	case geom.MultiLineString:
		return writeMultiLineString(w, byteOrder, g.(geom.MultiLineString))
	case geom.MultiPolygon:
		return writeMultiPolygon(w, byteOrder, g.(geom.MultiPolygon))
	case geom.GeometryCollection:
		return writeGeometryCollection(w, byteOrder, g.(geom.GeometryCollection))
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}

func Encode(g geom.Geom, byteOrder binary.ByteOrder) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
