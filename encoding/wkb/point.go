package wkb

import (
	"encoding/binary"
	"github.com/ctessum/geom"
	"io"
)

func pointReader(r io.Reader, byteOrder binary.ByteOrder) (geom.Geom, error) {
	point := geom.Point{}
	if err := binary.Read(r, byteOrder, &point); err != nil {
		return nil, err
	}
	return point, nil
}

func readPoints(r io.Reader, byteOrder binary.ByteOrder) ([]geom.Point, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	points := make([]geom.Point, numPoints)
	if err := binary.Read(r, byteOrder, &points); err != nil {
		return nil, err
	}
	return points, nil
}

func writePoint(w io.Writer, byteOrder binary.ByteOrder, point geom.Point) error {
	return binary.Write(w, byteOrder, &point)
}

func writePoints(w io.Writer, byteOrder binary.ByteOrder, points []geom.Point) error {
	if err := binary.Write(w, byteOrder, uint32(len(points))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &points)
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
