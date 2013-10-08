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

func writePoints(w io.Writer, byteOrder binary.ByteOrder, points []geom.Point) error {
	if err := binary.Write(w, byteOrder, uint32(len(points))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &points)
}

func writePointZs(w io.Writer, byteOrder binary.ByteOrder, pointZs []geom.PointZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointZs))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &pointZs)
}

func writePointMs(w io.Writer, byteOrder binary.ByteOrder, pointMs []geom.PointM) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointMs))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &pointMs)
}

func writePointZMs(w io.Writer, byteOrder binary.ByteOrder, pointZMs []geom.PointZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointZMs))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &pointZMs)
}

func writePointss(w io.Writer, byteOrder binary.ByteOrder, pointss [][]geom.Point) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointss))); err != nil {
		return err
	}
	for _, points := range pointss {
		if err := writePoints(w, byteOrder, points); err != nil {
			return err
		}
	}
	return nil
}

func writePointZss(w io.Writer, byteOrder binary.ByteOrder, pointZss [][]geom.PointZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointZss))); err != nil {
		return err
	}
	for _, pointZs := range pointZss {
		if err := writePointZs(w, byteOrder, pointZs); err != nil {
			return err
		}
	}
	return nil
}

func writePointMss(w io.Writer, byteOrder binary.ByteOrder, pointMss [][]geom.PointM) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointMss))); err != nil {
		return err
	}
	for _, pointMs := range pointMss {
		if err := writePointMs(w, byteOrder, pointMs); err != nil {
			return err
		}
	}
	return nil
}

func writePointZMss(w io.Writer, byteOrder binary.ByteOrder, pointZMss [][]geom.PointZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointZMss))); err != nil {
		return err
	}
	for _, pointZMs := range pointZMss {
		if err := writePointZMs(w, byteOrder, pointZMs); err != nil {
			return err
		}
	}
	return nil
}

func writeLineString(w io.Writer, byteOrder binary.ByteOrder, lineString geom.LineString) error {
	return writePoints(w, byteOrder, lineString.Points)
}

func writeLineStringZ(w io.Writer, byteOrder binary.ByteOrder, lineStringZ geom.LineStringZ) error {
	return writePointZs(w, byteOrder, lineStringZ.Points)
}

func writeLineStringM(w io.Writer, byteOrder binary.ByteOrder, lineStringM geom.LineStringM) error {
	return writePointMs(w, byteOrder, lineStringM.Points)
}

func writeLineStringZM(w io.Writer, byteOrder binary.ByteOrder, lineStringZM geom.LineStringZM) error {
	return writePointZMs(w, byteOrder, lineStringZM.Points)
}

func writePolygon(w io.Writer, byteOrder binary.ByteOrder, polygon geom.Polygon) error {
	return writePointss(w, byteOrder, polygon.Rings)
}

func writePolygonZ(w io.Writer, byteOrder binary.ByteOrder, polygonZ geom.PolygonZ) error {
	return writePointZss(w, byteOrder, polygonZ.Rings)
}

func writePolygonM(w io.Writer, byteOrder binary.ByteOrder, polygonM geom.PolygonM) error {
	return writePointMss(w, byteOrder, polygonM.Rings)
}

func writePolygonZM(w io.Writer, byteOrder binary.ByteOrder, polygonZM geom.PolygonZM) error {
	return writePointZMss(w, byteOrder, polygonZM.Rings)
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