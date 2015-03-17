package main

//go:generate go run generate.go

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
		"../geom.go",
		"../bounds.go",
		"../point.go",
		"../linestring.go",
		"../polygon.go",
		"../multipoint.go",
		"../multilinestring.go",
		"../multipolygon.go",
		"../geometrycollection.go",
		"../similar.go",
		"../encoding/wkb/wkb.go",
		"../encoding/wkb/point.go",
		"../encoding/wkb/linestring.go",
		"../encoding/wkb/polygon.go",
		"../encoding/wkb/multipoint.go",
		"../encoding/wkb/multilinestring.go",
		"../encoding/wkb/multipolygon.go",
		"../encoding/wkb/geometrycollection.go",
		"../encoding/wkt/point.go",
		"../encoding/wkt/linestring.go",
		"../encoding/wkt/multilinestring.go",
		"../encoding/wkt/polygon.go",
		"../encoding/wkt/multipolygon.go",
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
