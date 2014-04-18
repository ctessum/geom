package geomop

import (
	"fmt"
	"github.com/ctessum/carto"
	"github.com/twpayne/gogeom/geom"
	"image/color"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

var shape1, shape2 geom.T

func init() {
	shape1 = geom.Polygon{Rings: [][]geom.Point{{geom.Point{0, 0}, geom.Point{1, 0},
		geom.Point{1, 1}, geom.Point{0, 1}, geom.Point{0, 0}}}}
	shape2 = geom.Polygon{Rings: [][]geom.Point{{geom.Point{.25, .25},
		geom.Point{.75, .25}, geom.Point{.75, .75}, geom.Point{.25, .75},
		geom.Point{.25, .25}}}}
}

func TestArea(t *testing.T) {
	a := Area(shape1)
	if different(a, 1.) {
		t.Log(a)
		t.Fail()
	}
	a = Area(shape2)
	if different(a, 0.25) {
		t.Log(a)
		t.Fail()
	}
}

func TestDifference(t *testing.T) {
	FixOrientation(shape1)
	FixOrientation(shape2)
	shape := Construct(shape1, shape2, DIFFERENCE)
	t.Log(shape)
	a := Area(shape)
	if different(a, 0.75) {
		t.Log(a)
		t.Fail()
	}
	drawShapes(shape2, shape1, shape, "difference.png")
}

func TestUnion(t *testing.T) {
	shape := Construct(shape1, shape2, UNION)
	t.Log(shape)
	a := Area(shape)
	if different(a, 1.) {
		t.Log(a)
		t.Fail()
	}
	drawShapes(shape2, shape1, shape, "union.png")
}

func TestIntersection(t *testing.T) {
	shape := Construct(shape1, shape2, INTERSECTION)
	t.Log(shape)
	a := Area(shape)
	if different(a, 0.25) {
		t.Log(a)
		t.Fail()
	}
	drawShapes(shape2, shape1, shape, "intersection.png")
}

func TestXor(t *testing.T) {
	shape := Construct(shape1, shape2, XOR)
	t.Log(shape)
	a := Area(shape)
	if different(a, 0.75) {
		t.Log(a)
		t.Fail()
	}
	drawShapes(shape2, shape1, shape, "xor.png")
}

func TestDifficultShapes1(t *testing.T) {
	difficultShapeA := geom.T(geom.Polygon{[][]geom.Point{{
		{-949.9671190511435, -776530.876383242}, {-971.3149450758938, -776530.876383242},
		{-987.3186852218928, -776530.876383242}, {-971.3143310546875, -776530.9375},
		{-949.9663696289062, -776530.9375}, {-932.567192575676, -776530.876383242}},
		{{-400.82663417423987, -776530.9375}, {-914.386962890625, -776530.9375},
			{-931.7797537201152, -776530.876383242}, {-914.3856589011848, -776530.876383242}},
		{{1847.9894266200718, -776543.8761907278}, {1828.7784833281767, -776543.4981672468},
			{1449.5499414053597, -776539.7180058745}, {1828.78173828125, -776543.5625},
			{1847.9920654296875, -776543.9375}, {1870.3248060378032, -776543.8761907278}},
		{{1892.1079170496669, -776544.3119677962}, {1871.1446207564968, -776543.8840402034},
			{1892.11181640625, -776544.375}, {1906.2987491577942, -776544.5347901696}},
		{{2237.068105247765, -776547.7454548368}, {2225.847412109375, -776547.625},
			{2225.844897193834, -776547.6249235813}}}})

	difficultShapeB := geom.T(geom.Polygon{[][]geom.Point{{
		{-4000, -780000}, {-4000, -776000}, {0, -776000}, {0, -780000}, {-4000, -780000}}}})
	FixOrientation(difficultShapeA)
	FixOrientation(difficultShapeB)

	solution := Construct(difficultShapeA, difficultShapeB, INTERSECTION)
	drawShapes(difficultShapeA, difficultShapeB, solution, "DifficultShapes1.png")
	t.Log(Area(solution))
}

func TestDifficultShapes2(t *testing.T) {
	shape := geom.T(geom.Polygon{[][]geom.Point{
		{{1.825033375e+06, -277681.0625}, {1.824031125e+06, -277704.15625},
			{1.823029e+06, -277727.21875}, {1.82202675e+06, -277750.28125},
			{1.822009375e+06, -276760.3125}, {1.8230115e+06, -276737.25},
			{1.822994e+06, -275747.28125}, {1.821991875e+06, -275770.34375},
			{1.821974375e+06, -274780.4375}, {1.822976625e+06, -274757.34375},
			{1.82397875e+06, -274734.28125}, {1.82396125e+06, -273744.34375},
			{1.822959125e+06, -273767.4375}, {1.822941625e+06, -272777.53125},
			{1.822924125e+06, -271787.625}, {1.822906625e+06, -270797.75},
			{1.823908875e+06, -270774.65625}, {1.823891375e+06, -269784.78125},
			{1.8248935e+06, -269761.6875}, {1.824911e+06, -270751.5625},
			{1.825913125e+06, -270728.46875}, {1.825930625e+06, -271718.34375},
			{1.82693275e+06, -271695.25}, {1.82691525e+06, -270705.34375},
			{1.82689775e+06, -269715.46875}, {1.826880375e+06, -268725.59375},
			{1.825878125e+06, -268748.71875}, {1.82587775e+06, -268727.3125},
			{1.8264545e+06, -268620.5625}, {1.826543375e+06, -268605.375},
			{1.82657925e+06, -268598.875}, {1.827439125e+06, -268449.4375},
			{1.82787675e+06, -268378.90625}, {1.8278825e+06, -268702.46875},
			{1.8279e+06, -269692.34375}, {1.828902125e+06, -269669.21875},
			{1.828919625e+06, -270659.09375}, {1.82992175e+06, -270635.96875},
			{1.82993925e+06, -271625.875}, {1.830941375e+06, -271602.75},
			{1.831943625e+06, -271579.59375}, {1.831961125e+06, -272569.53125},
			{1.83296325e+06, -272546.375}, {1.83298075e+06, -273536.34375},
			{1.83299825e+06, -274526.28125}, {1.83301575e+06, -275516.28125},
			{1.83303325e+06, -276506.25}, {1.833050875e+06, -277496.25},
			{1.834053e+06, -277473.125}, {1.83505525e+06, -277450},
			{1.835037625e+06, -276460}, {1.835020125e+06, -275470},
			{1.83602225e+06, -275446.84375}, {1.83600475e+06, -274456.84375},
			{1.837007e+06, -274433.6875}, {1.838009125e+06, -274410.53125},
			{1.83901125e+06, -274387.375}, {1.8400135e+06, -274364.1875},
			{1.841015625e+06, -274341}, {1.84201775e+06, -274317.8125},
			{1.84302e+06, -274294.625}, {1.8430375e+06, -275284.625},
			{1.842035375e+06, -275307.8125}, {1.842052875e+06, -276297.84375},
			{1.84105075e+06, -276321.03125}, {1.84106825e+06, -277311.0625},
			{1.841085875e+06, -278301.125}, {1.841103375e+06, -279291.1875},
			{1.841121e+06, -280281.28125}, {1.84011875e+06, -280304.4375},
			{1.840136375e+06, -281294.53125}, {1.839134125e+06, -281317.6875},
			{1.839116625e+06, -280327.59375}, {1.839099e+06, -279337.5},
			{1.8390815e+06, -278347.4375}, {1.839063875e+06, -277357.40625},
			{1.83806175e+06, -277380.5625}, {1.83807925e+06, -278370.59375},
			{1.837077125e+06, -278393.75}, {1.836074875e+06, -278416.90625},
			{1.8360925e+06, -279406.9375}, {1.83509025e+06, -279430.0625},
			{1.83510775e+06, -280420.125}, {1.83611e+06, -280397},
			{1.8361275e+06, -281387.09375}, {1.835125375e+06, -281410.21875},
			{1.834123125e+06, -281433.34375}, {1.833121e+06, -281456.4375},
			{1.833103375e+06, -280466.375}, {1.834105625e+06, -280443.25},
			{1.834088125e+06, -279453.1875}, {1.833085875e+06, -279476.3125},
			{1.833068375e+06, -278486.28125}, {1.832066125e+06, -278509.40625},
			{1.83208375e+06, -279499.4375}, {1.8310815e+06, -279522.53125},
			{1.831099e+06, -280512.59375}, {1.8311165e+06, -281502.65625},
			{1.830114375e+06, -281525.75}, {1.830096875e+06, -280535.6875},
			{1.83007925e+06, -279545.65625}, {1.83006175e+06, -278555.625},
			{1.829059625e+06, -278578.71875}, {1.829042125e+06, -277588.71875},
			{1.828039875e+06, -277611.8125}, {1.828057375e+06, -278601.8125},
			{1.828074875e+06, -279591.84375}, {1.82707275e+06, -279614.90625},
			{1.82705525e+06, -278624.90625}, {1.826053e+06, -278648},
			{1.825050875e+06, -278671.0625}, {1.825033375e+06, -277681.0625}},
		{{1.82396125e+06, -273744.34375}, {1.824963375e+06, -273721.25},
			{1.824946e+06, -272731.34375}, {1.8249285e+06, -271741.4375},
			{1.82392625e+06, -271764.53125}, {1.82394375e+06, -272754.4375},
			{1.82396125e+06, -273744.34375}},
		{{1.837042e+06, -276413.6875}, {1.836039875e+06, -276436.84375},
			{1.836057375e+06, -277426.875}, {1.837059625e+06, -277403.71875},
			{1.837042e+06, -276413.6875}},
		{{1.837042e+06, -276413.6875}, {1.83804425e+06, -276390.53125},
			{1.838026625e+06, -275400.53125}, {1.8370245e+06, -275423.6875},
			{1.837042e+06, -276413.6875}},
		{{1.825033375e+06, -277681.0625}, {1.8260355e+06, -277658},
			{1.826018e+06, -276668}, {1.825015875e+06, -276691.09375},
			{1.825033375e+06, -277681.0625}},
		{{1.82993925e+06, -271625.875}, {1.828937125e+06, -271649},
			{1.828954625e+06, -272638.90625}, {1.82995675e+06, -272615.78125},
			{1.82993925e+06, -271625.875}},
		{{1.828004875e+06, -275631.84375}, {1.827987375e+06, -274641.90625},
			{1.827969875e+06, -273651.96875}, {1.82696775e+06, -273675.0625},
			{1.82698525e+06, -274665}, {1.82700275e+06, -275654.96875},
			{1.828004875e+06, -275631.84375}},
		{{1.8260005e+06, -275678.03125}, {1.825983125e+06, -274688.09375},
			{1.824980875e+06, -274711.1875}, {1.824998375e+06, -275701.125},
			{1.8260005e+06, -275678.03125}},
		{{1.832031125e+06, -276529.375}, {1.832013625e+06, -275539.40625},
			{1.831011375e+06, -275562.53125}, {1.83000925e+06, -275585.625},
			{1.83002675e+06, -276575.625}, {1.831029e+06, -276552.5},
			{1.832031125e+06, -276529.375}}}})

	b := shape.Bounds(nil)
	bounds := geom.T(geom.Polygon{[][]geom.Point{{
		{b.Min.X, b.Min.Y}, {b.Max.X, b.Min.Y}, {b.Max.X, b.Max.Y},
		{b.Min.X, b.Max.Y}, {b.Min.X, b.Min.Y}}}})
	FixOrientation(shape)
	FixOrientation(bounds)
	intersection := Construct(shape, bounds, INTERSECTION)
	if different(Area(intersection), Area(shape)) {
		t.Fail()
		t.Log(Area(intersection), Area(shape))
	}
	drawShapes(bounds, shape, intersection, "DifficultShapes2.png")
}

func TestDifficultShapes4(t *testing.T) {
	a := geom.T(geom.Polygon{[][]geom.Point{
		{{-1.05479925e+06, 357453.593}, {-1.05279975e+06, 357459.593}, {-1.05280525e+06, 358450.906}, {-1.0518055e+06, 358453.906}, {-1.051811e+06, 359445.25}, {-1.0518165e+06, 360436.593}, {-1.051822e+06, 361427.968}, {-1.052821875e+06, 361424.937}, {-1.05281625e+06, 360433.562}, {-1.05281075e+06, 359442.218}, {-1.053810625e+06, 359439.218}, {-1.054810375e+06, 359436.218}, {-1.054815875e+06, 360427.562}, {-1.055815625e+06, 360424.562}, {-1.0568155e+06, 360421.562}, {-1.05681e+06, 359430.218}, {-1.055810125e+06, 359433.218}, {-1.055804625e+06, 358441.906}, {-1.054804875e+06, 358444.906}, {-1.05479925e+06, 357453.593}}}})
	b := geom.T(geom.Polygon{[][]geom.Point{
		{{-1.08e+06, 324000}, {-1.08e+06, 360000}, {-1.044e+06, 360000}, {-1.044e+06, 324000}, {-1.08e+06, 324000}}}})

	FixOrientation(a)
	FixOrientation(b)

	intersection := Construct(a, b, INTERSECTION)
	diff := Construct(a, b, DIFFERENCE)
	if different(Area(a), Area(intersection)+Area(diff)) {
		t.Fail()
		t.Log(Area(a), Area(intersection)+Area(diff))
	}
	drawShapes(b, a, intersection, "DifficultShapes4.png")
	drawShapes(b, a, diff, "DifficultShapes4diff.png")
}

var spiral = geom.T(geom.LineString{[]geom.Point{
	{158.69048, 156.42586}, {144.01645, 156.42586}, {139.1901, 161.57183}, {139.1901, 169.9358}, {139.1901, 180.95427}, {150.53931, 194.58874}, {169.42641, 194.58874}, {194.23167, 192.66117}, {210.35714, 175.22916}, {210.35714, 147.61905}, {210.35714, 139.5671}, {202.97619, 92.261905}, {151.64502, 92.261905}, {100.31385, 92.261905}, {84.545455, 139.9026}, {84.545455, 163.72294}, {84.545455, 187.54329}, {106.35281, 238.87446}, {162.38095, 238.87446}, {218.40909, 238.87446}, {248.2684, 188.54978}, {248.2684, 150.63853}, {248.2684, 112.72727}, {216.3961, 58.376623}, {153.65801, 58.376623}, {90.919913, 58.376623}, {54.015152, 113.39827}, {54.015152, 160.36797}, {54.015152, 207.33766}, {92.597403, 267.05628}, {162.38095, 267.05628}, {232.1645, 267.05628}, {274.77273, 201.6342}, {274.77273, 152.31602}, {274.77273, 102.99784}, {233.171, 34.220779}, {154.6645, 34.220779}, {76.158009, 34.220779}, {30.194805, 103.66883}, {30.194805, 159.69697}, {30.194805, 215.72511}, {76.829004, 288.52814}, {163.38745, 288.52814}, {249.94589, 288.52814}, {295.90909, 210.35714}, {295.90909, 151.64502}, {295.90909, 92.9329}, {243.90693, 13.084415}, {155.3355, 13.084415}, {119.40098, 13.084415}, {97.739911, 26.744043}, {97.739911, 26.744043}}})

func TestLine(t *testing.T) {
	shape := geom.T(geom.Polygon{[][]geom.Point{
		{{0, 0}, {0, 200}, {200, 200}, {200, 0}, {0, 0}}}})
	FixOrientation(shape)

	intersection := Construct(spiral, shape, INTERSECTION)
	drawLine(spiral, shape, intersection, "line.png")
}

func TestLines(t *testing.T) {
	line := geom.T(geom.LineString{[]geom.Point{
		{0, 0}, {200, 200}, {100, 200}, {300, 330}, {100, 50}}})

	intersection := Construct(spiral, line, INTERSECTION)
	drawLines(spiral, line, intersection, "lines.png")

}

func drawShapes(a, b, c geom.T, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	carto.DrawShapes(f,
		[]color.NRGBA{{0, 0, 0, 255}, {0, 0, 0, 255}, {0, 0, 0, 255}},
		[]color.NRGBA{{255, 0, 0, 127}, {0, 255, 0, 127},
			{0, 0, 0, 200}},
		1, 0, a, b, c)
	f.Close()
}

func drawLine(a, b, c geom.T, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	carto.DrawShapes(f,
		[]color.NRGBA{{255, 0, 0, 150}, {0, 0, 0, 127}, {0, 0, 255, 100}},
		[]color.NRGBA{{255, 0, 0, 0}, {0, 255, 0, 127},
			{0, 0, 0, 0}},
		4, 0, a, b, c)
	f.Close()
}

func drawLines(a, b, c geom.T, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	carto.DrawShapes(f,
		[]color.NRGBA{{255, 0, 0, 255}, {0, 0, 0, 255}, {0, 0, 255, 255}},
		[]color.NRGBA{{255, 0, 0, 0}, {0, 255, 0, 0},
			{0, 0, 0, 0}},
		4, 8, a, b, c)
	f.Close()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandomPoly(maxWidth, maxHeight float64, vertCnt int) geom.T {
	var result geom.Polygon
	result.Rings = make([][]geom.Point, 1)
	result.Rings[0] = make([]geom.Point, vertCnt)
	for i := 0; i < vertCnt-1; i++ {
		result.Rings[0][i] = geom.Point{rand.Float64() * maxWidth,
			rand.Float64() * maxHeight}
	}
	result.Rings[0][vertCnt-1] = result.Rings[0][0]
	return geom.T(result)
}

// Test operations on random (likely self-intersecting) polygons.
// In output images, intersection is red, xor is green, and union
// is grey (all colors are partially transparent). There should
// be no overlapping green or red areas, and every green or
// red area should be covered in grey. There should also be
// no grey areas that do not overlap with green or red.
func TestRandom(t *testing.T) {
	for i := 0; i < 10; i++ {

		// Generate random subject and clip polygons ...
		numVerticies := 100
		subj := RandomPoly(640, 480, numVerticies)
		clip := RandomPoly(640, 480, numVerticies)
		//subj := Paths([]Path{Path([]*IntPoint{
		//	{76, 351},{460, 441},{136, 71},{76, 351}})})
		//clip:= Paths([]Path{Path([]*IntPoint{
		//	{17, 230},{325, 5},{475, 30},{17, 230}})})

		clipTypes := map[string]Op{"intersection": INTERSECTION,
			"union": UNION, "xor": XOR}
		areas := make(map[string]float64)
		solutions := make(map[string]geom.T)

		for clipType, ct := range clipTypes {
			solutions[clipType] = Construct(subj, clip, ct)
			areas[clipType] = Area(solutions[clipType])
		}
		drawShapes(solutions["intersection"], solutions["xor"],
			solutions["union"], fmt.Sprintf("random_%v.png", i))

		if different(areas["union"], areas["intersection"]+areas["xor"]) {
			t.Logf("%v\t%10.1f%10.1f\tFail", i, areas["union"],
				areas["intersection"]+areas["xor"])
			t.Fail()
		} else {
			t.Logf("%v\t%10.1f%10.1f\tPass", i, areas["union"],
				areas["intersection"]+areas["xor"])
		}
	}
}

func different(a, b float64) bool {
	if math.Abs(a-b)/b > 0.001 {
		return true
	} else {
		return false
	}
}
