package carto

import (
	"github.com/twpayne/gogeom/geom"
	"image/color"
	"math/rand"
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

func TestCmapPos(t *testing.T) {
	x := make([]float64, 100)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Float64()
	}
	c := NewColorMap(LinCutoff)
	c.AddArray(x)
	c.Set()
	f, err := os.Create("legendTestPos.svg")
	if err != nil {
		panic(err)
	}
	err = c.Legend(f, "Units!")
	if err != nil {
		panic(err)
	}
	f.Close()
}

func TestCmapNeg(t *testing.T) {
	x := make([]float64, 100)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Float64()-1
	}
	c := NewColorMap(LinCutoff)
	c.AddArray(x)
	c.Set()
	f, err := os.Create("legendTestNeg.svg")
	if err != nil {
		panic(err)
	}
	err = c.Legend(f, "Units!")
	if err != nil {
		panic(err)
	}
	f.Close()
}

func TestCmapPosNeg(t *testing.T) {
	x := make([]float64, 100)
	for i := 0; i < len(x); i++ {
		x[i] = (rand.Float64() - 0.5) * 64.
	}
	c := NewColorMap(LinCutoff)
	c.AddArray(x)
	c.Set()
	f, err := os.Create("legendTestPosNeg.svg")
	if err != nil {
		panic(err)
	}
	err = c.Legend(f, "Units!")
	if err != nil {
		panic(err)
	}
	f.Close()
}
