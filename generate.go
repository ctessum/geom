package main

import (
	"flag"
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

	var clean = flag.Bool("clean", false, "")
	flag.Parse()

	var filenames = []string{
		"geom/geom.go",
		"geom/bounds.go",
		"geom/point.go",
		"geom/linestring.go",
		"geom/polygon.go",
		"geom/multipoint.go",
		"geom/multilinestring.go",
		"geom/multipolygon.go",
		"geom/similar.go",
		"geom/encoding/wkb/wkb.go",
		"geom/encoding/wkb/point.go",
		"geom/encoding/wkb/linestring.go",
		"geom/encoding/wkb/polygon.go",
		"geom/encoding/wkb/multipoint.go",
		"geom/encoding/wkb/multilinestring.go",
		"geom/encoding/wkb/multipolygon.go",
		"geom/encoding/wkt/point.go",
		"geom/encoding/wkt/linestring.go",
		"geom/encoding/wkt/polygon.go",
	}

	if *clean {
		for _, filename := range filenames {
			os.Remove(filename)
		}
	} else {
		for _, filename := range filenames {
			if err := generate(filename); err != nil {
				log.Fatalf("%s: %s", filename, err)
			}
		}
	}

}
