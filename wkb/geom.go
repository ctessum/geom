package wkb

import (
	"encoding/binary"
	"io"
)

type Geom interface {
	wkbGeometryType() uint32
	wkbWrite(io.Writer, binary.ByteOrder) error
}

type Point struct {
	X float64
	Y float64
}

func (Point) wkbGeometryType() uint32 {
	return wkbPoint
}

type PointZ struct {
	X float64
	Y float64
	Z float64
}

func (PointZ) wkbGeometryType() uint32 {
	return wkbPointZ
}

type PointM struct {
	X float64
	Y float64
	M float64
}

func (PointM) wkbGeometryType() uint32 {
	return wkbPointM
}

type PointZM struct {
	X float64
	Y float64
	Z float64
	M float64
}

func (PointZM) wkbGeometryType() uint32 {
	return wkbPointZM
}

type LineString struct {
	Points []Point
}

func (LineString) wkbGeometryType() uint32 {
	return wkbLineString
}

type LineStringZ struct {
	Points []PointZ
}

func (LineStringZ) wkbGeometryType() uint32 {
	return wkbLineStringZ
}

type LineStringM struct {
	Points []PointM
}

func (LineStringM) wkbGeometryType() uint32 {
	return wkbLineStringM
}

type LineStringZM struct {
	Points []PointZM
}

func (LineStringZM) wkbGeometryType() uint32 {
	return wkbLineStringZM
}
