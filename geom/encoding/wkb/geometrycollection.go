package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)

func geometryCollectionReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numGeometries uint32
	if err := binary.Read(r, byteOrder, &numGeometries); err != nil {
		return nil, err
	}
	geoms := make([]geom.Geom, numGeometries)
	for i := uint32(0); i < numGeometries; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			geoms[i], ok = g.(geom.Geom)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.GeometryCollection{geoms}, nil
}

func writeGeometryCollection(w io.Writer, byteOrder binary.ByteOrder, geometryCollection geom.GeometryCollection) error {
	if err := binary.Write(w, byteOrder, uint32(len(geometryCollection.Geoms))); err != nil {
		return err
	}
	for _, geom := range geometryCollection.Geoms {
		if err := Write(w, byteOrder, geom); err != nil {
			return err
		}
	}
	return nil
}

func geometryCollectionZReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numGeometries uint32
	if err := binary.Read(r, byteOrder, &numGeometries); err != nil {
		return nil, err
	}
	geomZs := make([]geom.GeomZ, numGeometries)
	for i := uint32(0); i < numGeometries; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			geomZs[i], ok = g.(geom.GeomZ)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.GeometryCollectionZ{geomZs}, nil
}

func writeGeometryCollectionZ(w io.Writer, byteOrder binary.ByteOrder, geometryCollectionZ geom.GeometryCollectionZ) error {
	if err := binary.Write(w, byteOrder, uint32(len(geometryCollectionZ.Geoms))); err != nil {
		return err
	}
	for _, geomZ := range geometryCollectionZ.Geoms {
		if err := Write(w, byteOrder, geomZ); err != nil {
			return err
		}
	}
	return nil
}

func geometryCollectionMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numGeometries uint32
	if err := binary.Read(r, byteOrder, &numGeometries); err != nil {
		return nil, err
	}
	geomMs := make([]geom.GeomM, numGeometries)
	for i := uint32(0); i < numGeometries; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			geomMs[i], ok = g.(geom.GeomM)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.GeometryCollectionM{geomMs}, nil
}

func writeGeometryCollectionM(w io.Writer, byteOrder binary.ByteOrder, geometryCollectionM geom.GeometryCollectionM) error {
	if err := binary.Write(w, byteOrder, uint32(len(geometryCollectionM.Geoms))); err != nil {
		return err
	}
	for _, geomM := range geometryCollectionM.Geoms {
		if err := Write(w, byteOrder, geomM); err != nil {
			return err
		}
	}
	return nil
}

func geometryCollectionZMReader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numGeometries uint32
	if err := binary.Read(r, byteOrder, &numGeometries); err != nil {
		return nil, err
	}
	geomZMs := make([]geom.GeomZM, numGeometries)
	for i := uint32(0); i < numGeometries; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			geomZMs[i], ok = g.(geom.GeomZM)
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.GeometryCollectionZM{geomZMs}, nil
}

func writeGeometryCollectionZM(w io.Writer, byteOrder binary.ByteOrder, geometryCollectionZM geom.GeometryCollectionZM) error {
	if err := binary.Write(w, byteOrder, uint32(len(geometryCollectionZM.Geoms))); err != nil {
		return err
	}
	for _, geomZM := range geometryCollectionZM.Geoms {
		if err := Write(w, byteOrder, geomZM); err != nil {
			return err
		}
	}
	return nil
}
