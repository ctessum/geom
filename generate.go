package main

import (
	"io"
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

func main() {

	var filenames = []string{
		"geom/point.go",
		"geom/linestring.go",
		"geom/polygon.go",
		"geom/multipoint.go",
		"geom/encoding/wkb/point.go",
		"geom/encoding/wkb/linestring.go",
		"geom/encoding/wkb/polygon.go",
		"geom/encoding/wkb/multipoint.go",
		"geom/encoding/wkt/point.go",
		"geom/encoding/wkt/linestring.go",
		"geom/encoding/wkt/polygon.go",
	}

	var vars = struct {
		Dims []dim
	}{dims}

	for _, filename := range filenames {
		t, err := template.ParseFiles(filename + ".tmpl")
		if err != nil {
			log.Fatalf("parsing template %s: %s", filename, err)
		}
		var w io.WriteCloser
		if w, err = os.Create(filename); err != nil {
			log.Fatalf("creating %s: %s", filename, err)
		}
		defer w.Close()
		if err := t.Execute(w, vars); err != nil {
			log.Fatalf("executing %s: %s", filename, err)
		}
	}

}
