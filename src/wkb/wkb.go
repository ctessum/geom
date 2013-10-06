package wkb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"geom"
	"io"
	"reflect"
)

const (
	wkbXDR = 0
	wkbNDR = 1
)

const (
	wkbPoint                = 1
	wkbPointM               = 2001
	wkbPointZ               = 1001
	wkbPointZM              = 3001
	wkbLineString           = 2
	wkbLineStringM          = 2002
	wkbLineStringZ          = 1002
	wkbLineStringZM         = 3002
	wkbPolygon              = 3
	wkbPolygonM             = 2003
	wkbPolygonZ             = 1003
	wkbPolygonZM            = 3003
	wkbMultiPoint           = 4
	wkbMultiPointM          = 2004
	wkbMultiPointZ          = 1004
	wkbMultiPointZM         = 3004
	wkbMultiLineString      = 5
	wkbMultiLineStringM     = 2005
	wkbMultiLineStringZ     = 1005
	wkbMultiLineStringZM    = 3005
	wkbMultiPolygon         = 6
	wkbMultiPolygonM        = 2006
	wkbMultiPolygonZ        = 1006
	wkbMultiPolygonZM       = 3006
	wkbGeometryCollection   = 7
	wkbGeometryCollectionM  = 2007
	wkbGeometryCollectionZ  = 1007
	wkbGeometryCollectionZM = 3007
	wkbPolyhedralSurface    = 15
	wkbPolyhedralSurfaceM   = 2015
	wkbPolyhedralSurfaceZ   = 1015
	wkbPolyhedralSurfaceZM  = 3015
	wkbTIN                  = 16
	wkbTINM                 = 2016
	wkbTINZ                 = 1016
	wkbTINZM                = 3016
	wkbTriangle             = 17
	wkbTriangleM            = 2017
	wkbTriangleZ            = 1017
	wkbTriangleZM           = 3017
)

var (
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "wkb: unsupported type: " + e.Type.String()
}

type wkbReader func(io.Reader, binary.ByteOrder) (geom.T, error)

var wkbReaders = map[uint32]wkbReader{
	wkbPoint:        pointReader,
	wkbPointZ:       pointZReader,
	wkbPointM:       pointMReader,
	wkbPointZM:      pointZMReader,
	wkbLineString:   lineStringReader,
	wkbLineStringZ:  lineStringZReader,
	wkbLineStringM:  lineStringMReader,
	wkbLineStringZM: lineStringZMReader,
	wkbPolygon:      polygonReader,
	wkbPolygonZ:     polygonZReader,
	wkbPolygonM:     polygonMReader,
	wkbPolygonZM:    polygonZMReader,
}

func readLinearRing(r io.Reader, byteOrder binary.ByteOrder) (geom.LinearRing, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make(geom.LinearRing, numPoints)
	if err := binary.Read(r, byteOrder, &points); err != nil {
		return nil, err
	}
	return points, nil
}

func readLinearRingZ(r io.Reader, byteOrder binary.ByteOrder) (geom.LinearRingZ, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZs := make(geom.LinearRingZ, numPoints)
	if err := binary.Read(r, byteOrder, &pointZs); err != nil {
		return nil, err
	}
	return pointZs, nil
}

func readLinearRingM(r io.Reader, byteOrder binary.ByteOrder) (geom.LinearRingM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointMs := make(geom.LinearRingM, numPoints)
	if err := binary.Read(r, byteOrder, &pointMs); err != nil {
		return nil, err
	}
	return pointMs, nil
}

func readLinearRingZM(r io.Reader, byteOrder binary.ByteOrder) (geom.LinearRingZM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZMs := make(geom.LinearRingZM, numPoints)
	if err := binary.Read(r, byteOrder, &pointZMs); err != nil {
		return nil, err
	}
	return pointZMs, nil
}

func pointReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	point := geom.Point{}
	if err := binary.Read(r, byteOrder, &point); err != nil {
		return nil, err
	}
	return point, nil
}

func pointZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZ := geom.PointZ{}
	if err := binary.Read(r, byteOrder, &pointZ); err != nil {
		return nil, err
	}
	return pointZ, nil
}

func pointMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointM := geom.PointM{}
	if err := binary.Read(r, byteOrder, &pointM); err != nil {
		return nil, err
	}
	return pointM, nil
}

func pointZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZM := geom.PointZM{}
	if err := binary.Read(r, byteOrder, &pointZM); err != nil {
		return nil, err
	}
	return pointZM, nil
}

func lineStringReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	points, err := readLinearRing(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineString{points}, nil
}

func lineStringZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZs, err := readLinearRingZ(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringZ{pointZs}, nil
}

func lineStringMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointMs, err := readLinearRingM(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringM{pointMs}, nil
}

func lineStringZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZMs, err := readLinearRingZM(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineStringZM{pointZMs}, nil
}

func polygonReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make([]geom.LinearRing, numRings)
	for i := uint32(0); i < numRings; i++ {
		if points, err := readLinearRing(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = points
		}
	}
	return geom.Polygon{rings}, nil
}

func polygonZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringZs := make([]geom.LinearRingZ, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZs, err := readLinearRingZ(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringZs[i] = pointZs
		}
	}
	return geom.PolygonZ{ringZs}, nil
}

func polygonMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringMs := make([]geom.LinearRingM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointMs, err := readLinearRingM(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringMs[i] = pointMs
		}
	}
	return geom.PolygonM{ringMs}, nil
}

func polygonZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringZMs := make([]geom.LinearRingZM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZMs, err := readLinearRingZM(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringZMs[i] = pointZMs
		}
	}
	return geom.PolygonZM{ringZMs}, nil
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
