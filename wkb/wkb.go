package geom

import (
	"encoding/binary"
	"fmt"
	"io"
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

type WKBGeom interface {
	Geom
	wkbGeometryType() uint32
	wkbWrite(io.Writer, binary.ByteOrder) error
}

func (Point) wkbGeometryType() uint32 {
	return wkbPoint
}

func (PointZ) wkbGeometryType() uint32 {
	return wkbPointZ
}

func (PointM) wkbGeometryType() uint32 {
	return wkbPointM
}

func (PointZM) wkbGeometryType() uint32 {
	return wkbPointZM
}

func (LineString) wkbGeometryType() uint32 {
	return wkbLineString
}

func (LineStringZ) wkbGeometryType() uint32 {
	return wkbLineStringZ
}

func (LineStringM) wkbGeometryType() uint32 {
	return wkbLineStringM
}

func (LineStringZM) wkbGeometryType() uint32 {
	return wkbLineStringZM
}

func (Polygon) wkbGeometryType() uint32 {
	return wkbPolygon
}

func (PolygonZ) wkbGeometryType() uint32 {
	return wkbPolygonZ
}

func (PolygonM) wkbGeometryType() uint32 {
	return wkbPolygonM
}

func (PolygonZM) wkbGeometryType() uint32 {
	return wkbPolygonZM
}


type wkbReader func(io.Reader, binary.ByteOrder) (Geom, error)

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

func readLinearRing(r io.Reader, byteOrder binary.ByteOrder) (LinearRing, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make(LinearRing, numPoints)
	if err := binary.Read(r, byteOrder, &points); err != nil {
		return nil, err
	}
	return points, nil
}

func readLinearRingZ(r io.Reader, byteOrder binary.ByteOrder) (LinearRingZ, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZs := make(LinearRingZ, numPoints)
	if err := binary.Read(r, byteOrder, &pointZs); err != nil {
		return nil, err
	}
	return pointZs, nil
}

func readLinearRingM(r io.Reader, byteOrder binary.ByteOrder) (LinearRingM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointMs := make(LinearRingM, numPoints)
	if err := binary.Read(r, byteOrder, &pointMs); err != nil {
		return nil, err
	}
	return pointMs, nil
}

func readLinearRingZM(r io.Reader, byteOrder binary.ByteOrder) (LinearRingZM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZMs := make(LinearRingZM, numPoints)
	if err := binary.Read(r, byteOrder, &pointZMs); err != nil {
		return nil, err
	}
	return pointZMs, nil
}

func pointReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	point := Point{}
	if err := binary.Read(r, byteOrder, &point); err != nil {
		return nil, err
	}
	return point, nil
}

func pointZReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	pointZ := PointZ{}
	if err := binary.Read(r, byteOrder, &pointZ); err != nil {
		return nil, err
	}
	return pointZ, nil
}

func pointMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	pointM := PointM{}
	if err := binary.Read(r, byteOrder, &pointM); err != nil {
		return nil, err
	}
	return pointM, nil
}

func pointZMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	pointZM := PointZM{}
	if err := binary.Read(r, byteOrder, &pointZM); err != nil {
		return nil, err
	}
	return pointZM, nil
}

func lineStringReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	points, err := readLinearRing(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return LineString{points}, nil
}

func lineStringZReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	pointZs, err := readLinearRingZ(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return LineStringZ{pointZs}, nil
}

func lineStringMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	pointMs, err := readLinearRingM(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return LineStringM{pointMs}, nil
}

func lineStringZMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	pointZMs, err := readLinearRingZM(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return LineStringZM{pointZMs}, nil
}

func polygonReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make([]LinearRing, numRings)
	for i := uint32(0); i < numRings; i++ {
		if points, err := readLinearRing(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = points
		}
	}
	return Polygon{rings}, nil
}

func polygonZReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringZs := make([]LinearRingZ, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZs, err := readLinearRingZ(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringZs[i] = pointZs
		}
	}
	return PolygonZ{ringZs}, nil
}

func polygonMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringMs := make([]LinearRingM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointMs, err := readLinearRingM(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringMs[i] = pointMs
		}
	}
	return PolygonM{ringMs}, nil
}

func polygonZMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ringZMs := make([]LinearRingZM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZMs, err := readLinearRingZM(r, byteOrder); err != nil {
			return nil, err
		} else {
			ringZMs[i] = pointZMs
		}
	}
	return PolygonZM{ringZMs}, nil
}

func WKBRead(r io.Reader) (Geom, error) {

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

func writeMany(w io.Writer, byteOrder binary.ByteOrder, data ...interface{}) error {
	for _, datum := range data {
		if err := binary.Write(w, byteOrder, datum); err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRing(w io.Writer, byteOrder binary.ByteOrder, linearRing []Point) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRing))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRing)
}

func writeLinearRingZ(w io.Writer, byteOrder binary.ByteOrder, linearRingZ []PointZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingZ))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRingZ)
}

func writeLinearRingM(w io.Writer, byteOrder binary.ByteOrder, linearRingM []PointM) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingM))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRingM)
}

func writeLinearRingZM(w io.Writer, byteOrder binary.ByteOrder, linearRingZM []PointZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(linearRingZM))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &linearRingZM)
}

func (point Point) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return binary.Write(w, byteOrder, &point)
}

func (pointZ PointZ) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return binary.Write(w, byteOrder, &pointZ)
}

func (pointM PointM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return binary.Write(w, byteOrder, &pointM)
}

func (pointZM PointZM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return binary.Write(w, byteOrder, &pointZM)
}

func (lineString LineString) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writeLinearRing(w, byteOrder, lineString.Points)
}

func (lineStringZ LineStringZ) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writeLinearRingZ(w, byteOrder, lineStringZ.Points)
}

func (lineStringM LineStringM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writeLinearRingM(w, byteOrder, lineStringM.Points)
}

func (lineStringZM LineStringZM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writeLinearRingZM(w, byteOrder, lineStringZM.Points)
}

func (polygon Polygon) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	if err := binary.Write(w, byteOrder, uint32(len(polygon.Rings))); err != nil {
		return err
	}
	for _, ring := range polygon.Rings {
		if err := writeLinearRing(w, byteOrder, ring); err != nil {
			return err
		}
	}
	return nil
}

func (polygonZ PolygonZ) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	if err := binary.Write(w, byteOrder, uint32(len(polygonZ.Rings))); err != nil {
		return err
	}
	for _, ring := range polygonZ.Rings {
		if err := writeLinearRingZ(w, byteOrder, ring); err != nil {
			return err
		}
	}
	return nil
}

func (polygonM PolygonM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	if err := binary.Write(w, byteOrder, uint32(len(polygonM.Rings))); err != nil {
		return err
	}
	for _, ring := range polygonM.Rings {
		if err := writeLinearRingM(w, byteOrder, ring); err != nil {
			return err
		}
	}
	return nil
}

func (polygonZM PolygonZM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	err := binary.Write(w, byteOrder, uint32(len(polygonZM.Rings)))
	if err != nil {
		return err
	}
	for _, ring := range polygonZM.Rings {
		if err := writeLinearRingZM(w, byteOrder, ring); err != nil {
			return err
		}
	}
	return nil
}

func WKBWrite(w io.Writer, byteOrder binary.ByteOrder, g WKBGeom) error {
	var wkbByteOrder uint8
	switch byteOrder {
	case XDR:
		wkbByteOrder = wkbXDR
	case NDR:
		wkbByteOrder = wkbNDR
	default:
		return fmt.Errorf("unsupported byte order %v", byteOrder)
	}
	if err := writeMany(w, byteOrder, wkbByteOrder, g.wkbGeometryType()); err != nil {
		return err
	}
	return g.wkbWrite(w, byteOrder)
}
