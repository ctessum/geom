package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"io"
)

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

func Unmarshal(buf []byte) (geom.T, error) {
	return Read(bytes.NewBuffer(buf))
}
