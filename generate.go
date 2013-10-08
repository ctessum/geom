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

var types = []struct {
	filename  string
	name      string
	template  string
	Instances []dim
}{
	{"geom/point.go", "Point", point, dims},
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
