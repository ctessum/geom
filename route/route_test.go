package route

import (
	"fmt"
	"testing"

	"github.com/ctessum/geom"
)

var (
	link1 = geom.LineString{
		geom.Point{X: 0, Y: 0},
		geom.Point{X: 0, Y: 1},
		geom.Point{X: 1, Y: 1},
		geom.Point{X: 8, Y: 1},
		geom.Point{X: 8, Y: 4},
	}
	link3 = geom.LineString{
		geom.Point{X: 7.999999999999998, Y: 4}, // floating point error
		geom.Point{X: 8, Y: -6},
	}
)

func ExampleRoute() {
	link1 := geom.LineString{
		geom.Point{X: 0, Y: 0},
		geom.Point{X: 0, Y: 1},
		geom.Point{X: 1, Y: 1},
		geom.Point{X: 8, Y: 1},
		geom.Point{X: 8, Y: 4},
	} // Distance = 12
	link2 := geom.LineString{
		geom.Point{X: 8, Y: 4}, // link2 beginning == link1 end
		geom.Point{X: 8, Y: -6},
	} // Distance = 10
	startingPoint := geom.Point{X: 0, Y: -1}
	endingPoint := geom.Point{X: 6, Y: -6}
	net := NewNetwork(Time)
	net.AddLink(link1, 6)
	net.AddLink(link2, 2)
	route, distance, time, startDistance, endDistance :=
		net.ShortestRoute(startingPoint, endingPoint)
	fmt.Println(route, distance, time, startDistance, endDistance)
	// Output: [[{0 0} {0 1} {1 1} {8 1} {8 4}] [{8 4} {8 -6}]] 22 7 1 2
}

func TestFloatingPoint(t *testing.T) {
	net := NewNetwork(Time)
	net.AddLink(link1, 6)
	net.AddLink(link3, 2)
	startingPoint := geom.Point{X: 0, Y: -1}
	endingPoint := geom.Point{X: 6, Y: -6}
	route, _, _, _, _ := net.ShortestRoute(startingPoint, endingPoint)
	want := "[[{0 0} {0 1} {1 1} {8 1} {8 4}] [{7.999999999999998 4} {8 -6}]]"
	got := fmt.Sprint(route)
	if want != got {
		t.Errorf("Want: %v; got: %v", want, got)
	}
}
