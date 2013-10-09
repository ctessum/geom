package main

import (
	"log"
	"os"
	"text/template"
)

type dim struct {
	Z string
	M string
}

var dims = []dim{
	{Z: "", M: ""},
	{Z: "Z", M: ""},
	{Z: "", M: "M"},
	{Z: "Z", M: "M"},
}

const point = `package geom
{{range .Instances}}
type Point{{.Z}}{{.M}} struct {
	X float64
	Y float64{{with .Z}}
	Z float64{{end}}{{with .M}}
	M float64{{end}}
}

func (point{{.Z}}{{.M}} Point{{.Z}}{{.M}}) Bounds() *Bounds {
	return NewBoundsPoint{{.Z}}{{.M}}(point{{.Z}}{{.M}})
}
{{end}}`

const pointWKB = `package wkb

import (
	"encoding/binary"
	"github.com/twpayne/gogeom/geom"
	"io"
)
{{range .Instances}}
func point{{.Z}}{{.M}}Reader(r io.Reader, byteOrder binary.ByteOrder) (geom.T, error) {
	point{{.Z}}{{.M}} := geom.Point{{.Z}}{{.M}}{}
	if err := binary.Read(r, byteOrder, &point{{.Z}}{{.M}}); err != nil {
		return nil, err
	}
	return point{{.Z}}{{.M}}, nil
}

func writePoint{{.Z}}{{.M}}(w io.Writer, byteOrder binary.ByteOrder, point{{.Z}}{{.M}} geom.Point{{.Z}}{{.M}}) error {
	return binary.Write(w, byteOrder, &point{{.Z}}{{.M}})
}

func writePoint{{.Z}}{{.M}}s(w io.Writer, byteOrder binary.ByteOrder, point{{.Z}}{{.M}}s []geom.Point{{.Z}}{{.M}}) error {
	if err := binary.Write(w, byteOrder, uint32(len(point{{.Z}}{{.M}}s))); err != nil {
		return err
	}
	return binary.Write(w, byteOrder, &point{{.Z}}{{.M}}s)
}

func writePoint{{.Z}}{{.M}}ss(w io.Writer, byteOrder binary.ByteOrder, point{{.Z}}{{.M}}ss [][]geom.Point{{.Z}}{{.M}}) error {
	if err := binary.Write(w, byteOrder, uint32(len(point{{.Z}}{{.M}}ss))); err != nil {
		return err
	}
	for _, point{{.Z}}{{.M}}s := range point{{.Z}}{{.M}}ss {
		if err := writePoint{{.Z}}{{.M}}s(w, byteOrder, point{{.Z}}{{.M}}s); err != nil {
			return err
		}
	}
	return nil

}
{{end}}`

const lineString = `package geom
{{range .Instances}}
type LineString{{.Z}}{{.M}} struct {
	Points []Point{{.Z}}{{.M}}
}

func (lineString{{.Z}}{{.M}} LineString{{.Z}}{{.M}}) Bounds() *Bounds {
	return NewBounds().ExtendPoint{{.Z}}{{.M}}s(lineString{{.Z}}{{.M}}.Points)
}
{{end}}`

const polygon = `package geom
{{range .Instances}}
type Polygon{{.Z}}{{.M}} struct {
	Rings [][]Point{{.Z}}{{.M}}
}

func (polygon{{.Z}}{{.M}} Polygon{{.Z}}{{.M}}) Bounds() *Bounds {
	return NewBounds().ExtendPoint{{.Z}}{{.M}}ss(polygon{{.Z}}{{.M}}.Rings)
}
{{end}}`

const multiPoint = `package geom
{{range .Instances}}
type MultiPoint{{.Z}}{{.M}} struct {
	Points []Point{{.Z}}{{.M}}
}

func (multiPoint{{.Z}}{{.M}} MultiPoint{{.Z}}{{.M}}) Bounds() *Bounds {
	return NewBounds().ExtendPoint{{.Z}}{{.M}}s(multiPoint{{.Z}}{{.M}}.Points)
}
{{end}}`

var types = []struct {
	filename  string
	name      string
	template  string
	Instances []dim
}{
	{"geom/point.go", "Point", point, dims},
	{"geom/linestring.go", "LineString", lineString, dims},
	{"geom/polygon.go", "Polygon", polygon, dims},
	{"geom/multipoint.go", "MultiPoint", multiPoint, dims},
	{"geom/encoding/wkb/point.go", "Point", pointWKB, dims},
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
