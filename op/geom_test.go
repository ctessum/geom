// Copyright (c) 2011 Mateusz Czapliński
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package op

import (
	"fmt"
	"github.com/ctessum/geom"
	"math"
	. "testing"
)

func verify(t *T, cond bool, format string, args ...interface{}) {
	if !cond {
		t.Errorf(format, args...)
	}
}

func circa(f, g float64) bool {
	//TODO: (f-g)/g < 1e-6  ?
	return math.Abs(f-g) < 1e-6
}

func TestPoint(t *T) {
	verify(t, PointEquals(geom.Point{0, 0}, geom.Point{0, 0}),
		"Expected equal points")
	verify(t, PointEquals(geom.Point{1, 2}, geom.Point{1, 2}),
		"Expected equal points")
	verify(t, PointEquals(geom.Point{1, 2}, geom.Point{1.000000000001, 2}),
		"Expected equal points")
	verify(t, circa(lengthToOrigin(geom.Point{3, 4}), 5), "Expected length 5")
}

func TestContourAdd(t *T) {
	c := contour{}
	pp := []geom.Point{{1, 2}, {3, 4}, {5, 6}}
	for i := range pp {
		c = append(c, pp[i])
	}
	verify(t, len(c) == len(pp), "Expected all points in contour")
	for i := range pp {
		verify(t, PointEquals(c[i], pp[i]), "Wrong point at position %d", i)
	}
}

func TestContourBoundingBox(t *T) {
	// TODO
}

func TestContourSegment(t *T) {
	c := contour([]geom.Point{{1, 2}, {3, 4}, {5, 6}})
	segeq := func(s1, s2 segment) bool {
		return PointEquals(s1.start, s2.start) && PointEquals(s1.end, s2.end)
	}
	verify(t, segeq(c.segment(0), segment{geom.Point{1, 2}, geom.Point{3, 4}}),
		"Expected segment 0")
	verify(t, segeq(c.segment(1), segment{geom.Point{3, 4}, geom.Point{5, 6}}),
		"Expected segment 1")
	verify(t, segeq(c.segment(2), segment{geom.Point{5, 6}, geom.Point{1, 2}}),
		"Expected segment 2")
}

func TestContourSegmentError1(t *T) {
	c := contour([]geom.Point{{1, 2}, {3, 4}, {5, 6}})

	defer func() {
		verify(t, recover() != nil, "Expected error")
	}()
	_ = c.segment(3)
}

type pointresult struct {
	p      geom.Point
	result int
}

func TestContourContains(t *T) {
	var cases1 []pointresult
	c1 := contour([]geom.Point{{0, 0}, {10, 0}, {0, 10}})
	c2 := contour([]geom.Point{{0, 0}, {0, 10}, {10, 0}}) // opposite rotation
	cases1 = []pointresult{
		{geom.Point{1, 1}, 1},
		{geom.Point{2, .1}, 1},
		{geom.Point{10, 10}, 0},
		{geom.Point{11, 0}, 0},
		{geom.Point{0, 11}, 0},
		{geom.Point{-1, -1}, 0},
	}
	for i, v := range cases1 {
		verify(t, pointInPoly(v.p, c1) == v.result, "Expected %v for point %d for c1", v.result, i)
		verify(t, pointInPoly(v.p, c2) == v.result, "Expected %v for point %d for c2", v.result, i)
	}
}

func ExamplePolygon_Construct() {
	subject := geom.Geom(geom.Polygon{
		{{1, 1}, {1, 2}, {2, 2}, {2, 1}}}) // small square
	clipping := geom.Geom(geom.Polygon{
		{{0, 0}, {0, 3}, {3, 0}}}) // overlapping triangle

	// Calculate the intersection
	result, err := Construct(subject, clipping, INTERSECTION)
	handle(err)

	out := []string{}
	for _, point := range result.(geom.Polygon)[0] {
		out = append(out, fmt.Sprintf("%v", point))
	}
	fmt.Println(out)
	// Output: [{1 2} {1 1} {2 1} {1 2}]
}
