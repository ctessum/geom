package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ctessum/gogeom/geom"
	"io"
	"reflect"
)

const (
	wkbXDR = 0
	wkbNDR = 1
)

const (
	wkbPoint                = 1
	wkbPointZ               = 1001
	wkbPointM               = 2001
	wkbPointZM              = 3001
	wkbLineString           = 2
	wkbLineStringZ          = 1002
	wkbLineStringM          = 2002
	wkbLineStringZM         = 3002
	wkbPolygon              = 3
	wkbPolygonZ             = 1003
	wkbPolygonM             = 2003
	wkbPolygonZM            = 3003
	wkbMultiPoint           = 4
	wkbMultiPointZ          = 1004
	wkbMultiPointM          = 2004
	wkbMultiPointZM         = 3004
	wkbMultiLineString      = 5
	wkbMultiLineStringZ     = 1005
	wkbMultiLineStringM     = 2005
	wkbMultiLineStringZM    = 3005
	wkbMultiPolygon         = 6
	wkbMultiPolygonZ        = 1006
	wkbMultiPolygonM        = 2006
	wkbMultiPolygonZM       = 3006
	wkbGeometryCollection   = 7
	wkbGeometryCollectionZ  = 1007
	wkbGeometryCollectionM  = 2007
	wkbGeometryCollectionZM = 3007
	wkbPolyhedralSurface    = 15
	wkbPolyhedralSurfaceZ   = 1015
	wkbPolyhedralSurfaceM   = 2015
	wkbPolyhedralSurfaceZM  = 3015
	wkbTIN                  = 16
	wkbTINZ                 = 1016
	wkbTINM                 = 2016
	wkbTINZM                = 3016
	wkbTriangle             = 17
	wkbTriangleZ            = 1017
	wkbTriangleM            = 2017
	wkbTriangleZM           = 3017
)

var (
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

type UnexpectedGeometryError struct {
	Geom geom.T
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

type wkbReader func(io.Reader, binary.ByteOrder) (geom.T, error)

var wkbReaders map[uint32]wkbReader

func init() {
	wkbReaders = make(map[uint32]wkbReader)
	wkbReaders[wkbPoint] = pointReader
	wkbReaders[wkbPointZ] = pointZReader
	wkbReaders[wkbPointM] = pointMReader
	wkbReaders[wkbPointZM] = pointZMReader
	wkbReaders[wkbLineString] = lineStringReader
	wkbReaders[wkbLineStringZ] = lineStringZReader
	wkbReaders[wkbLineStringM] = lineStringMReader
	wkbReaders[wkbLineStringZM] = lineStringZMReader
	wkbReaders[wkbPolygon] = polygonReader
	wkbReaders[wkbPolygonZ] = polygonZReader
	wkbReaders[wkbPolygonM] = polygonMReader
	wkbReaders[wkbPolygonZM] = polygonZMReader
	wkbReaders[wkbMultiPoint] = multiPointReader
	wkbReaders[wkbMultiPointZ] = multiPointZReader
	wkbReaders[wkbMultiPointM] = multiPointMReader
	wkbReaders[wkbMultiPointZM] = multiPointZMReader
	wkbReaders[wkbMultiLineString] = multiLineStringReader
	wkbReaders[wkbMultiLineStringZ] = multiLineStringZReader
	wkbReaders[wkbMultiLineStringM] = multiLineStringMReader
	wkbReaders[wkbMultiLineStringZM] = multiLineStringZMReader
	wkbReaders[wkbMultiPolygon] = multiPolygonReader
	wkbReaders[wkbMultiPolygonZ] = multiPolygonZReader
	wkbReaders[wkbMultiPolygonM] = multiPolygonMReader
	wkbReaders[wkbMultiPolygonZM] = multiPolygonZMReader
	wkbReaders[wkbGeometryCollection] = geometryCollectionReader
	wkbReaders[wkbGeometryCollectionZ] = geometryCollectionZReader
	wkbReaders[wkbGeometryCollectionM] = geometryCollectionMReader
	wkbReaders[wkbGeometryCollectionZM] = geometryCollectionZMReader
}

func Read(r io.Reader) (geom.T, error) {

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
		return nil, fmt.Errorf("invalid byte order %u", wkbByteOrder)
	}

	var wkbGeometryType uint32
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return nil, err
	}

	if reader, ok := wkbReaders[wkbGeometryType]; ok {
		return reader(r, byteOrder)
	} else {
		return nil, fmt.Errorf("unsupported geometry type %u", wkbGeometryType)
	}

}

func Decode(buf []byte) (geom.T, error) {
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
	case geom.MultiLineString:
		wkbGeometryType = wkbMultiLineString
	case geom.MultiLineStringZ:
		wkbGeometryType = wkbMultiLineStringZ
	case geom.MultiLineStringM:
		wkbGeometryType = wkbMultiLineStringM
	case geom.MultiLineStringZM:
		wkbGeometryType = wkbMultiLineStringZM
	case geom.MultiPolygon:
		wkbGeometryType = wkbMultiPolygon
	case geom.MultiPolygonZ:
		wkbGeometryType = wkbMultiPolygonZ
	case geom.MultiPolygonM:
		wkbGeometryType = wkbMultiPolygonM
	case geom.MultiPolygonZM:
		wkbGeometryType = wkbMultiPolygonZM
	case geom.GeometryCollection:
		wkbGeometryType = wkbGeometryCollection
	case geom.GeometryCollectionZ:
		wkbGeometryType = wkbGeometryCollectionZ
	case geom.GeometryCollectionM:
		wkbGeometryType = wkbGeometryCollectionM
	case geom.GeometryCollectionZM:
		wkbGeometryType = wkbGeometryCollectionZM
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
	case geom.MultiLineString:
		return writeMultiLineString(w, byteOrder, g.(geom.MultiLineString))
	case geom.MultiLineStringZ:
		return writeMultiLineStringZ(w, byteOrder, g.(geom.MultiLineStringZ))
	case geom.MultiLineStringM:
		return writeMultiLineStringM(w, byteOrder, g.(geom.MultiLineStringM))
	case geom.MultiLineStringZM:
		return writeMultiLineStringZM(w, byteOrder, g.(geom.MultiLineStringZM))
	case geom.MultiPolygon:
		return writeMultiPolygon(w, byteOrder, g.(geom.MultiPolygon))
	case geom.MultiPolygonZ:
		return writeMultiPolygonZ(w, byteOrder, g.(geom.MultiPolygonZ))
	case geom.MultiPolygonM:
		return writeMultiPolygonM(w, byteOrder, g.(geom.MultiPolygonM))
	case geom.MultiPolygonZM:
		return writeMultiPolygonZM(w, byteOrder, g.(geom.MultiPolygonZM))
	case geom.GeometryCollection:
		return writeGeometryCollection(w, byteOrder, g.(geom.GeometryCollection))
	case geom.GeometryCollectionZ:
		return writeGeometryCollectionZ(w, byteOrder, g.(geom.GeometryCollectionZ))
	case geom.GeometryCollectionM:
		return writeGeometryCollectionM(w, byteOrder, g.(geom.GeometryCollectionM))
	case geom.GeometryCollectionZM:
		return writeGeometryCollectionZM(w, byteOrder, g.(geom.GeometryCollectionZM))
	default:
		return &UnsupportedGeometryError{reflect.TypeOf(g)}
	}
}

func Encode(g geom.T, byteOrder binary.ByteOrder) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	if err := Write(w, byteOrder, g); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
