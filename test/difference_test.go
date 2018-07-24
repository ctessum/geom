package test

import (
	"fmt"
	"testing"

	"math"

	"github.com/ctessum/geom"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

type Rectangle struct {
	Width, Depth float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GenFloat64F() gopter.Gen {
	return gen.Float64().Map(func(x float64) float64 { return math.Floor(x) })
}

func GenFloat64RangeF(min float64, max float64) gopter.Gen {
	return gen.Float64Range(min, max).Map(func(x float64) float64 { return math.Floor(x) })
}

func ClosedPolygon(points []geom.Point) geom.MultiPolygon {
	if points[0] != points[len(points)-1] {
		points = append(points, points[0])
	}

	return geom.MultiPolygon([]geom.Polygon{geom.Polygon{points}})
}

func newRectangle(width float64, depth float64) Rectangle {
	return Rectangle{Width: math.Floor(width), Depth: math.Floor(depth)}
}

func deconstructRectangle(d Rectangle) (float64, float64) {
	return d.Width, d.Depth
}

func GenBigRectangle() gopter.Gen {
	return gopter.DeriveGen(
		newRectangle,
		deconstructRectangle,
		GenFloat64RangeF(100, 300),
		GenFloat64RangeF(400, 1500),
	)
}

func GenRectangle() gopter.Gen {
	return gopter.DeriveGen(
		newRectangle,
		deconstructRectangle,
		GenFloat64RangeF(10, 200),
		GenFloat64RangeF(10, 200),
	)
}

func getOrderedPoints(points []geom.Point, bounds *geom.Bounds) []geom.Point {
	var comparator utils.Comparator = func(p, q interface{}) int {
		a := p.(geom.Point)
		b := q.(geom.Point)
		r0 := a.Y - b.Y
		if r0 == 0 {
			r0 = a.X - b.X
		}
		return int(r0)
	}

	treepoints := treemap.NewWith(comparator)

	for _, p := range points {
		v, found := treepoints.Get(p)
		if found {
			treepoints.Put(p, v.(int)+1)
		} else {
			treepoints.Put(p, 1)
		}
	}

	topRight := bounds.Max
	result := []geom.Point{}
	for !treepoints.Empty() {
		ip, _ := treepoints.Min()
		p := ip.(geom.Point)
		if p.X < topRight.X && p.Y < topRight.Y {
			result = append(result, p)
		}
		treepoints.Remove(p)
	}

	return result
}

func GenCuttingRectangles() gopter.Gen {
	genParams := gopter.DefaultGenParameters()
	genParams.MaxSize = 10
	genParams.MinSize = 1
	sliceGen := gen.SliceOf(GenRectangle())
	sliceGen(genParams)
	return sliceGen
}

func polygonFromRectangle(rect Rectangle, position geom.Point) geom.MultiPolygon {
	rectPoints := []geom.Point{
		geom.Point{X: position.X, Y: position.Y},
		geom.Point{X: position.X + rect.Width, Y: position.Y},
		geom.Point{X: position.X + rect.Width, Y: position.Y + rect.Depth},
		geom.Point{X: position.X, Y: position.Y + rect.Depth},
	}
	return ClosedPolygon(rectPoints)
}

func getPoints(shape geom.Geom) []geom.Point {
	pointsFunc := shape.Points()
	var points []geom.Point

	for i := 0; i < shape.Len(); i++ {
		p := pointsFunc()
		points = append(points, p)
	}

	return points
}

func CutShapeOk(bigRectangle Rectangle, cuttingRectangles []Rectangle) bool {
	shape := polygonFromRectangle(bigRectangle, geom.Point{X: 0, Y: 0})
	for _, rect := range cuttingRectangles {
		for _, point := range getOrderedPoints(getPoints(shape), shape.Bounds()) {
			polygon := polygonFromRectangle(rect, point)
			difference := polygon.Difference(shape)
			intersection := polygon.Intersection(shape)
			if math.Round(intersection.Area()+difference.Area()) != polygon.Area() {
				// Bug with difference or intersection
				fmt.Printf("------ Bad difference -----\n")
				fmt.Printf("polygon1 %+#v\n", polygon)
				fmt.Printf("polygon2 %+v\n", shape)
				fmt.Printf("polygon1 \\ polygon2 = %+v\n", difference)
				fmt.Printf("(%v = %v + %v) != %v\n", intersection.Area()+difference.Area(), intersection.Area(), difference.Area(), polygon.Area())
				fmt.Printf("---------------------------\n")
				return false
			}
			if difference.Len() == 0 {
				newShape := geom.MultiPolygon([]geom.Polygon{shape.Difference(polygon)})
				if math.Round(newShape.Area()+polygon.Area()) != shape.Area() {
					fmt.Printf("Bad difference\n")
					fmt.Printf("polygon1 %+v\n", shape)
					fmt.Printf("polygon2 %+v\n", polygon)
					fmt.Printf("polygon1 \\ polygon2 = %+v\n", newShape)
					fmt.Printf("(%v = %v + %v) != %v\n", newShape.Area()+polygon.Area(), newShape.Area(), polygon.Area(), shape.Area())
					return false
				}
				shape = newShape
			}
		}
	}
	return true
}

func Test_RectangleProps(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234)
	properties := gopter.NewProperties(parameters)

	properties.Property("Test difference", prop.ForAll(
		CutShapeOk,
		GenBigRectangle().WithLabel("big rectangle"),
		GenCuttingRectangles().WithLabel("cutting rectangles"),
	))

	properties.TestingRun(t)
}
