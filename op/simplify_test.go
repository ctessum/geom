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

func TestSimplifyInfiniteLoop1(t *testing.T) {
	geometry := geom.Polygon{[]geom.Point{geom.Point{X: -871773.1638742175, Y: 497165.8489278648}, geom.Point{X: -871786.058897913, Y: 497051.36283788737}, geom.Point{X: -871812.6300702088, Y: 496843.16848807875}, geom.Point{X: -871817.2594889546, Y: 496804.2641064115}, geom.Point{X: -871828.8518461098, Y: 496707.8440110069}, geom.Point{X: -871852.4900296698, Y: 496511.2579906229}, geom.Point{X: -871868.0439736051, Y: 496404.69711395353}, geom.Point{X: -871974.6604566738, Y: 496416.7107209433}, geom.Point{X: -872080.2841608367, Y: 496430.28268054314}, geom.Point{X: -872185.3395748375, Y: 496443.8991893595}, geom.Point{X: -872291.2272915924, Y: 496457.2826704979}, geom.Point{X: -872396.018682872, Y: 496470.4231297448}, geom.Point{X: -872499.7279522702, Y: 496482.53964638617}, geom.Point{X: -872578.9629476898, Y: 496491.58048275113}, geom.Point{X: -872587.0595366888, Y: 496492.4496874893}, geom.Point{X: -872593.1615890632, Y: 496493.189034136}, geom.Point{X: -872651.0113261654, Y: 496500.8693724377}, geom.Point{X: -872657.2059697689, Y: 496501.50819449686}, geom.Point{X: -872664.1138388801, Y: 496502.2334463196}, geom.Point{X: -872653.2429897713, Y: 496583.9830456376}, geom.Point{X: -872641.9582658032, Y: 496669.1483742343}, geom.Point{X: -872631.8505936791, Y: 496753.2266395176}, geom.Point{X: -872611.5133181036, Y: 496923.7164973216}, geom.Point{X: -872598.2396354877, Y: 497035.24964090344}, geom.Point{X: -872585.005110706, Y: 497147.1231088713}, geom.Point{X: -872571.3476570707, Y: 497259.16910962015}, geom.Point{X: -872463.8327270573, Y: 497246.7009226391}, geom.Point{X: -872354.5218117528, Y: 497233.79319054075}, geom.Point{X: -872247.2050374286, Y: 497221.01688759495}, geom.Point{X: -872130.8947989381, Y: 497206.8175871102}, geom.Point{X: -872029.1911243061, Y: 497194.16539798304}, geom.Point{X: -871936.3613199592, Y: 497183.25994469505}, geom.Point{X: -871924.8707343122, Y: 497182.53941812366}, geom.Point{X: -871878.9516291074, Y: 497176.64415429346}, geom.Point{X: -871773.1638742175, Y: 497165.8489278648}, geom.Point{X: -999943.8377197147, Y: 242782.39624921232}, geom.Point{X: -999947.5075495535, Y: 242738.22203087248}, geom.Point{X: -999958.0564939477, Y: 242680.11889008153}}}
	_, err := Simplify(geometry, 100.)
	if err != nil {
		t.Fatal(err)
	}
}
