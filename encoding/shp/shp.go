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

const (
	// intLength is the integer length to use when creating shapefiles
	intLength = 10

	// floatLength is the float length to use when creating shapefiles
	floatLength = 10

	// floatPrecision is the float precision to use when creating shapefiles
	floatPrecision = 10

	// stringLength is the length of the string to use when creating shapefiles
	stringLength = 50
)

// Decoder is a wrapper around the github.com/jonas-p/go-shp shapefile
// reader.
type Decoder struct {
	shp.Reader
	row          int
	fieldIndices map[string]int
	err          error
}

func NewDecoder(filename string) (*Decoder, error) {
	r := new(Decoder)
	rr, err := shp.Open(filename)
	r.Reader = *rr
	return r, err
}

func (r *Decoder) Close() {
	r.Reader.Close()
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
func (r *Decoder) DecodeRow(rec interface{}) bool {
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
		gI := reflect.TypeOf((*geom.T)(nil)).Elem()
		if fType.Type.Implements(gI) {
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
func (r Decoder) Error() error {
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

func (r Decoder) setFieldToAttribute(fValue reflect.Value,
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

// Encode is a wrapper around the github.com/jonas-p/go-shp shapefile
// reader.
type Encoder struct {
	shp.Writer
	fieldIndices []int
	geomIndex    int
	row          int
}

// NewEncoder creates a new encoder using the path to the output shapefile
// and a data archetype which is a struct whose fields will become the
// fields in the output shapefile. The archetype struct must also contain
// a field that holds a concrete geometry type by which to set the shape type
// in the output shapefile.
func NewEncoder(filename string, archetype interface{}) (*Encoder, error) {
	var err error
	e := new(Encoder)

	t := reflect.TypeOf(archetype)
	if t.Kind() != reflect.Struct {
		panic("Archetype must be a struct")
	}

	var shpType shp.ShapeType
	var shpFields []shp.Field
	for i := 0; i < t.NumField(); i++ {
		sField := t.Field(i)
		switch sField.Type.Kind() {
		case reflect.Int:
			shpFields = append(shpFields, shp.NumberField(sField.Name, intLength))
			e.fieldIndices = append(e.fieldIndices, i)
		case reflect.Float64:
			shpFields = append(shpFields,
				shp.FloatField(sField.Name, floatLength, floatPrecision))
			e.fieldIndices = append(e.fieldIndices, i)
		case reflect.String:
			shpFields = append(shpFields,
				shp.StringField(sField.Name, stringLength))
			e.fieldIndices = append(e.fieldIndices, i)
		case reflect.Struct, reflect.Slice:
			switch sField.Name {
			case "Point":
				shpType = shp.POINT
				e.geomIndex = i
			case "LineString":
				shpType = shp.POLYLINE
				e.geomIndex = i
			case "Polygon":
				shpType = shp.POLYGON
				e.geomIndex = i
			case "MultiPoint":
				shpType = shp.MULTIPOINT
				e.geomIndex = i
			case "PointZ":
				shpType = shp.POINTZ
				e.geomIndex = i
			case "LineStringZ":
				shpType = shp.POLYLINEZ
				e.geomIndex = i
			case "PolygonZ":
				shpType = shp.POLYGONZ
				e.geomIndex = i
			case "MultiPolygonZ":
				shpType = shp.MULTIPOINTZ
				e.geomIndex = i
			case "PointM":
				shpType = shp.POINTM
				e.geomIndex = i
			case "LineStringM":
				shpType = shp.POLYLINEM
				e.geomIndex = i
			case "PolygonM":
				shpType = shp.POLYGONM
				e.geomIndex = i
			case "MultiPointM":
				shpType = shp.MULTIPOINTM
				e.geomIndex = i
				//shpType = shp.MULTIPATCH
			}
		default:
			panic(fmt.Sprintf("Invalid type `%v` for field `%v`.",
				sField.Type.Kind(), sField.Name))
		}
	}
	if shpType == shp.NULL {
		panic("Did not find a shape field in the archetype struct")
	}

	w, err := shp.Create(filename, shpType)
	if err != nil {
		return nil, err
	}
	e.Writer = *w
	e.Writer.SetFields(shpFields)
	return e, nil
}

func (e *Encoder) Close() {
	e.Writer.Close()
}

// Encode encodes the data in a struct as a shapefile record.
// d must be of the same type as the archetype struct that was used to
// initialize the encoder.
func (e *Encoder) Encode(d interface{}) error {
	v := reflect.ValueOf(d)
	for i, j := range e.fieldIndices {
		e.Writer.WriteAttribute(e.row, i, v.Field(j).Interface())
	}

	shape, err := geom2Shp(v.Field(e.geomIndex).Interface().(geom.T))
	if err != nil {
		return err
	}
	e.Writer.Write(shape)
	e.row++
	return nil
}
