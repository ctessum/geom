package geom

import (
	"reflect"
	"testing"
)

func TestMultipolygon_Difference(t *testing.T) {
	mp1 := MultiPolygon{Polygon{
		Path{
			Point{X: 99, Y: 164}, Point{X: 99, Y: 108}, Point{X: 114, Y: 108}, Point{X: 114, Y: 0},
			Point{X: 121, Y: 0}, Point{X: 121, Y: 164}, Point{X: 99, Y: 164},
		},
		Path{
			Point{X: 0, Y: 499}, Point{X: 0, Y: 488},
			Point{X: 88, Y: 488}, Point{X: 88, Y: 465}, Point{X: 97, Y: 465},
			Point{X: 97, Y: 326}, Point{X: 79, Y: 326}, Point{X: 79, Y: 258},
			Point{X: 121, Y: 258}, Point{X: 121, Y: 499}, Point{X: 0, Y: 499},
		},
	}}

	mp2 := MultiPolygon{Polygon{Path{
		Point{X: 114, Y: 0}, Point{X: 161, Y: 0}, Point{X: 161, Y: 168},
		Point{X: 114, Y: 168}, Point{X: 114, Y: 0},
	}}}

	difference := mp2.Difference(mp1)
	want := Polygon{Path{
		Point{X: 114, Y: 168}, Point{X: 114, Y: 164}, Point{X: 121, Y: 164},
		Point{X: 121, Y: 0}, Point{X: 161, Y: 0}, Point{X: 161, Y: 168},
		Point{X: 114, Y: 168},
	}}
	if !reflect.DeepEqual(want, difference) {
		t.Errorf("want %+v, have %+v", want, difference)
	}
}
