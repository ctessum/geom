package martinez

import (
	"reflect"
	"testing"

	"github.com/Workiva/go-datastructures/queue"
)

func TestQueue(t *testing.T) {
	t.Run("queue should process lest(by x) sweep event first", func(t *testing.T) {
		var queue = queue.NewPriorityQueue(0, false)
		var e1 = &SweepEvent{point: point{x: 0.0, y: 0.0}}
		var e2 = &SweepEvent{point: point{x: 0.5, y: 0.5}}
		queue.Put(e1)
		queue.Put(e2)

		e, err := queue.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(e1, e[0].(*SweepEvent)) {
			t.Error("should equal e1")
		}
		e, err = queue.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(e2, e[0].(*SweepEvent)) {
			t.Error("should equal e2")
		}
	})

	t.Run("queue should process lest(by y) sweep event first", func(t *testing.T) {
		var queue = queue.NewPriorityQueue(0, false)
		var e1 = &SweepEvent{point: point{x: 0.0, y: 0.0}}
		var e2 = &SweepEvent{point: point{x: 0.0, y: 0.5}}

		queue.Put(e1)
		queue.Put(e2)

		e, err := queue.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(e1, e[0].(*SweepEvent)) {
			t.Error("should equal e1")
		}
		e, err = queue.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(e2, e[0].(*SweepEvent)) {
			t.Error("should equal e2")
		}
	})

	t.Run("queue should pop least(by left prop) sweep event first", func(t *testing.T) {
		var queue = queue.NewPriorityQueue(0, false)
		var e1 = &SweepEvent{point: point{x: 0.0, y: 0.0}}
		var e2 = &SweepEvent{point: point{x: 0.0, y: 0.0}}

		queue.Put(e1)
		queue.Put(e2)

		e, err := queue.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(e2, e[0].(*SweepEvent)) {
			t.Error("should equal e2")
		}
		e, err = queue.Get(1)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(e1, e[0].(*SweepEvent)) {
			t.Error("should equal e1")
		}
	})
}

func TestSweepEvent(t *testing.T) {
	t.Run("sweep event comparison x coordinates", func(t *testing.T) {
		var e1 = &SweepEvent{point: point{x: 0.0, y: 0.0}}
		var e2 = &SweepEvent{point: point{x: 0.5, y: 0.5}}

		if e1.Compare(e2) != -1 {
			t.Error("should equal -1")
		}
		if e2.Compare(e1) != 1 {
			t.Error("should equal 1")
		}
	})

	t.Run("sweep event comparison y coordinates", func(t *testing.T) {
		var e1 = &SweepEvent{point: point{x: 0.0, y: 0.0}}
		var e2 = &SweepEvent{point: point{x: 0.0, y: 0.5}}

		if e1.Compare(e2) != -1 {
			t.Error("should equal -1")
		}
		if e2.Compare(e1) != 1 {
			t.Error("should equal 1")
		}
	})

	t.Run("sweep event comparison not left first", func(t *testing.T) {
		var e1 = &SweepEvent{point: point{x: 0.0, y: 0.0}, left: true}
		var e2 = &SweepEvent{point: point{x: 0.0, y: 0.0}, left: false}

		if e1.Compare(e2) != 1 {
			t.Error("should equal 1")
		}
		if e2.Compare(e1) != -1 {
			t.Error("should equal -1")
		}
	})

	t.Run("sweep event comparison shared start point not collinear edges", func(t *testing.T) {
		var e1 = NewSweepEvent(
			point{x: 0.0, y: 0.0},
			true,
			NewSweepEvent(point{x: 1, y: 1}, false, nil, false, normal),
			false,
			normal,
		)
		var e2 = NewSweepEvent(
			point{x: 0.0, y: 0.0},
			true,
			NewSweepEvent(point{x: 2, y: 3}, false, nil, false, normal),
			false,
			normal,
		)

		if e1.Compare(e2) != -1 {
			t.Error("lower is processed first")
		}
		if e2.Compare(e1) != 1 {
			t.Error("higher is processed second")
		}
	})

	t.Run("sweep event comparison collinear edges", func(t *testing.T) {
		var e1 = NewSweepEvent(
			point{x: 0.0, y: 0.0},
			true,
			NewSweepEvent(point{x: 1, y: 1}, false, nil, false, normal),
			true,
			normal,
		)
		var e2 = NewSweepEvent(
			point{x: 0.0, y: 0.0},
			true,
			NewSweepEvent(point{x: 2, y: 2}, false, nil, false, normal),
			false,
			normal,
		)

		if e1.Compare(e2) != -1 {
			t.Error("clipping is processed first")
		}
		if e2.Compare(e1) != 1 {
			t.Error("subject is processed second")
		}
	})
}
