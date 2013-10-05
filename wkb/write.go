package wkb

import (
	"encoding/binary"
	"fmt"
	"io"
)

func writeMany(w io.Writer, byteOrder binary.ByteOrder, data ...interface{}) error {
	for _, datum := range data {
		err := binary.Write(w, byteOrder, datum)
		if err != nil {
			return err
		}
	}
	return nil
}

func writePoint(w io.Writer, byteOrder binary.ByteOrder, point Point) error {
	return writeMany(w, byteOrder, point.X, point.Y)
}

func writePointZ(w io.Writer, byteOrder binary.ByteOrder, pointZ PointZ) error {
	return writeMany(w, byteOrder, pointZ.X, pointZ.Y, pointZ.Z)
}

func writePointM(w io.Writer, byteOrder binary.ByteOrder, pointM PointM) error {
	return writeMany(w, byteOrder, pointM.X, pointM.Y, pointM.M)
}

func writePointZM(w io.Writer, byteOrder binary.ByteOrder, pointZM PointZM) error {
	return writeMany(w, byteOrder, pointZM.X, pointZM.Y, pointZM.Z, pointZM.M)
}

func writeLinearRing(w io.Writer, byteOrder binary.ByteOrder, linearRing []Point) error {
	binary.Write(w, byteOrder, uint32(len(linearRing)))
	for _, point := range linearRing {
		err := writePoint(w, byteOrder, point)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRingZ(w io.Writer, byteOrder binary.ByteOrder, linearRingZ []PointZ) error {
	binary.Write(w, byteOrder, uint32(len(linearRingZ)))
	for _, pointZ := range linearRingZ {
		err := writePointZ(w, byteOrder, pointZ)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRingM(w io.Writer, byteOrder binary.ByteOrder, linearRingM []PointM) error {
	binary.Write(w, byteOrder, uint32(len(linearRingM)))
	for _, pointM := range linearRingM {
		err := writePointM(w, byteOrder, pointM)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeLinearRingZM(w io.Writer, byteOrder binary.ByteOrder, linearRingZM []PointZM) error {
	binary.Write(w, byteOrder, uint32(len(linearRingZM)))
	for _, pointZM := range linearRingZM {
		err := writePointZM(w, byteOrder, pointZM)
		if err != nil {
			return err
		}
	}
	return nil
}

func (point Point) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writePoint(w, byteOrder, point)
}

func (pointZ PointZ) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writePointZ(w, byteOrder, pointZ)
}

func (pointM PointM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writePointM(w, byteOrder, pointM)
}

func (pointZM PointZM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	return writePointZM(w, byteOrder, pointZM)
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
	err := binary.Write(w, byteOrder, uint32(len(polygon.Rings)))
	if err != nil {
		return err
	}
	for _, ring := range polygon.Rings {
		err = writeLinearRing(w, byteOrder, ring)
		if err != nil {
			return err
		}
	}
	return nil
}

func (polygonZ PolygonZ) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	err := binary.Write(w, byteOrder, uint32(len(polygonZ.Rings)))
	if err != nil {
		return err
	}
	for _, ring := range polygonZ.Rings {
		err = writeLinearRingZ(w, byteOrder, ring)
		if err != nil {
			return err
		}
	}
	return nil
}

func (polygonM PolygonM) wkbWrite(w io.Writer, byteOrder binary.ByteOrder) error {
	err := binary.Write(w, byteOrder, uint32(len(polygonM.Rings)))
	if err != nil {
		return err
	}
	for _, ring := range polygonM.Rings {
		err = writeLinearRingM(w, byteOrder, ring)
		if err != nil {
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
		err = writeLinearRingZM(w, byteOrder, ring)
		if err != nil {
			return err
		}
	}
	return nil
}

func Write(w io.Writer, byteOrder binary.ByteOrder, g Geom) error {
	var wkbByteOrder uint8
	switch byteOrder {
	case XDR:
		wkbByteOrder = wkbXDR
	case NDR:
		wkbByteOrder = wkbNDR
	default:
		return fmt.Errorf("unsupported byte order %v", byteOrder)
	}
	err := writeMany(w, byteOrder, wkbByteOrder, g.wkbGeometryType())
	if err != nil {
		return err
	}
	return g.wkbWrite(w, byteOrder)
}
