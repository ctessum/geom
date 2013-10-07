package wkt

import (
	"reflect"
)

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func (e UnsupportedGeometryError) Error() string {
	return "wkt: unsupported geometry type: " + e.Type.String()
}
