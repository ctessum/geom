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

var vars = struct {
	Dims []dim
}{dims}

func generate(filename string) error {
	t, err := template.ParseFiles(filename + ".tmpl")
	if err != nil {
		return err
	}
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := t.Execute(w, vars); err != nil {
		return err
	}
	return nil
}

func main() {

	var filenames = []string{
		"geom/bounds.go",
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

	for _, filename := range filenames {
		if err := generate(filename); err != nil {
			log.Fatalf("%s: %s", filename, err)
		}
	}

}
