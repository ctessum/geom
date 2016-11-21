package geom

import (
	"reflect"
	"testing"
)

func TestPoint_Buffer(t *testing.T) {
	p := Point{X: 1, Y: 1}
	buf := p.Buffer(1, 4)
	want := Polygon{[]Point{
		Point{X: 2, Y: 1},
		Point{X: 1, Y: 2},
		Point{X: 0, Y: 1.0000000000000002},
		Point{X: 0.9999999999999998, Y: 0},
	}}
	if !reflect.DeepEqual(buf, want) {
		t.Errorf("have: %#v, want %#v", buf, want)
	}
}
