package geom

import (
	"encoding/binary"
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
