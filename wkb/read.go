package wkb

import (
	"encoding/binary"
	"fmt"
	"io"
)

func readMany(r io.Reader, byteOrder binary.ByteOrder, data ...interface{}) error {
	for _, datum := range data {
		err := binary.Read(r, byteOrder, datum)
		if err != nil {
			return err
		}
	}
	return nil
}

func readPoint(r io.Reader, byteOrder binary.ByteOrder, p *Point) error {
	err := readMany(r, byteOrder, &p.X, &p.Y)
	if err != nil {
		return err
	}
	return nil
}

func readPointZ(r io.Reader, byteOrder binary.ByteOrder, p *PointZ) error {
	err := readMany(r, byteOrder, &p.X, &p.Y, &p.Z)
	if err != nil {
		return err
	}
	return nil
}

func readPointM(r io.Reader, byteOrder binary.ByteOrder, p *PointM) error {
	err := readMany(r, byteOrder, &p.X, &p.Y, &p.M)
	if err != nil {
		return err
	}
	return nil
}

func readPointZM(r io.Reader, byteOrder binary.ByteOrder, p *PointZM) error {
	err := readMany(r, byteOrder, &p.X, &p.Y, &p.Z, &p.M)
	if err != nil {
		return err
	}
	return nil
}

func readLinearRing(r io.Reader, byteOrder binary.ByteOrder) ([]Point, error) {
	var numPoints uint32
	err := binary.Read(r, byteOrder, &numPoints)
	if err != nil {
		return nil, err
	}
	points := make([]Point, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		err = readPoint(r, byteOrder, &points[i])
		if err != nil {
			return nil, err
		}
	}
	return points, nil
}

func readLinearRingZ(r io.Reader, byteOrder binary.ByteOrder) ([]PointZ, error) {
	var numPoints uint32
	err := binary.Read(r, byteOrder, &numPoints)
	if err != nil {
		return nil, err
	}
	pointZs := make([]PointZ, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		err = readPointZ(r, byteOrder, &pointZs[i])
		if err != nil {
			return nil, err
		}
	}
	return pointZs, nil
}

func readLinearRingM(r io.Reader, byteOrder binary.ByteOrder) ([]PointM, error) {
	var numPoints uint32
	err := binary.Read(r, byteOrder, &numPoints)
	if err != nil {
		return nil, err
	}
	pointMs := make([]PointM, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		err = readPointM(r, byteOrder, &pointMs[i])
		if err != nil {
			return nil, err
		}
	}
	return pointMs, nil
}

func readLinearRingZM(r io.Reader, byteOrder binary.ByteOrder) ([]PointZM, error) {
	var numPoints uint32
	err := binary.Read(r, byteOrder, &numPoints)
	if err != nil {
		return nil, err
	}
	pointZMs := make([]PointZM, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		err = readPointZM(r, byteOrder, &pointZMs[i])
		if err != nil {
			return nil, err
		}
	}
	return pointZMs, nil
}

func Read(r io.Reader) (Geom, error) {

	var wkbByteOrder uint8
	err := binary.Read(r, binary.LittleEndian, &wkbByteOrder)
	if err != nil {
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
	err = binary.Read(r, byteOrder, &wkbGeometryType)
	if err != nil {
		return nil, err
	}
	switch wkbGeometryType {
	case wkbPoint:
		point := Point{}
		err = readPoint(r, byteOrder, &point)
		if err != nil {
			return nil, err
		}
		return point, nil
	case wkbLineString:
		points, err := readLinearRing(r, byteOrder)
		if err != nil {
			return nil, err
		}
		return LineString{points}, nil
	case wkbPointZ:
		pointZ := PointZ{}
		err = readPointZ(r, byteOrder, &pointZ)
		if err != nil {
			return nil, err
		}
		return pointZ, nil
	case wkbLineStringZ:
		pointZs, err := readLinearRingZ(r, byteOrder)
		if err != nil {
			return nil, err
		}
		return LineStringZ{pointZs}, nil
	case wkbPointM:
		pointM := PointM{}
		err = readPointM(r, byteOrder, &pointM)
		if err != nil {
			return nil, err
		}
		return pointM, nil
	case wkbLineStringM:
		pointMs, err := readLinearRingM(r, byteOrder)
		if err != nil {
			return nil, err
		}
		return LineStringM{pointMs}, nil
	case wkbPointZM:
		pointZM := PointZM{}
		err = readPointZM(r, byteOrder, &pointZM)
		if err != nil {
			return nil, err
		}
		return pointZM, nil
	case wkbLineStringZM:
		pointZMs, err := readLinearRingZM(r, byteOrder)
		if err != nil {
			return nil, err
		}
		return LineStringZM{pointZMs}, nil
	default:
		return nil, fmt.Errorf("unsupported geometry type: %v", wkbGeometryType)
	}

}
