package shp

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ctessum/geom"
	"github.com/jonas-p/go-shp"
)

// Tag to use for matching struct fields with shapefile attributes.
// Case insensitive.
const tag = "shp"

// Reader is a wrapper around the github.com/jonas-p/go-shp shapefile
// reader.
type Reader struct {
	shp.Reader
	row          int
	fieldIndices map[string]int
	err          error
}

// DecodeRow decodes a shapefile row into a struct. The input
// value rec must be a pointer to a struct. The function will
// attempt to match the struct fields to shapefile data.
// It will read the shape data into any struct fields that
// implement the geom.T interface. It will read attribute
// data into any struct fields whose `shp` tag or field names
// that match an attribute name in the shapefile (case insensitive).
// Only exported fields will be matched, and all matched fields
// must be of either string, int, or float64 types.
// The return value is true if there are still more records
// to be read from the shapefile.
// Be sure to call r.Error() after reading is finished
// to check for any errors that may have occured.
func (r Reader) DecodeRow(rec interface{}) bool {
	run := r.Next()
	if !run || r.err != nil {
		return false
	}
	// Figure out the indices of the attribute fields
	if r.fieldIndices == nil {
		r.fieldIndices = make(map[string]int)
		for i, f := range r.Fields() {
			name := strings.ToLower(shpFieldName2String(f.Name))
			r.fieldIndices[name] = i
		}
	}
	v, t := getRecInfo(rec)
	_, shape := r.Shape()
	for i := 0; i < v.NumField(); i++ {
		fType := t.Field(i)
		fValue := v.Field(i)
		fName := strings.ToLower(fType.Name)
		tagName := strings.ToLower(fType.Tag.Get(tag))

		// First, check if this is a geometry field
		var gref geom.T
		if fType.Type.Implements(reflect.TypeOf(gref)) {
			_, g, err := shp2Geom(0, shape)
			if err != nil {
				r.err = err
				return false
			}
			fValue.Set(reflect.ValueOf(g))

			// Then, check the tag name
		} else if j, ok := r.fieldIndices[tagName]; ok {
			r.setFieldToAttribute(fValue, fType.Type, j)

			// Finally, check the struct field name
		} else if j, ok := r.fieldIndices[fName]; ok {
			r.setFieldToAttribute(fValue, fType.Type, j)
		}

	}
	r.row++
	return run
}

// Error returns any errors that have been encountered while decoding
// a shapfile.
func (r Reader) Error() error {
	return r.err
}

func getRecInfo(rec interface{}) (reflect.Value, reflect.Type) {
	t := reflect.TypeOf(rec)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("rec must be a pointer to a "+
			"struct, not a %v.", t.Kind()))
	}
	v := reflect.Indirect(reflect.ValueOf(rec))
	if tt := v.Type().Kind(); tt != reflect.Struct {
		panic(fmt.Sprintf("rec must be a struct, not a %v.", tt))
	}
	return v, v.Type()
}

// ShpFieldName2String converts the shapefile field name into a
// string that can be more easily dealt with.
func shpFieldName2String(name [11]byte) string {
	b := bytes.Trim(name[:], "\x00")
	n := bytes.Index(b, []byte{0})
	if n == -1 {
		n = len(b)
	}
	return strings.TrimSpace(string(b[0:n]))
}

// ShpAttrbute2Float converts a shapefile attribute (which may contain
// "\x00" characters to a float.
func shpAttributeToFloat(attr string) (float64, error) {
	f, err := strconv.ParseFloat(strings.Trim(attr, "\x00"), 64)
	return f, err
}

// ShpAttrbute2Int converts a shapefile attribute (which may contain
// "\x00" characters to an int.
func shpAttributeToInt(attr string) (int64, error) {
	i, err := strconv.ParseInt(strings.Trim(attr, "\x00"), 10, 64)
	return i, err
}

func (r Reader) setFieldToAttribute(fValue reflect.Value,
	fType reflect.Type, index int) {
	dataStr := r.ReadAttribute(r.row, index)
	switch fType.Kind() {
	case reflect.Float64:
		d, err := shpAttributeToFloat(dataStr)
		if err != nil {
			r.err = err
			return
		}
		fValue.SetFloat(d)
	case reflect.Int:
		d, err := shpAttributeToInt(dataStr)
		if err != nil {
			r.err = err
			return
		}
		fValue.SetInt(d)
	case reflect.String:
		fValue.SetString(dataStr)
	default:
		panic("Struct field type can only be float64, int, or string.")
	}

}
