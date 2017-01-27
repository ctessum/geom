package martinez

import "testing"

func TestSweepEvent(t *testing.T) {
	t.Run("isBelow", func(t *testing.T) {
		var s1 = NewSweepEvent(
			point{x: 0, y: 0},
			true,
			NewSweepEvent(point{x: 1, y: 1}, false, nil, false, normal),
			false,
			normal,
		)
		var s2 = NewSweepEvent(
			point{x: 0, y: 1},
			false,
			NewSweepEvent(point{x: 0, y: 0}, false, nil, false, normal),
			false,
			normal,
		)

		if !s1.isBelow(point{x: 0, y: 1}) {
			t.Fail()
		}
		if !s1.isBelow(point{x: 1, y: 2}) {
			t.Fail()
		}
		if s1.isBelow(point{x: 0, y: 0}) {
			t.Fail()
		}
		if s1.isBelow(point{x: 5, y: -1}) {
			t.Fail()
		}

		if s2.isBelow(point{x: 0, y: 1}) {
			t.Fail()
		}
		if s2.isBelow(point{x: 1, y: 2}) {
			t.Fail()
		}
		if s2.isBelow(point{x: 0, y: 0}) {
			t.Fail()
		}
		if s2.isBelow(point{x: 5, y: -1}) {
			t.Fail()
		}
	})

	t.Run("isAbove", func(t *testing.T) {
		var s1 = NewSweepEvent(
			point{x: 0, y: 0},
			true,
			NewSweepEvent(point{x: 1, y: 1}, false, nil, false, normal),
			false,
			normal,
		)
		var s2 = NewSweepEvent(
			point{x: 0, y: 1},
			false,
			NewSweepEvent(point{x: 0, y: 0}, false, nil, false, normal),
			false,
			normal,
		)

		if s1.isAbove(point{x: 0, y: 1}) {
			t.Fail()
		}
		if s1.isAbove(point{x: 1, y: 2}) {
			t.Fail()
		}
		if !s1.isAbove(point{x: 0, y: 0}) {
			t.Fail()
		}
		if !s1.isAbove(point{x: 5, y: -1}) {
			t.Fail()
		}

		if !s2.isAbove(point{x: 0, y: 1}) {
			t.Fail()
		}
		if !s2.isAbove(point{x: 0, y: 1}) {
			t.Fail()
		}
		if !s2.isAbove(point{x: 0, y: 0}) {
			t.Fail()
		}
		if !s2.isAbove(point{x: 5, y: -1}) {
			t.Fail()
		}
	})

	t.Run("isVertical", func(t *testing.T) {
		if !NewSweepEvent(point{x: 0, y: 0}, true, NewSweepEvent(point{x: 0, y: 1}, false, nil, false, normal), false, normal).isVertical() {
			t.Fail()
		}
		if NewSweepEvent(point{x: 0, y: 0}, true, NewSweepEvent(point{x: 0.0001, y: 1}, false, nil, false, normal), false, normal).isVertical() {
			t.Fail()
		}
	})
}
