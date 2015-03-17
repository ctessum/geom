package carto

import (
	"github.com/ctessum/geom"
	"image/color"
	"math/rand"
	"os"
	"testing"
)

func TestMap(t *testing.T) {
	shape := geom.T(geom.Polygon([][]geom.Point{[]geom.Point{
		{1., 0.}, {2., 1.}, {1., 2.},
		{0., 1.}, {1., 0.}}}))
	f, err := os.Create("test.png")
	if err != nil {
		panic(err)
	}
	DrawShapes(f, []color.NRGBA{{0, 0, 0, 255}},
		[]color.NRGBA{{0, 255, 0, 127}}, 5, 0, shape)
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
	f, err := os.Create("legendTestPos.pdf")
	if err != nil {
		panic(err)
	}
	canvas := NewDefaultLegendCanvas()
	err = c.Legend(&canvas.Canvas, "Units!")
	if err != nil {
		panic(err)
	}
	err = canvas.WriteTo(f)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func TestCmapNeg(t *testing.T) {
	x := make([]float64, 100)
	for i := 0; i < len(x); i++ {
		x[i] = rand.Float64() - 1
	}
	c := NewColorMap(LinCutoff)
	c.AddArray(x)
	c.Set()
	f, err := os.Create("legendTestNeg.pdf")
	if err != nil {
		panic(err)
	}
	canvas := NewDefaultLegendCanvas()
	err = c.Legend(&canvas.Canvas, "Units!")
	if err != nil {
		panic(err)
	}
	err = canvas.WriteTo(f)
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
	f, err := os.Create("legendTestPosNeg.pdf")
	if err != nil {
		panic(err)
	}
	canvas := NewDefaultLegendCanvas()
	err = c.Legend(&canvas.Canvas, "Units!")
	if err != nil {
		panic(err)
	}
	err = canvas.WriteTo(f)
	if err != nil {
		panic(err)
	}
	f.Close()
}

func TestInterpolate(t *testing.T) {
	c := Optimized.interpolate(0)
	white := color.NRGBA{255, 255, 255, 255}
	if c != white {
		t.Errorf("Color %v should be white", c)
	}
	c = Optimized.interpolate(-1)
	blue := color.NRGBA{59, 76, 192, 255}
	if c != blue {
		t.Errorf("Color %v should be 59,76,192", c)
	}
	c = Optimized.interpolate(1)
	red := color.NRGBA{180, 4, 38, 255}
	if c != red {
		t.Errorf("Color %v should be 180,4,38", c)
	}
	c = Optimized.interpolate(-0.92)
	blue2 := color.NRGBA{71, 94, 207, 255}
	if c != blue2 {
		t.Errorf("Color %v should be 71,94,207", c)
	}
}
