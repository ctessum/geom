package geom

import (
	"reflect"
	"testing"
	"time"
)

func TestSimplify(t *testing.T) {
	type simplifyTest struct {
		input     LineString
		output    LineString
		tolerance float64
	}

	var line = simplifyTest{
		input: LineString{
			Point{153.52, 928.49},
			Point{240.79, 988.95},
			Point{323.34, 1014.40},
			Point{404.41, 1020.08},
			Point{475.60, 981.17},
			Point{497.37, 921.45},
			Point{546.26, 903.57},
			Point{598.10, 907.57},
			Point{655.31, 941.11},
			Point{679.28, 1004.20},
			Point{630.91, 1052.36},
			Point{581.17, 1029.23},
		},
		output: LineString{
			Point{153.52, 928.49},
			Point{404.41, 1020.08},
			Point{546.26, 903.57},
			Point{655.31, 941.11},
			Point{679.28, 1004.20},
			Point{630.91, 1052.36},
			Point{581.17, 1029.23},
		},
		tolerance: 30.0,
	}
	var spiral = simplifyTest{
		input: LineString{
			Point{70.57, 609.01},
			Point{102.21, 618.89},
			Point{125.19, 635.79},
			Point{133.07, 659.34},
			Point{134.86, 688.40},
			Point{121.04, 709.80},
			Point{104.15, 726.70},
			Point{80.45, 731.71},
			Point{56.40, 729.34},
			Point{37.86, 714.81},
			Point{22.83, 692.69},
			Point{23.19, 669.21},
			Point{33.42, 648.38},
			Point{49.74, 635.79},
			Point{84.03, 628.63},
			Point{115.31, 645.31},
			Point{118.75, 681.96},
			Point{109.94, 704.43},
			Point{84.39, 715.17},
			Point{60.20, 716.24},
			Point{42.37, 703.14},
			Point{34.64, 675.16},
			Point{46.31, 658.05},
			Point{69.50, 645.16},
			Point{85.68, 651.96},
			Point{98.78, 669.93},
			Point{92.84, 691.98},
			Point{68.07, 699.21},
			Point{72.58, 676.59},
		},
		output: LineString{
			Point{70.57, 609.01},
			Point{133.07, 659.34},
			Point{104.15, 726.70},
			Point{37.86, 714.81},
			Point{23.19, 669.21},
			Point{84.03, 628.63},
			Point{118.75, 681.96},
			Point{60.20, 716.24},
			Point{46.31, 658.05},
			Point{98.78, 669.93},
			Point{72.58, 676.59},
		},
		tolerance: 30.0,
	}

	for _, test := range []simplifyTest{line, spiral} {
		o := test.input.Simplify(test.tolerance)
		if !reflect.DeepEqual(o, test.output) {
			t.Errorf("%v should equal %v.", o, test.output)
			t.Logf("%v should equal %v.", o, test.output)
		}
	}
}

func TestSimplifyInfiniteLoop(t *testing.T) {
	// This is a self-intersecting shape.
	geometry := Polygon{[]Point{
		Point{X: -871773.1638742175, Y: 497165.8489278648},
		Point{X: -871974.6604566738, Y: 496416.7107209433},
		Point{X: -871878.9516291074, Y: 497176.64415429346},
		Point{X: -871773.1638742175, Y: 497165.8489278648},
		Point{X: -999958.0564939477, Y: 242680.11889008153}}}

	ch := make(chan int)
	go func() {
		geometry.Simplify(100.)
		ch <- 0
	}()

	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		t.Errorf("Simplify %+v timed out.", geometry)
	}
}
