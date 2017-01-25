package martinez

import (
	"fmt"
	"math"

	"github.com/Workiva/go-datastructures/queue"
	"github.com/glycerine/rbtree"
)

type operation int

// These are the available operations.
const (
	intersection operation = iota
	union
	difference
	xor
)

/**
 * @param  {<Array.<Number>} s1
 * @param  {<Array.<Number>} s2
 * @param  {Boolean}         isSubject
 * @param  {Queue}           eventQueue
 * @param  {Array.<Number>}  bbox
 */
func processSegment(s1, s2 Point, isSubject bool, depth int, eventQueue *queue.PriorityQueue, bbox Bounds) {
	// Possible degenerate condition.
	// if (equals(s1, s2)) return;

	var e1 = NewSweepEvent(s1, false, nil, isSubject, normal)
	var e2 = NewSweepEvent(s2, false, e1, isSubject, normal)
	e1.otherEvent = e2

	e1.contourID, e2.contourID = depth, depth

	if e1.Compare(e2) > 0 {
		e2.left = true
	} else {
		e1.left = true
	}

	bbox.Extend(s1)

	// Pushing it so the queue is sorted from left to right,
	// with object on the left having the highest priority.
	eventQueue.Put(e1)
	eventQueue.Put(e2)
}

func processPolygon(polygon Polygon, isSubject bool, queue *queue.PriorityQueue, bbox Bounds) {
	for i := 0; i < polygon.Len(); i++ {
		path := polygon.At(i)
		for j := 0; j < path.Len(); j++ {
			processSegment(path.At(i), path.At(i+1), isSubject, i+1, queue, bbox)
		}
	}
}

func fillQueue(subject, clipping Polygon, sbbox, cbbox Bounds) *queue.PriorityQueue {
	eventQueue := queue.NewPriorityQueue(0, false) // (null, compareEvents);

	processPolygon(subject, true, eventQueue, sbbox)
	processPolygon(clipping, false, eventQueue, cbbox)

	return eventQueue
}

func computeFields(event, prev *SweepEvent, sweepLine *rbtree.Tree, op operation) {
	// compute inOut and otherInOut fields
	if prev == nil {
		event.inOut = false
		event.otherInOut = true

		// previous line segment in sweepline belongs to the same polygon
	} else if event.isSubject == prev.isSubject {
		event.inOut = !prev.inOut
		event.otherInOut = prev.otherInOut

		// previous line segment in sweepline belongs to the clipping polygon
	} else {
		event.inOut = !prev.otherInOut
		if prev.isVertical() {
			event.otherInOut = !prev.inOut
		}
		event.otherInOut = prev.inOut
	}

	// compute prevInResult field
	if prev != nil {
		if !inResult(prev, op) || prev.isVertical() {
			event.prevInResult = prev.prevInResult
		}
		event.prevInResult = prev // TODO: not sure what this means.
	}
	// check if the line segment belongs to the Boolean operation
	event.inResult = inResult(event, op)
}

func inResult(event *SweepEvent, op operation) bool {
	switch event.edgeType {
	case normal:
		switch op {
		case intersection:
			return !event.otherInOut
		case union:
			return event.otherInOut
		case difference:
			return (event.isSubject && event.otherInOut) ||
				(!event.isSubject && !event.otherInOut)
		case xor:
			return true
		}
	case sameTransition:
		return op == intersection || op == union
	case differentTransition:
		return op == difference
	case nonContributing:
		return false
	}
	return false
}

/**
 * @param  {SweepEvent} se1
 * @param  {SweepEvent} se2
 * @param  {Queue}      queue
 * @return {Number}
 */
func possibleIntersection(se1, se2 *SweepEvent, queue *queue.PriorityQueue) int {
	// that disallows self-intersecting polygons,
	// did cost us half a day, so I'll leave it
	// out of respect
	// if (se1.isSubject === se2.isSubject) return;

	p1, _, nintersections := segmentIntersection(se1.point, se1.otherEvent.point, se2.point, se2.otherEvent.point, false) // TODO: I think this should be false but not sure.

	if nintersections == 0 {
		return 0
	} // no intersection

	// the line segments intersect at an endpoint of both line segments
	if (nintersections == 1) &&
		(equals(se1.point, se2.point) ||
			equals(se1.otherEvent.point, se2.otherEvent.point)) {
		return 0
	}

	if nintersections == 2 && se1.isSubject == se2.isSubject {
		if se1.contourID == se2.contourID {
			panic(fmt.Errorf("Edges of the same polygon overlap, %+v, %+v, %+v, %+v",
				se1.point, se1.otherEvent.point, se2.point, se2.otherEvent.point))
		}
		return 0
	}

	// The line segments associated to se1 and se2 intersect
	if nintersections == 1 {

		// if the intersection point is not an endpoint of se1
		if !equals(se1.point, p1) && !equals(se1.otherEvent.point, p1) {
			divideSegment(se1, p1, queue)
		}

		// if the intersection point is not an endpoint of se2
		if !equals(se2.point, p1) && !equals(se2.otherEvent.point, p1) {
			divideSegment(se2, p1, queue)
		}
		return 1
	}

	// The line segments associated to se1 and se2 overlap
	var events []*SweepEvent
	var leftCoincide = false
	var rightCoincide = false

	if equals(se1.point, se2.point) {
		leftCoincide = true // linked
	} else if se1.Compare(se2) == 1 {
		events = []*SweepEvent{se2, se1}
	} else {
		events = []*SweepEvent{se1, se2}
	}

	if equals(se1.otherEvent.point, se2.otherEvent.point) {
		rightCoincide = true
	} else if se1.otherEvent.Compare(se2.otherEvent) == 1 {
		events = append(events, se2.otherEvent, se1.otherEvent)
	} else {
		events = append(events, se1.otherEvent, se2.otherEvent)
	}

	if (leftCoincide && rightCoincide) || leftCoincide {
		// both line segments are equal or share the left endpoint
		se1.edgeType = nonContributing
		if se1.inOut == se2.inOut {
			se2.edgeType = sameTransition
		}
		se2.edgeType = differentTransition

		if leftCoincide && !rightCoincide {
			// honestly no idea, but changing events selection from [2, 1]
			// to [0, 1] fixes the overlapping self-intersecting polygons issue
			divideSegment(events[0].otherEvent, events[1].point, queue)
		}
		return 2
	}

	// the line segments share the right endpoint
	if rightCoincide {
		divideSegment(events[0], events[1].point, queue)
		return 3
	}

	// no line segment includes totally the other one
	if events[0] != events[3].otherEvent {
		divideSegment(events[0], events[1].point, queue)
		divideSegment(events[1], events[2].point, queue)
		return 3
	}

	// one line segment includes the other one
	divideSegment(events[0], events[1].point, queue)
	divideSegment(events[3].otherEvent, events[2].point, queue)

	return 3
}

/**
 * @param  {SweepEvent} se
 * @param  {Array.<Number>} p
 * @param  {Queue} queue
 * @return {Queue}
 */
func divideSegment(se *SweepEvent, p Point, queue *queue.PriorityQueue) *queue.PriorityQueue {
	var r = NewSweepEvent(p, false, se, se.isSubject, normal)
	var l = NewSweepEvent(p, true, se.otherEvent, se.isSubject, normal)

	if equals(se.point, se.otherEvent.point) {
		panic(fmt.Errorf("what is that? %+v", se))
	}

	r.contourID, l.contourID = se.contourID, se.contourID

	// avoid a rounding error. The left event would be processed after the right event
	if l.Compare(se.otherEvent) > 0 {
		se.otherEvent.left = true
		l.left = false
	}

	// avoid a rounding error. The left event would be processed after the right event
	// if (compareEvents(se, r) > 0) {}

	se.otherEvent.otherEvent = l
	se.otherEvent = r

	queue.Put(l)
	queue.Put(r)

	return queue
}

func subdivideSegments(eventQueue *queue.PriorityQueue, subject, clipping Polygon, sbbox, cbbox Bounds, op operation) []*SweepEvent {
	sweepLine := rbtree.NewTree(compareSegments)
	var sortedEvents []*SweepEvent

	var rightbound = math.Min(sbbox.Max().X(), cbbox.Max().X())

	var prev, next rbtree.Iterator

	for eventQueue.Len() > 0 {
		var eventI, err = eventQueue.Get(1)
		if err != nil {
			panic(err)
		}
		event := eventI[0].(*SweepEvent)
		sortedEvents = append(sortedEvents, event)

		// optimization by bboxes for intersection and difference goes here
		if (op == intersection && event.point.X() > rightbound) ||
			(op == difference && event.point.X() > sbbox.Max().X()) {
			break
		}

		if event.left {
			sweepLine.Insert(event)
			// _renderSweepLine(sweepLine, event.point, event);

			next = sweepLine.FindGE(event)
			prev = sweepLine.FindLE(event)
			event.iterator = sweepLine.FindGE(event)

			// Cannot get out of the tree what we just put there
			/*if prev == nil || next == nil {
				fmt.Println("brute")
				var iterators = findIterBrute(sweepLine)
				prev = iterators[0]
				next = iterators[1]
			}*/

			if !prev.Min() {
				prev = prev.Prev()
			} else {
				prev = sweepLine.Min().Next()
			}
			next = next.Next()

			computeFields(event, prev.Item().(*SweepEvent), sweepLine, op)

			if next.Item() != nil {
				if possibleIntersection(event, next.Item().(*SweepEvent), eventQueue) == 2 {
					computeFields(event, prev.Item().(*SweepEvent), sweepLine, op)
					computeFields(event, next.Item().(*SweepEvent), sweepLine, op)
				}
			}

			if prev.Item() != nil {
				if possibleIntersection(prev.Item().(*SweepEvent), event, eventQueue) == 2 {
					var prevprev = sweepLine.FindLE(prev)
					if !prevprev.Min() {
						prevprev = prevprev.Prev()
					} else {
						prevprev = sweepLine.Max()
						prevprev = prevprev.Next()
					}
					computeFields(prev.Item().(*SweepEvent), prevprev.Item().(*SweepEvent), sweepLine, op)
					computeFields(event, prev.Item().(*SweepEvent), sweepLine, op)
				}
			}
		} else {
			event = event.otherEvent
			next = sweepLine.FindGE(event)
			prev = sweepLine.FindLE(event)

			// _renderSweepLine(sweepLine, event.otherEvent.point, event);

			if prev.Limit() || next.Limit() {
				continue
			}

			if !prev.Min() {
				prev = prev.Prev()
			} else {
				prev = sweepLine.Min().Next()
			}
			next = next.Next()
			sweepLine.DeleteWithKey(event)

			//_renderSweepLine(sweepLine, event.otherEvent.point, event);

			if next.Item() != nil && prev.Item() != nil {
				possibleIntersection(prev.Item().(*SweepEvent), next.Item().(*SweepEvent), eventQueue)
			}
		}
	}
	return sortedEvents
}

/*func findIterBrute(sweepLine *rbtree.Tree) { //, q)
	var prev = sweepLine.iterator()
	var next = sweepLine.iterator()
	var it = sweepLine.iterator() //, data;
	for {
		data := it.next()
		if data == nil {
			break
		}
		prev.next()
		next.next()
		if data == event {
			break
		}
	}
	return prev, next
}*/

func changeOrientation(a path) path {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
	return a
}

/**
 * @param  {Array.<SweepEvent>} sortedEvents
 * @return {Array.<SweepEvent>}
 */
func orderEvents(sortedEvents []*SweepEvent) []*SweepEvent {
	var resultEvents []*SweepEvent
	for i := 0; i < len(sortedEvents); i++ {
		event := sortedEvents[i]
		if (event.left && event.inResult) || (!event.left && event.otherEvent.inResult) {
			resultEvents = append(resultEvents, event)
		}
	}

	// Due to overlapping edges the resultEvents array can be not wholly sorted
	sorted := false
	for !sorted {
		sorted = true
		for i := 0; i < len(resultEvents); i++ {
			if (i+1) < len(resultEvents) &&
				resultEvents[i].Compare(resultEvents[i+1]) == 1 {
				resultEvents[i], resultEvents[i+1] = resultEvents[i+1], resultEvents[i]
				sorted = false
			}
		}
	}

	for i := 0; i < len(resultEvents); i++ {
		resultEvents[i].pos = i
	}

	for i := 0; i < len(resultEvents); i++ {
		if !resultEvents[i].left {
			var temp = resultEvents[i].pos
			resultEvents[i].pos = resultEvents[i].otherEvent.pos
			resultEvents[i].otherEvent.pos = temp
		}
	}

	return resultEvents
}

/**
 * @param  {Array.<SweepEvent>} sortedEvents
 * @return {Array.<*>} polygons
 */
func connectEdges(sortedEvents []*SweepEvent) multiPath {
	var resultEvents = orderEvents(sortedEvents)

	// "false"-filled array
	processed := make([]bool, len(resultEvents))
	var result multiPath

	var depth []int
	var holeOf []int
	var isHole []bool

	for i := 0; i < len(resultEvents); i++ {
		if processed[i] {
			continue
		}

		var contour path
		result = append(result, contour)

		var ringID = len(result) - 1
		depth = append(depth, 0)
		holeOf = append(holeOf, -1)

		if resultEvents[i].prevInResult != nil {
			var lowerContourID = resultEvents[i].prevInResult.contourID
			if !resultEvents[i].prevInResult.resultInOut {
				//      addHole(result[lowerContourId], ringId);
				holeOf[ringID] = lowerContourID
				depth[ringID] = depth[lowerContourID] + 1
				isHole[ringID] = true
			} else if isHole[lowerContourID] {
				//    addHole(result[holeOf[lowerContourId]], ringId);
				holeOf[ringID] = holeOf[lowerContourID]
				depth[ringID] = depth[lowerContourID]
				isHole[ringID] = true
			}
		}

		var pos = i
		var initial = resultEvents[i].point
		contour = append(contour, initial)

		for pos >= i {
			processed[pos] = true

			if resultEvents[pos].left {
				resultEvents[pos].resultInOut = false
				resultEvents[pos].contourID = ringID
			} else {
				resultEvents[pos].otherEvent.resultInOut = true
				resultEvents[pos].otherEvent.contourID = ringID
			}

			pos = resultEvents[pos].pos
			processed[pos] = true

			contour = append(contour, resultEvents[pos].point)
			pos = nextPos(pos, resultEvents, processed)
		}

		if pos == -1 {
			pos = i
		}

		processed[pos] = true
		processed[resultEvents[pos].pos] = true
		resultEvents[pos].otherEvent.resultInOut = true
		resultEvents[pos].otherEvent.contourID = ringID

		// depth is even
		/* eslint-disable no-bitwise */
		if depth[ringID]%2 == 0 {
			changeOrientation(contour)
		}
		/* eslint-enable no-bitwise */
	}
	return result
}

/**
 * @param  {Number} pos
 * @param  {Array.<SweepEvent>} resultEvents
 * @param  {Array.<Boolean>}    processed
 * @return {Number}
 */
func nextPos(pos int, resultEvents []*SweepEvent, processed []bool) int {
	var newPos = pos + 1
	var length = len(resultEvents)
	for newPos < length &&
		equals(resultEvents[newPos].point, resultEvents[pos].point) {
		if !processed[newPos] {
			return newPos
		}
		newPos = newPos + 1
	}

	newPos = pos - 1

	for processed[newPos] {
		newPos = newPos - 1
	}
	return newPos
}

// trivialOperation checks for a trivial result in the case where
// one of the polygons is empty.
func trivialOperation(subject, clipping Polygon, op operation) Polygon {
	if subject.Len()*clipping.Len() != 0 {
		return nil
	}
	if op == intersection {
		return nil
	} else if op == difference {
		return subject
	} else if op == union || op == xor {
		if subject.Len() == 0 {
			return clipping
		}
		return subject
	}
	panic("invalid trivial operation")
}

// compareBBoxes checks for a trivial solution where the bounding boxes of
// the polygons don't overlap, and returns the trivial solution if it exists.
// Otherwise, it returns nil.
func compareBBoxes(subject, clipping Polygon, sbbox, cbbox Bounds, op operation) Polygon {
	if sbbox.Min().X() > cbbox.Max().X() ||
		cbbox.Min().X() > sbbox.Max().X() ||
		sbbox.Min().Y() > cbbox.Max().Y() ||
		cbbox.Min().Y() > sbbox.Max().Y() {
		if op == intersection {
			// The result here is nil if there is no overlap.
			return nil
		} else if op == difference {
			// The result here is the subject if there is no overlap.
			return subject
		} else if op == union || op == xor {
			// The result here is the combination of the subject and
			// clipping polygons if there is no overlap.
			result := make(multiPath, subject.Len()+clipping.Len())
			for i := 0; i < subject.Len(); i++ {
				result[i] = subject.At(i)
			}
			for i := 0; i < clipping.Len(); i++ {
				result[i+subject.Len()] = subject.At(i)
			}
			return result
		}
	}
	// There is no trivial solution
	return nil
}

func boolean(subject, clipping Polygon, op operation) Polygon {
	var trivial = trivialOperation(subject, clipping, op)
	if trivial != nil {
		return trivial
	}
	var sbbox = NewBounds()
	var cbbox = NewBounds()

	var eventQueue = fillQueue(subject, clipping, sbbox, cbbox)

	trivial = compareBBoxes(subject, clipping, sbbox, cbbox, op)
	if trivial != nil {
		return trivial
	}
	var sortedEvents = subdivideSegments(eventQueue, subject, clipping, sbbox, cbbox, op)
	return connectEdges(sortedEvents)
}
