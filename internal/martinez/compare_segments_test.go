package martinez

import (
	"reflect"
	"testing"

	"github.com/glycerine/rbtree"
)

func TestCompareSegments(t *testing.T) {
	t.Run("not collinear", func(t *testing.T) {

		t.Run("shared left point - right point first", func(t *testing.T) {
			var tree = rbtree.NewTree(compareSegments)
			var pt = point{x: 0.0, y: 0.0}
			var se1 = NewSweepEvent(
				pt,
				true,
				NewSweepEvent(point{x: 1, y: 1}, false, nil, false, normal),
				false,
				normal,
			)
			var se2 = NewSweepEvent(
				pt,
				true,
				NewSweepEvent(point{x: 2, y: 3}, false, nil, false, normal),
				false,
				normal,
			)

			tree.Insert(se1)
			tree.Insert(se2)

			if !reflect.DeepEqual(tree.Max().Item().(*SweepEvent).otherEvent.point, point{x: 2, y: 3}) {
				t.Fail()
			}
			if !reflect.DeepEqual(tree.Min().Item().(*SweepEvent).otherEvent.point, point{x: 1, y: 1}) {
				t.Fail()
			}
		})

		t.Run("different left point - right point y coord to sort", func(t *testing.T) {
			var tree = rbtree.NewTree(compareSegments)
			var se1 = NewSweepEvent(
				point{x: 0.0, y: 1},
				true,
				NewSweepEvent(point{x: 1, y: 1}, false, nil, false, normal),
				false,
				normal,
			)
			var se2 = NewSweepEvent(
				point{x: 0.0, y: 2},
				true,
				NewSweepEvent(point{x: 2, y: 3}, false, nil, false, normal),
				false,
				normal,
			)

			tree.Insert(se1)
			tree.Insert(se2)

			if !reflect.DeepEqual(tree.Min().Item().(*SweepEvent).otherEvent.point, point{x: 1, y: 1}) {
				t.Fail()
			}
			if !reflect.DeepEqual(tree.Max().Item().(*SweepEvent).otherEvent.point, point{x: 2, y: 3}) {
				t.Fail()
			}
		})

		t.Run("events order in sweep line", func(t *testing.T) {
			var se1 = NewSweepEvent(
				point{x: 0, y: 1},
				true,
				NewSweepEvent(point{x: 2, y: 1}, false, nil, false, normal),
				false,
				normal,
			)
			var se2 = NewSweepEvent(
				point{x: -1.0, y: 0},
				true,
				NewSweepEvent(point{x: 2, y: 3}, false, nil, false, normal),
				false,
				normal,
			)
			var se3 = NewSweepEvent(
				point{x: 0.0, y: 1},
				true,
				NewSweepEvent(point{x: 3, y: 4}, false, nil, false, normal),
				false,
				normal,
			)
			var se4 = NewSweepEvent(
				point{x: -1, y: 0},
				true,
				NewSweepEvent(point{x: 3, y: 1}, false, nil, false, normal),
				false,
				normal,
			)

			if se1.Compare(se2) != 1 {
				t.Fail()
			}
			if se2.isBelow(se1.point) != false {
				t.Fail()
			}
			if se2.isAbove(se1.point) != true {
				t.Fail()
			}

			if compareSegments(se1, se2) != -1 {
				t.Error("compare segments")
			}
			if compareSegments(se2, se1) != 1 {
				t.Error("compare segments inverted")
			}

			if se3.Compare(se4) != 1 {
				t.Fail()
			}
			if se4.isAbove(se3.point) != false {
				t.Fail()
			}
		})

		t.Run("first point is below", func(t *testing.T) {
			var se2 = NewSweepEvent(
				point{x: 0.0, y: 1},
				true,
				NewSweepEvent(point{x: 2, y: 1}, false, nil, false, normal),
				false,
				normal,
			)
			var se1 = NewSweepEvent(
				point{x: -1.0, y: 0},
				true,
				NewSweepEvent(point{x: 2, y: 3}, false, nil, false, normal),
				false,
				normal,
			)

			if se1.isBelow(se2.point) != false {
				t.Fail()
			}
			if compareSegments(se1, se2) != 1 {
				t.Error("compare segments")
			}
		})
	})

	t.Run("collinear segments", func(t *testing.T) {
		var se1 = NewSweepEvent(
			point{x: 1, y: 1},
			true,
			NewSweepEvent(point{x: 5, y: 1}, false, nil, false, normal),
			true,
			normal,
		)
		var se2 = NewSweepEvent(
			point{x: 2, y: 1},
			true,
			NewSweepEvent(point{x: 3, y: 1}, false, nil, false, normal),
			false,
			normal,
		)

		if se1.isSubject == se2.isSubject {
			t.Fail()
		}
		if compareSegments(se1, se2) != -1 {
			t.Fail()
		}
	})

	t.Run("collinear shared left point", func(t *testing.T) {
		var pt = point{x: 0, y: 1}
		var se1 = NewSweepEvent(
			pt,
			true,
			NewSweepEvent(point{x: 5, y: 1}, false, nil, false, normal),
			false,
			normal,
		)
		var se2 = NewSweepEvent(
			pt,
			true,
			NewSweepEvent(point{x: 3, y: 1}, false, nil, false, normal),
			false,
			normal,
		)

		se1.contourID = 1
		se2.contourID = 2

		if se1.isSubject != se2.isSubject {
			t.Fail()
		}
		if se1.point != se2.point {
			t.Fail()
		}

		if compareSegments(se1, se2) != -1 {
			t.Fail()
		}

		se1.contourID = 2
		se2.contourID = 1

		if compareSegments(se1, se2) != 1 {
			t.Fail()
		}
	})

	t.Run("collinear same polygon different left points", func(t *testing.T) {
		var se1 = NewSweepEvent(
			point{x: 1, y: 1},
			true,
			NewSweepEvent(point{x: 5, y: 1}, false, nil, false, normal),
			true,
			normal,
		)
		var se2 = NewSweepEvent(
			point{x: 2, y: 1},
			true,
			NewSweepEvent(point{x: 3, y: 1}, false, nil, false, normal),
			true,
			normal,
		)

		if se1.isSubject != se2.isSubject {
			t.Fail()
		}
		if se1.point == se2.point {
			t.Fail()
		}
		if compareSegments(se1, se2) != -1 {
			t.Fail()
		}
		if compareSegments(se2, se1) != 1 {
			t.Fail()
		}
	})
}
