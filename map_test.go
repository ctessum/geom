package gis

import (
	"github.com/twpayne/gogeom/geom"
	"image/color"
	"os"
	"testing"
)

func TestMap(t *testing.T) {
	shape := geom.T(geom.Polygon{[][]geom.Point{[]geom.Point{{1., 0.}, {2., 1.}, {1., 2.},
		{0., 1.}, {1., 0.}}}})
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	DrawShapes(f, []color.NRGBA{{0, 0, 0, 255}},
		[]color.NRGBA{{255, 255, 255, 127}}, 5, 0, shape)
	f.Close()
}
