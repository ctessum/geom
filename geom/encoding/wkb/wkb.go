package wkb

import (
	"encoding/binary"
	"fmt"
	"github.com/twpayne/gogeom/geom"
	"reflect"
)

const (
	wkbXDR = 0
	wkbNDR = 1
)

const (
	wkbPoint                = 1
	wkbPointZ               = 1001
	wkbPointM               = 2001
	wkbPointZM              = 3001
	wkbLineString           = 2
	wkbLineStringZ          = 1002
	wkbLineStringM          = 2002
	wkbLineStringZM         = 3002
	wkbPolygon              = 3
	wkbPolygonZ             = 1003
	wkbPolygonM             = 2003
	wkbPolygonZM            = 3003
	wkbMultiPoint           = 4
	wkbMultiPointZ          = 1004
	wkbMultiPointM          = 2004
	wkbMultiPointZM         = 3004
	wkbMultiLineString      = 5
	wkbMultiLineStringZ     = 1005
	wkbMultiLineStringM     = 2005
	wkbMultiLineStringZM    = 3005
	wkbMultiPolygon         = 6
	wkbMultiPolygonZ        = 1006
	wkbMultiPolygonM        = 2006
	wkbMultiPolygonZM       = 3006
	wkbGeometryCollection   = 7
	wkbGeometryCollectionZ  = 1007
	wkbGeometryCollectionM  = 2007
	wkbGeometryCollectionZM = 3007
	wkbPolyhedralSurface    = 15
	wkbPolyhedralSurfaceZ   = 1015
	wkbPolyhedralSurfaceM   = 2015
	wkbPolyhedralSurfaceZM  = 3015
	wkbTIN                  = 16
	wkbTINZ                 = 1016
	wkbTINM                 = 2016
	wkbTINZM                = 3016
	wkbTriangle             = 17
	wkbTriangleZ            = 1017
	wkbTriangleM            = 2017
	wkbTriangleZM           = 3017
)

var (
	XDR = binary.BigEndian
	NDR = binary.LittleEndian
)

type UnexpectedGeometryError struct {
	Geom geom.T
}

func (e UnexpectedGeometryError) Error() string {
	return fmt.Sprintf("wkb: unexpected geometry %v", e.Geom)
}

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "wkb: unsupported type: " + e.Type.String()
}
