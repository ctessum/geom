package wkb

import (
	"encoding/binary"
	"github.com/ctessum/gogeom/geom"
	"io"
)

func pointReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
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

func pointZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZ := geom.PointZ{}
	if err := binary.Read(r, byteOrder, &pointZ); err != nil {
		return nil, err
	}
	return pointZ, nil
}

func readPointZs(r io.Reader, byteOrder binary.ByteOrder) ([]geom.PointZ, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZs := make([]geom.PointZ, numPoints)
	if err := binary.Read(r, byteOrder, &pointZs); err != nil {
		return nil, err
	}
	return pointZs, nil
}

func writePointZ(w io.Writer, byteOrder binary.ByteOrder, pointZ geom.PointZ) error {
	return binary.Write(w, byteOrder, &pointZ)
}

func writePointZs(w io.Writer, byteOrder binary.ByteOrder, pointZs []geom.PointZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointZs))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &pointZs)
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

func pointMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointM := geom.PointM{}
	if err := binary.Read(r, byteOrder, &pointM); err != nil {
		return nil, err
	}
	return pointM, nil
}

func readPointMs(r io.Reader, byteOrder binary.ByteOrder) ([]geom.PointM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointMs := make([]geom.PointM, numPoints)
	if err := binary.Read(r, byteOrder, &pointMs); err != nil {
		return nil, err
	}
	return pointMs, nil
}

func writePointM(w io.Writer, byteOrder binary.ByteOrder, pointM geom.PointM) error {
	return binary.Write(w, byteOrder, &pointM)
}

func writePointMs(w io.Writer, byteOrder binary.ByteOrder, pointMs []geom.PointM) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointMs))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &pointMs)
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

func pointZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	pointZM := geom.PointZM{}
	if err := binary.Read(r, byteOrder, &pointZM); err != nil {
		return nil, err
	}
	return pointZM, nil
}

func readPointZMs(r io.Reader, byteOrder binary.ByteOrder) ([]geom.PointZM, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	pointZMs := make([]geom.PointZM, numPoints)
	if err := binary.Read(r, byteOrder, &pointZMs); err != nil {
		return nil, err
	}
	return pointZMs, nil
}

func writePointZM(w io.Writer, byteOrder binary.ByteOrder, pointZM geom.PointZM) error {
	return binary.Write(w, byteOrder, &pointZM)
}

func writePointZMs(w io.Writer, byteOrder binary.ByteOrder, pointZMs []geom.PointZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(pointZMs))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &pointZMs)
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
