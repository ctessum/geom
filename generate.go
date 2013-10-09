package main

import (
	"log"
	"os"
	"text/template"
)

type dim struct {
	Z  string
	M  string
	ZM string
}

var dims = []dim{
	{Z: "", M: "", ZM: ""},
	{Z: "Z", M: "", ZM: "Z"},
	{Z: "", M: "M", ZM: "M"},
	{Z: "Z", M: "M", ZM: "ZM"},
}

const point = `package geom
{{range .Dims}}
type Point{{.ZM}} struct {
	X float64
	Y float64{{with .Z}}
	Z float64{{end}}{{with .M}}
	M float64{{end}}
}

func (point{{.ZM}} Point{{.ZM}}) Bounds() *Bounds {
	return NewBoundsPoint{{.ZM}}(point{{.ZM}})
}
{{end}}`

const pointWKB = `package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)
{{range .Dims}}
func point{{.ZM}}Reader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	point{{.ZM}} := geom.Point{{.ZM}}{}
	if err := binary.Read(r, byteOrder, &point{{.ZM}}); err != nil {
		return nil, err
	}
	return point{{.ZM}}, nil
}

func readPoint{{.ZM}}s(r io.Reader, byteOrder binary.ByteOrder) ([]geom.Point{{.ZM}}, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	point{{.ZM}}s := make([]geom.Point{{.ZM}}, numPoints)
	if err := binary.Read(r, byteOrder, &point{{.ZM}}s); err != nil {
		return nil, err
	}
	return point{{.ZM}}s, nil
}

func writePoint{{.ZM}}(w io.Writer, byteOrder binary.ByteOrder, point{{.ZM}} geom.Point{{.ZM}}) error {
	return binary.Write(w, byteOrder, &point{{.ZM}})
}

func writePoint{{.ZM}}s(w io.Writer, byteOrder binary.ByteOrder, point{{.ZM}}s []geom.Point{{.ZM}}) error {
	if err := binary.Write(w, byteOrder, uint32(len(point{{.ZM}}s))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &point{{.ZM}}s)
}

func writePoint{{.ZM}}ss(w io.Writer, byteOrder binary.ByteOrder, point{{.ZM}}ss [][]geom.Point{{.ZM}}) error {
	if err := binary.Write(w, byteOrder, uint32(len(point{{.ZM}}ss))); err != nil {
		return err
	}
	for _, point{{.ZM}}s := range point{{.ZM}}ss {
		if err := writePoint{{.ZM}}s(w, byteOrder, point{{.ZM}}s); err != nil {
			return err
		}
	}
	return nil

}
{{end}}`

const pointWKT = `package wkt

import (
	"github.com/twpayne/gogeom/geom"
	"strconv"
)
{{range .Dims}}
func appendPoint{{.ZM}}Coords(dst []byte, point{{.ZM}} *geom.Point{{.ZM}}) []byte {
	dst = strconv.AppendFloat(dst, point{{.ZM}}.X, 'g', -1, 64)
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, point{{.ZM}}.Y, 'g', -1, 64){{if .Z}}
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, point{{.ZM}}.Z, 'g', -1, 64){{end}}{{if .M}}
	dst = append(dst, ' ')
	dst = strconv.AppendFloat(dst, point{{.ZM}}.M, 'g', -1, 64){{end}}
	return dst
}

func appendPoint{{.ZM}}sCoords(dst []byte, point{{.ZM}}s []geom.Point{{.ZM}}) []byte {
	for i, point{{.ZM}} := range point{{.ZM}}s {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = appendPoint{{.ZM}}Coords(dst, &point{{.ZM}})
	}
	return dst
}

func appendPoint{{.ZM}}ssCoords(dst []byte, point{{.ZM}}ss [][]geom.Point{{.ZM}}) []byte {
	for i, point{{.ZM}}s := range point{{.ZM}}ss {
		if i != 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '(')
		dst = appendPoint{{.ZM}}sCoords(dst, point{{.ZM}}s)
		dst = append(dst, ')')
	}
	return dst
}

func appendPoint{{.ZM}}WKT(dst []byte, point{{.ZM}} *geom.Point{{.ZM}}) []byte {
	dst = append(dst, []byte("POINT{{.ZM}}(")...)
	dst = appendPoint{{.ZM}}Coords(dst, point{{.ZM}})
	dst = append(dst, ')')
	return dst
}
{{end}}`

const lineString = `package geom
{{range .Dims}}
type LineString{{.ZM}} struct {
	Points []Point{{.ZM}}
}

func (lineString{{.ZM}} LineString{{.ZM}}) Bounds() *Bounds {
	return NewBounds().ExtendPoint{{.ZM}}s(lineString{{.ZM}}.Points)
}
{{end}}`

const lineStringWKB = `package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)
{{range .Dims}}
func lineString{{.ZM}}Reader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	point{{.ZM}}s, err := readPoint{{.ZM}}s(r, byteOrder)
	if err != nil {
		return nil, err
	}
	return geom.LineString{{.ZM}}{point{{.ZM}}s}, nil
}

func writeLineString{{.ZM}}(w io.Writer, byteOrder binary.ByteOrder, lineString{{.ZM}} geom.LineString{{.ZM}}) error {
	return writePoint{{.ZM}}s(w, byteOrder, lineString{{.ZM}}.Points)
}
{{end}}`

const lineStringWKT = `package wkt

import (
	"github.com/twpayne/gogeom/geom"
)
{{range .Dims}}
func appendLineString{{.ZM}}WKT(dst []byte, lineString{{.ZM}} *geom.LineString{{.ZM}}) []byte {
	dst = append(dst, []byte("LINESTRING{{.ZM}}(")...)
	dst = appendPoint{{.ZM}}sCoords(dst, lineString{{.ZM}}.Points)
	dst = append(dst, ')')
	return dst
}
{{end}}`

const polygon = `package geom
{{range .Dims}}
type Polygon{{.ZM}} struct {
	Rings [][]Point{{.ZM}}
}

func (polygon{{.ZM}} Polygon{{.ZM}}) Bounds() *Bounds {
	return NewBounds().ExtendPoint{{.ZM}}ss(polygon{{.ZM}}.Rings)
}
{{end}}`

const polygonWKB = `package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)
{{range .Dims}}
func polygon{{.ZM}}Reader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numRings uint32
	if err := binary.Read(r, byteOrder, &numRings); err != nil {
		return nil, err
	}
	ring{{.ZM}}s := make([][]geom.Point{{.ZM}}, numRings)
	for i := uint32(0); i < numRings; i++ {
		if point{{.ZM}}s, err := readPoint{{.ZM}}s(r, byteOrder); err != nil {
			return nil, err
		} else {
			ring{{.ZM}}s[i] = point{{.ZM}}s
		}
	}
	return geom.Polygon{{.ZM}}{ring{{.ZM}}s}, nil
}

func writePolygon{{.ZM}}(w io.Writer, byteOrder binary.ByteOrder, polygon{{.ZM}} geom.Polygon{{.ZM}}) error {
	return writePoint{{.ZM}}ss(w, byteOrder, polygon{{.ZM}}.Rings)
}
{{end}}`

const polygonWKT = `package wkt

import (
	"github.com/twpayne/gogeom/geom"
)
{{range .Dims}}
func appendPolygon{{.ZM}}WKT(dst []byte, polygon{{.ZM}} *geom.Polygon{{.ZM}}) []byte {
	dst = append(dst, []byte("POLYGON{{.ZM}}(")...)
	dst = appendPoint{{.ZM}}ssCoords(dst, polygon{{.ZM}}.Rings)
	dst = append(dst, ')')
	return dst
}
{{end}}`

const multiPoint = `package geom
{{range .Dims}}
type MultiPoint{{.ZM}} struct {
	Points []Point{{.ZM}}
}

func (multiPoint{{.ZM}} MultiPoint{{.ZM}}) Bounds() *Bounds {
	return NewBounds().ExtendPoint{{.ZM}}s(multiPoint{{.ZM}}.Points)
}
{{end}}`

const multiPointWKB = `package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)
{{range .Dims}}
func multiPoint{{.ZM}}Reader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	var numPoints uint32
	if err := binary.Read(r, byteOrder, &numPoints); err != nil {
		return nil, err
	}
	point{{.ZM}}s := make([]geom.Point{{.ZM}}, numPoints)
	for i := uint32(0); i < numPoints; i++ {
		if g, err := Read(r); err == nil {
			var ok bool
			point{{.ZM}}s[i], ok = g.(geom.Point{{.ZM}})
			if !ok {
				return nil, &UnexpectedGeometryError{g}
			}
		} else {
			return nil, err
		}
	}
	return geom.MultiPoint{{.ZM}}{point{{.ZM}}s}, nil
}

func writeMultiPoint{{.ZM}}(w io.Writer, byteOrder binary.ByteOrder, multiPoint{{.ZM}} geom.MultiPoint{{.ZM}}) error {
	if err := binary.Write(w, byteOrder, uint32(len(multiPoint{{.ZM}}.Points))); err != nil {
		return err
	}
	for _, point{{.ZM}} := range multiPoint{{.ZM}}.Points {
		if err := Write(w, byteOrder, point{{.ZM}}); err != nil {
			return err
		}
	}
	return nil
}
{{end}}`

var types = []struct {
	filename string
	name     string
	template string
	Dims     []dim
}{
	{"geom/point.go", "Point", point, dims},
	{"geom/linestring.go", "LineString", lineString, dims},
	{"geom/polygon.go", "Polygon", polygon, dims},
	{"geom/multipoint.go", "MultiPoint", multiPoint, dims},
	{"geom/encoding/wkb/point.go", "PointWKB", pointWKB, dims},
	{"geom/encoding/wkb/linestring.go", "LineStringWKB", lineStringWKB, dims},
	{"geom/encoding/wkb/polygon.go", "PolygonWKB", polygonWKB, dims},
	{"geom/encoding/wkb/multipoint.go", "MultiPointWKB", multiPointWKB, dims},
	{"geom/encoding/wkt/point.go", "PointWKT", pointWKT, dims},
	{"geom/encoding/wkt/linestring.go", "LineStringWKT", lineStringWKT, dims},
	{"geom/encoding/wkt/polygon.go", "PolygonWKT", polygonWKT, dims},
}

func main() {

	for _, typ := range types {
		t := template.Must(template.New(typ.name).Parse(typ.template))
		if w, err := os.Create(typ.filename); err == nil {
			if err = t.Execute(w, typ); err != nil {
				log.Println("executing template:", err)
			}
		} else {
			log.Println("creating:", err)
		}
	}

}
