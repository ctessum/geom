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

type LinearRing []Point
type LinearRingZ []PointZ
type LinearRingM []PointM
type LinearRingZM []PointZM

type LineString struct {
	Points LinearRing
}

func (LineString) wkbGeometryType() uint32 {
	return wkbLineString
}

type LineStringZ struct {
	Points LinearRingZ
}

func (LineStringZ) wkbGeometryType() uint32 {
	return wkbLineStringZ
}

type LineStringM struct {
	Points LinearRingM
}

func (LineStringM) wkbGeometryType() uint32 {
	return wkbLineStringM
}

type LineStringZM struct {
	Points LinearRingZM
}

func (LineStringZM) wkbGeometryType() uint32 {
	return wkbLineStringZM
}

type Polygon struct {
	Rings []LinearRing
}

func (Polygon) wkbGeometryType() uint32 {
	return wkbPolygon
}

type PolygonZ struct {
	Rings []LinearRingZ
}

func (PolygonZ) wkbGeometryType() uint32 {
	return wkbPolygonZ
}

type PolygonM struct {
	Rings []LinearRingM
}

func (PolygonM) wkbGeometryType() uint32 {
	return wkbPolygonM
}

type PolygonZM struct {
	Rings []LinearRingZM
}

func (PolygonZM) wkbGeometryType() uint32 {
	return wkbPolygonZM
}
