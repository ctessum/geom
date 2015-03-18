package op

import (
	"reflect"
	"testing"

	"github.com/ctessum/geom"
)

func TestSimplify(t *testing.T) {
	type simplifyTest struct {
		input     geom.LineString
		output    geom.LineString
		tolerance float64
	}

	var line = simplifyTest{
		input: geom.LineString{
			geom.Point{153.52, 928.49},
			geom.Point{240.79, 988.95},
			geom.Point{323.34, 1014.40},
			geom.Point{404.41, 1020.08},
			geom.Point{475.60, 981.17},
			geom.Point{497.37, 921.45},
			geom.Point{546.26, 903.57},
			geom.Point{598.10, 907.57},
			geom.Point{655.31, 941.11},
			geom.Point{679.28, 1004.20},
			geom.Point{630.91, 1052.36},
			geom.Point{581.17, 1029.23},
		},
		output: geom.LineString{
			geom.Point{153.52, 928.49},
			geom.Point{404.41, 1020.08},
			geom.Point{546.26, 903.57},
			geom.Point{655.31, 941.11},
			geom.Point{679.28, 1004.20},
			geom.Point{630.91, 1052.36},
			geom.Point{581.17, 1029.23},
		},
		tolerance: 30.0,
	}
	var spiral = simplifyTest{
		input: geom.LineString{
			geom.Point{70.57, 609.01},
			geom.Point{102.21, 618.89},
			geom.Point{125.19, 635.79},
			geom.Point{133.07, 659.34},
			geom.Point{134.86, 688.40},
			geom.Point{121.04, 709.80},
			geom.Point{104.15, 726.70},
			geom.Point{80.45, 731.71},
			geom.Point{56.40, 729.34},
			geom.Point{37.86, 714.81},
			geom.Point{22.83, 692.69},
			geom.Point{23.19, 669.21},
			geom.Point{33.42, 648.38},
			geom.Point{49.74, 635.79},
			geom.Point{84.03, 628.63},
			geom.Point{115.31, 645.31},
			geom.Point{118.75, 681.96},
			geom.Point{109.94, 704.43},
			geom.Point{84.39, 715.17},
			geom.Point{60.20, 716.24},
			geom.Point{42.37, 703.14},
			geom.Point{34.64, 675.16},
			geom.Point{46.31, 658.05},
			geom.Point{69.50, 645.16},
			geom.Point{85.68, 651.96},
			geom.Point{98.78, 669.93},
			geom.Point{92.84, 691.98},
			geom.Point{68.07, 699.21},
			geom.Point{72.58, 676.59},
		},
		output: geom.LineString{
			geom.Point{70.57, 609.01},
			geom.Point{133.07, 659.34},
			geom.Point{104.15, 726.70},
			geom.Point{37.86, 714.81},
			geom.Point{23.19, 669.21},
			geom.Point{84.03, 628.63},
			geom.Point{118.75, 681.96},
			geom.Point{60.20, 716.24},
			geom.Point{46.31, 658.05},
			geom.Point{98.78, 669.93},
			geom.Point{72.58, 676.59},
		},
		tolerance: 30.0,
	}

	for _, test := range []simplifyTest{line, spiral} {
		o, err := Simplify(test.input, test.tolerance)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(o, test.output) {
			t.Errorf("%v should equal %v.", o, test.output)
			t.Logf("%v should equal %v.", o, test.output)
		}
	}
}
