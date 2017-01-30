package martinez

import (
	"reflect"
	"testing"

	"github.com/glycerine/rbtree"
)

func TestSweepLine(t *testing.T) {
	// Two triangles
	s := NewPolygon(NewPath(NewPoint(20, -23.5), NewPoint(170, 74), NewPoint(226, -113.5), NewPoint(20, -23.5)))
	c := NewPolygon(NewPath(NewPoint(54.5, -170.5), NewPoint(140.5, 33.5), NewPoint(239.5, -198), NewPoint(54.4, -170.5)))

	type sweepTest struct {
		*SweepEvent
		name string
	}

	var EF = NewSweepEvent(s.At(0).At(0), true, NewSweepEvent(s.At(0).At(2), false, nil, false, normal), true, normal)
	var EG = NewSweepEvent(s.At(0).At(0), true, NewSweepEvent(s.At(0).At(1), false, nil, false, normal), true, normal)

	var tree = rbtree.NewTree(compareSegments)
	tree.Insert(EF)
	tree.Insert(EG)

	if !reflect.DeepEqual(tree.FindGE(EF).Item().(*SweepEvent), EF) {
		t.Error("not able to retrieve node")
	}
	if !reflect.DeepEqual(tree.Min().Item().(*SweepEvent), EF) {
		t.Error("EF is at the begin")
	}
	if !reflect.DeepEqual(tree.Max().Item().(*SweepEvent), EG) {
		t.Error("EG is at the end")
	}

	var it = tree.FindGE(EF)
	it = it.Next()

	if !reflect.DeepEqual(it.Item().(*SweepEvent), EG) {
		t.Error("EG")
	}

	it = tree.FindLE(EG)
	it = it.Prev()

	if !reflect.DeepEqual(it.Item().(*SweepEvent), EF) {
		t.Error("EF")
	}

	var DA = NewSweepEvent(c.At(0).At(0), true, NewSweepEvent(c.At(0).At(2), false, nil, false, normal), true, normal)
	var DC = NewSweepEvent(c.At(0).At(0), true, NewSweepEvent(c.At(0).At(1), false, nil, false, normal), true, normal)

	tree.Insert(DA)
	tree.Insert(DC)

	var begin = tree.Min()

	if !reflect.DeepEqual(begin.Item().(*SweepEvent), DA) {
		t.Error("DA")
	}
	begin = begin.Next()
	if !reflect.DeepEqual(begin.Item().(*SweepEvent), DC) {
		t.Error("DC")
	}
	begin = begin.Next()
	if !reflect.DeepEqual(begin.Item().(*SweepEvent), EF) {
		t.Error("EF")
	}
	begin = begin.Next()
	if !reflect.DeepEqual(begin.Item().(*SweepEvent), EG) {
		t.Error("EG")
	}
}
