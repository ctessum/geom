// Copyright 2012 Daniel Connelly.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"math"
	"testing"

	"github.com/ctessum/geom"
)

const EPS = 0.000000001

func TestDist(t *testing.T) {
	p := geom.Point{1, 2}
	q := geom.Point{4, 6}
	dst := math.Sqrt(25)
	if d := dist(p, q); d != dst {
		t.Errorf("dist(%v, %v) = %v; expected %v", p, q, d, dist)
	}
}

func TestNewRect(t *testing.T) {
	p := geom.Point{1.0, -2.5}
	q := geom.Point{3.5, 5.5}
	lengths := geom.Point{2.5, 8.0}

	rect, err := newRect(p, lengths)
	if err != nil {
		t.Errorf("Error on NewRect(%v, %v): %v", p, lengths, err)
	}
	if d := dist(p, rect.Min); d > EPS {
		t.Errorf("Expected p == rect.p")
	}
	if d := dist(q, rect.Max); d > EPS {
		t.Errorf("Expected q == rect.q")
	}
}

func TestNewRectDistError(t *testing.T) {
	p := geom.Point{1.0, -2.5}
	lengths := geom.Point{2.5, -8.0}
	_, err := newRect(p, lengths)
	if _, ok := err.(DistError); !ok {
		t.Errorf("Expected distError on NewRect(%v, %v)", p, lengths)
	}
}

func TestRectSize(t *testing.T) {
	p := geom.Point{1.0, -2.5}
	lengths := geom.Point{2.5, 8.0}
	rect, _ := newRect(p, lengths)
	sze := lengths.X * lengths.Y
	actual := size(rect)
	if sze != actual {
		t.Errorf("Expected %v.size() == %v, got %v", rect, size, actual)
	}
}

func TestRectMargin(t *testing.T) {
	p := geom.Point{1.0, -2.5}
	lengths := geom.Point{2.5, 8.0}
	rect, _ := newRect(p, lengths)
	size := 2*2.5 + 2*8.0
	actual := margin(rect)
	if size != actual {
		t.Errorf("Expected %v.margin() == %v, got %v", rect, size, actual)
	}
}

func TestContainsPoint(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths := geom.Point{6.2, 1.1}
	rect, _ := newRect(p, lengths)

	q := geom.Point{4.5, -1.7}
	if yes := containsPoint(rect, q); !yes {
		t.Errorf("Expected %v contains %v", rect, q)
	}
}

func TestDoesNotContainPoint(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths := geom.Point{6.2, 1.1}
	rect, _ := newRect(p, lengths)

	q := geom.Point{4.5, -0.7}
	if yes := containsPoint(rect, q); yes {
		t.Errorf("Expected %v doesn't contain %v", rect, q)
	}
}

func TestContainsRect(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths1 := geom.Point{6.2, 1.1}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{4.1, -1.9}
	lengths2 := geom.Point{3.2, 0.6}
	rect2, _ := newRect(q, lengths2)

	if yes := containsRect(rect1, rect2); !yes {
		t.Errorf("Expected %v.containsRect(%v", rect1, rect2)
	}
}

func TestDoesNotContainRectOverlaps(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths1 := geom.Point{6.2, 1.1}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{4.1, -1.9}
	lengths2 := geom.Point{3.2, 1.4}
	rect2, _ := newRect(q, lengths2)

	if yes := containsRect(rect1, rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestDoesNotContainRectDisjoint(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths1 := geom.Point{6.2, 1.1}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{1.2, -19.6}
	lengths2 := geom.Point{2.2, 5.9}
	rect2, _ := newRect(q, lengths2)

	if yes := containsRect(rect1, rect2); yes {
		t.Errorf("Expected %v doesn't contain %v", rect1, rect2)
	}
}

func TestNoIntersection(t *testing.T) {
	p := geom.Point{1, 2}
	lengths1 := geom.Point{1, 1}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{-1, -2}
	lengths2 := geom.Point{2.5, 3}
	rect2, _ := newRect(q, lengths2)

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if intersect(rect1, rect2) {
		t.Errorf("Expected intersect(%v, %v) == false", rect1, rect2)
	}
}

func TestNoIntersectionJustTouches(t *testing.T) {
	p := geom.Point{1, 2}
	lengths1 := geom.Point{1, 1}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{-1, -2}
	lengths2 := geom.Point{2.5, 4}
	rect2, _ := newRect(q, lengths2)

	// rect1 and rect2 fail to overlap in just one dimension (second)

	if !intersect(rect1, rect2) {
		t.Errorf("Expected intersect(%v, %v) == true", rect1, rect2)
	}
}

func TestContainmentIntersection(t *testing.T) {
	p := geom.Point{1, 2}
	lengths1 := geom.Point{1, 1}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{1, 2.2}
	lengths2 := geom.Point{0.5, 0.5}
	rect2, _ := newRect(q, lengths2)

	r := geom.Point{1, 2.2}
	s := geom.Point{1.5, 2.7}

	if !intersect(rect1, rect2) {
		t.Errorf("intersect(%v, %v) != %v, %v", rect1, rect2, r, s)
	}
}

func TestOverlapIntersection(t *testing.T) {
	p := geom.Point{1, 2}
	lengths1 := geom.Point{1, 2.5}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{1, 4}
	lengths2 := geom.Point{3, 2}
	rect2, _ := newRect(q, lengths2)

	r := geom.Point{1, 4}
	s := geom.Point{2, 4.5}

	if !intersect(rect1, rect2) {
		t.Errorf("intersect(%v, %v) != %v, %v", rect1, rect2, r, s)
	}
}

func TestToRect(t *testing.T) {
	x := geom.Point{3.7, -2.4}
	tol := 0.05
	rect := ToRect(x, tol)

	p := geom.Point{3.65, -2.45}
	q := geom.Point{3.75, -2.35}
	d1 := dist(p, rect.Min)
	d2 := dist(q, rect.Max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("Expected %v.ToRect(%v) == %v, %v, got %v", x, tol, p, q, rect)
	}
}

func TestBoundingBox(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths1 := geom.Point{1, 15}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{-6.5, 4.7}
	lengths2 := geom.Point{4, 5}
	rect2, _ := newRect(q, lengths2)

	r := geom.Point{-6.5, -2.4}
	s := geom.Point{4.7, 12.6}

	bb := boundingBox(rect1, rect2)
	d1 := dist(r, bb.Min)
	d2 := dist(s, bb.Max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBox(%v, %v) != %v, %v, got %v", rect1, rect2, r, s, bb)
	}
}

func TestBoundingBoxContains(t *testing.T) {
	p := geom.Point{3.7, -2.4}
	lengths1 := geom.Point{1, 15}
	rect1, _ := newRect(p, lengths1)

	q := geom.Point{4.0, 0.0}
	lengths2 := geom.Point{0.56, 6.222222}
	rect2, _ := newRect(q, lengths2)

	bb := boundingBox(rect1, rect2)
	d1 := dist(rect1.Min, bb.Min)
	d2 := dist(rect1.Max, bb.Max)
	if d1 > EPS || d2 > EPS {
		t.Errorf("boundingBox(%v, %v) != %v, got %v", rect1, rect2, rect1, bb)
	}
}

func TestMinDistZero(t *testing.T) {
	p := geom.Point{1, 2}
	r := ToRect(p, 1)
	if d := minDist(p, r); d > EPS {
		t.Errorf("Expected %v.minDist(%v) == 0, got %v", p, r, d)
	}
}

func TestMinDistPositive(t *testing.T) {
	p := geom.Point{1, 2}
	r := geom.Bounds{geom.Point{-1, -4}, geom.Point{2, -2}}
	expected := float64((-2 - 2) * (-2 - 2))
	if d := minDist(p, &r); math.Abs(d-expected) > EPS {
		t.Errorf("Expected %v.minDist(%v) == %v, got %v", p, r, expected, d)
	}
}

func TestMinMaxdist(t *testing.T) {
	p := geom.Point{-3, -2}
	r := geom.Bounds{geom.Point{0, 0}, geom.Point{1, 2}}

	// furthest points from p on the faces closest to p in each dimension
	q1 := geom.Point{0, 2}
	q2 := geom.Point{1, 0}

	// find the closest distance from p to one of these furthest points
	d1 := dist(p, q1)
	d2 := dist(p, q2)
	expected := math.Min(d1*d1, d2*d2)

	if d := minMaxDist(p, &r); math.Abs(d-expected) > EPS {
		t.Errorf("Expected %v.minMaxDist(%v) == %v, got %v", p, r, expected, d)
	}
}
