package wkb

import (
	"encoding/binary"
	"fmt"
	"io"
)

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

func readLinearRing(r io.Reader, byteOrder binary.ByteOrder) ([]Point, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make([]Point, numPoints)
	if err := binary.Read(r, byteOrder, &points); err != nil {
		return nil, err
	}
	return points, nil
}

func readLinearRingZ(r io.Reader, byteOrder binary.ByteOrder) ([]PointZ, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZs := make([]PointZ, numPoints)
	if err := binary.Read(r, byteOrder, &pointZs); err != nil {
		return nil, err
	}
	return pointZs, nil
}

func readLinearRingM(r io.Reader, byteOrder binary.ByteOrder) ([]PointM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointMs := make([]PointM, numPoints)
	if err := binary.Read(r, byteOrder, &pointMs); err != nil {
		return nil, err
	}
	return pointMs, nil
}

func readLinearRingZM(r io.Reader, byteOrder binary.ByteOrder) ([]PointZM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZMs := make([]PointZM, numPoints)
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
	rings := make([]LinearRingZ, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZs, err := readLinearRingZ(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = pointZs
		}
	}
	return PolygonZ{rings}, nil
}

func polygonMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make([]LinearRingM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointMs, err := readLinearRingM(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = pointMs
		}
	}
	return PolygonM{rings}, nil
}

func polygonZMReader(r io.Reader, byteOrder binary.ByteOrder) (Geom, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	rings := make([]LinearRingZM, numRings)
	for i := uint32(0); i < numRings; i++ {
		if pointZMs, err := readLinearRingZM(r, byteOrder); err != nil {
			return nil, err
		} else {
			rings[i] = pointZMs
		}
	}
	return PolygonZM{rings}, nil
}

func Read(r io.Reader) (Geom, error) {

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
