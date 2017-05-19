package martinez

import (
	"reflect"
	"testing"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/glycerine/rbtree"
)

func TestDivideSegments(t *testing.T) {
	// Two shapes
	s := NewPolygon(NewPath(NewPoint(16, 292), NewPoint(298, 359), NewPoint(153, 203.5), NewPoint(16, 282)))
	c := NewPolygon(NewPath(NewPoint(56, 181), NewPoint(153, 294.5), NewPoint(241.5, 229.5), NewPoint(108.5, 120), NewPoint(56, 181)))

	t.Run("divide 2 segments", func(t *testing.T) {
		var se1 = NewSweepEvent(point{x: 0, y: 0}, true, NewSweepEvent(point{x: 5, y: 5}, false, nil, false, normal), true, normal)
		var se2 = NewSweepEvent(point{x: 0, y: 5}, true, NewSweepEvent(point{x: 5, y: 0}, false, nil, false, normal), false, normal)
		var q = queue.NewPriorityQueue(0, false)

		q.Put(se1)
		q.Put(se2)

		p1, _, _ := segmentIntersection(
			se1.point, se1.otherEvent.point,
			se2.point, se2.otherEvent.point,
			false,
		)

		divideSegment(se1, p1, q)
		divideSegment(se2, p1, q)

		if q.Len() != 6 {
			t.Error("subdivided in 4 segments by intersection point")
		}
	})

	t.Run("possible intersections", func(t *testing.T) {
		var q = queue.NewPriorityQueue(0, false)

		var se1 = NewSweepEvent(s.At(0).At(3), true, NewSweepEvent(s.At(0).At(2), false, nil, false, normal), true, normal)
		var se2 = NewSweepEvent(c.At(0).At(0), true, NewSweepEvent(c.At(0).At(1), false, nil, false, normal), false, normal)

		// console.log(se1.point, se1.left, se1.otherEvent.point, se1.otherEvent.left);
		// console.log(se2.point, se2.left, se2.otherEvent.point, se2.otherEvent.left);

		if v := possibleIntersection(se1, se2, q); v != 1 {
			t.Errorf("possibleIntersection want 1, have %d", v)
		}
		if q.Len() != 4 {
			t.Errorf("queue length want 4, have %d", q.Len())
		}

		eI, err := q.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		var e = eI[0].(*SweepEvent)
		p := point{x: 100.79403384562251, y: 233.41363754101192}
		if !reflect.DeepEqual(e.point, p) {
			t.Errorf("have %+v, want %+v", e.point, p)
		}
		p = point{x: 153, y: 203.5}
		if !reflect.DeepEqual(e.otherEvent.point, p) {
			t.Errorf("have %+v, want %+v", e.point, p)
		}

		eI, err = q.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		e = eI[0].(*SweepEvent)
		if !reflect.DeepEqual(e.point, point{x: 100.79403384562251, y: 233.41363754101192}) {
			t.Fail()
		}
		if !reflect.DeepEqual(e.otherEvent.point, point{x: 56, y: 181}) {
			t.Fail()
		}

		eI, err = q.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		e = eI[0].(*SweepEvent)
		if !reflect.DeepEqual(e.point, point{x: 100.79403384562251, y: 233.41363754101192}) {
			t.Fail()
		}
		if !reflect.DeepEqual(e.otherEvent.point, point{x: 153, y: 294.5}) {
			t.Fail()
		}

		eI, err = q.Get(1)
		if err != nil {
			t.Fatal(err)
		}
		e = eI[0].(*SweepEvent)
		if !reflect.DeepEqual(e.point, point{x: 100.79403384562251, y: 233.41363754101192}) {
			t.Fail()
		}
		if !reflect.DeepEqual(e.otherEvent.point, point{x: 16, y: 282}) {
			t.Fail()
		}
	})

	t.Run("possible intersections on 2 polygons", func(t *testing.T) {

		var bbox = NewBounds()
		var q = fillQueue(s, c, bbox, bbox)
		var p0 = point{x: 16, y: 282}
		var p1 = point{x: 298, y: 359}
		var p2 = point{x: 156, y: 203.5}

		var te = NewSweepEvent(p0, true, nil, true, normal)
		var te2 = NewSweepEvent(p1, false, te, false, normal)
		te.otherEvent = te2

		var te3 = NewSweepEvent(p0, true, nil, true, normal)
		var te4 = NewSweepEvent(p2, true, te3, false, normal)
		te3.otherEvent = te4

		var tr = rbtree.NewTree(compareSegments)

		if !tr.Insert(te) {
			t.Error("insert")
		}
		if !tr.Insert(te3) {
			t.Error("insert")
		}

		if !reflect.DeepEqual(tr.FindGE(te), te) {
			t.Fail()
		}
		if !reflect.DeepEqual(tr.FindGE(te3), te3) {
			t.Fail()
		}

		if compareSegments(te, te3) != 1 {
			t.Fail()
		}
		if compareSegments(te3, te) != -1 {
			t.Fail()
		}

		var segments = subdivideSegments(q, bbox, bbox, 0)
		var leftSegments []*SweepEvent
		for i := 0; i < len(segments); i++ {
			if segments[i].left {
				leftSegments = append(leftSegments, segments[i])
			}
		}

		if len(leftSegments) != 11 {
			t.Fail()
		}

		var E = point{x: 16, y: 282}
		var I = point{x: 100.79403384562252, y: 233.41363754101192}
		var G = point{x: 298, y: 359}
		var C = point{x: 153, y: 294.5}
		var J = point{x: 203.36313843035356, y: 257.5101243166895}
		var F = point{x: 153, y: 203.5}
		var D = point{x: 56, y: 181}
		var A = point{x: 108.5, y: 120}
		var B = point{x: 241.5, y: 229.5}

		type sweepEventTest struct {
			l, r                        point
			inOut, otherInOut, inResult bool
			prevInResult                *sweepEventTest
		}

		intervals := map[string]sweepEventTest{
			"EI": sweepEventTest{l: E, r: I, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
			"IF": sweepEventTest{l: I, r: F, inOut: false, otherInOut: false, inResult: true, prevInResult: nil},
			"FJ": sweepEventTest{l: F, r: J, inOut: false, otherInOut: false, inResult: true, prevInResult: nil},
			"JG": sweepEventTest{l: J, r: G, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
			"EG": sweepEventTest{l: E, r: G, inOut: true, otherInOut: true, inResult: false, prevInResult: nil},
			"DA": sweepEventTest{l: D, r: A, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
			"AB": sweepEventTest{l: A, r: B, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
			"JB": sweepEventTest{l: J, r: B, inOut: true, otherInOut: true, inResult: false, prevInResult: nil},

			"CJ": sweepEventTest{l: C, r: J, inOut: true, otherInOut: false, inResult: true, prevInResult: &sweepEventTest{l: F, r: J}},
			"IC": sweepEventTest{l: I, r: C, inOut: true, otherInOut: false, inResult: true, prevInResult: &sweepEventTest{l: I, r: F}},

			"DI": sweepEventTest{l: D, r: I, inOut: true, otherInOut: true, inResult: false, prevInResult: nil},
		}

		checkContain := func(interval string) {
			data := intervals[interval]
			for i := 0; i < len(leftSegments); i++ {
				seg := leftSegments[i]
				if equals(seg.point, data.l) &&
					equals(seg.otherEvent.point, data.r) &&
					seg.inOut == data.inOut &&
					seg.otherInOut == data.otherInOut &&
					seg.inResult == data.inResult &&
					((seg.prevInResult == nil && data.prevInResult == nil) ||
						(equals(seg.prevInResult.point, data.prevInResult.l) &&
							equals(seg.prevInResult.otherEvent.point, data.prevInResult.r))) {
					t.Logf("pass %s", interval)
					continue
				}
			}
			t.Errorf("fail %s", interval)
		}

		for key := range intervals {
			checkContain(key)
		}
	})
}
