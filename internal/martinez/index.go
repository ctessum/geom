package martinez

import (
  "github.com/Workiva/go-datastructures/queue"
  "github.com/glycerine/rbtree"
)

type operation int
const (
  INTERSECTION  operation = iota
 UNION
 DIFFERENCE
 XOR
)

type EMPTY  struct{}

//var Queue           = require('tinyqueue');
//var Tree            = require('bintrees').RBTree;


/**
 * @param  {<Array.<Number>} s1
 * @param  {<Array.<Number>} s2
 * @param  {Boolean}         isSubject
 * @param  {Queue}           eventQueue
 * @param  {Array.<Number>}  bbox
 */
func processSegment(s1, s2 Segment, isSubject bool, depth int, eventQueue *queue.PriorityQueue, bbox Bounds) {
  // Possible degenerate condition.
  // if (equals(s1, s2)) return;

  var e1 = NewSweepEvent(s1, false, undefined, isSubject);
  var e2 = NewSweepEvent(s2, false, e1,        isSubject);
  e1.otherEvent = e2;

  e1.contourId, e2.contourId = depth, depth;

  if (compareEvents(e1, e2) > 0) {
    e2.left = true;
  } else {
    e1.left = true;
  }

  bbox.Extend(s1.Start())

  // Pushing it so the queue is sorted from left to right,
  // with object on the left having the highest priority.
  eventQueue.Put(e1);
  eventQueue.Put(e2);
}

func processPolygon(polygon Polygon, isSubject bool, queue *queue.PriorityQueue, bbox Bounds) {
for i := 0; i < polygon.Len(); i++ {
  path := polygon.At(i)
  for j := 0 ; j<path.Len(); j++ {
    processSegment(path.At(i), path.At(i + 1), isSubject, i +1, queue, bbox);
  }
}
}

func fillQueue(subject, clipping Polygon, sbbox, cbbox Bounds) *queue.PriorityQueue {
  eventQueue := queue.NewPriorityQueue(0, false) // (null, compareEvents);

  processPolygon(subject,  true,  eventQueue, sbbox);
  processPolygon(clipping, false, eventQueue, cbbox);

  return eventQueue;
}


func computeFields(event, prev *SweepEvent, sweepLine *rbtree.Tree, op operation) {
  // compute inOut and otherInOut fields
  if (prev == nil) {
    event.inOut      = false;
    event.otherInOut = true;

  // previous line segment in sweepline belongs to the same polygon
  } else if (event.isSubject == prev.isSubject) {
    event.inOut      = !prev.inOut;
    event.otherInOut = prev.otherInOut;

  // previous line segment in sweepline belongs to the clipping polygon
  } else {
    event.inOut      = !prev.otherInOut;
    if prev.isVertical() {
      event.otherInOut =   !prev.inOut
      }
      event.otherInOut = prev.inOut;
  }

  // compute prevInResult field
  if (prev != nil) {
    if !inResult(prev, operation) || prev.isVertical() {
    event.prevInResult = prev.prevInResult
  }
    event.prevInResult = prev; // TODO: not sure what this means.
  }
  // check if the line segment belongs to the Boolean operation
  event.inResult = inResult(event, operation);
}


func inResult(event *SweepEvent, op operation) {
  switch (event.edgeType) {
    case NORMAL:
      switch (operation) {
        case INTERSECTION:
          return !event.otherInOut;
        case UNION:
          return event.otherInOut;
        case DIFFERENCE:
          return (event.isSubject && event.otherInOut) ||
                 (!event.isSubject && !event.otherInOut);
        case XOR:
          return true;
      }
    case SAME_TRANSITION:
      return operation == INTERSECTION || operation == UNION;
    case DIFFERENT_TRANSITION:
      return operation == DIFFERENCE;
    case edgeType.NON_CONTRIBUTING:
      return false;
  }
  return false;
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

  p1, p2, nintersections := segmentIntersection(se1.point, se1.otherEvent.point, se2.point, se2.otherEvent.point, false); // TODO: I think this should be false but not sure.

  if (nintersections == 0) {return 0} // no intersection

  // the line segments intersect at an endpoint of both line segments
  if ((nintersections == 1) &&
      (equals(se1.point, se2.point) ||
       equals(se1.otherEvent.point, se2.otherEvent.point))) {
    return 0;
  }

  if (nintersections == 2 && se1.isSubject == se2.isSubject){
    if(se1.contourId == se2.contourId){
    panic(fmt.Errorf("Edges of the same polygon overlap, %+v, %+v, %+v, %+v",
      se1.point, se1.otherEvent.point, se2.point, se2.otherEvent.point));
    }
    return 0;
  }

  // The line segments associated to se1 and se2 intersect
  if (nintersections == 1) {

    // if the intersection point is not an endpoint of se1
    if (!equals(se1.point, p1) && !equals(se1.otherEvent.point, p1)) {
      divideSegment(se1, p1, queue);
    }

    // if the intersection point is not an endpoint of se2
    if (!equals(se2.point, p1) && !equals(se2.otherEvent.point, p1)) {
      divideSegment(se2, p1, queue);
    }
    return 1;
  }

  // The line segments associated to se1 and se2 overlap
  var events        []*SweepEvent;
  var leftCoincide  = false;
  var rightCoincide = false;

  if (equals(se1.point, se2.point)) {
    leftCoincide = true; // linked
  } else if (compareEvents(se1, se2) == 1) {
    events = []*SweepEvent{se2, se1};
  } else {
    events = []*SweepEvent{se1, se2};
  }

  if (equals(se1.otherEvent.point, se2.otherEvent.point)) {
    rightCoincide = true;
  } else if (compareEvents(se1.otherEvent, se2.otherEvent) == 1) {
    events = append(events, se2.otherEvent, se1.otherEvent);
  } else {
    events = append(events, se1.otherEvent, se2.otherEvent);
  }

  if ((leftCoincide && rightCoincide) || leftCoincide) {
    // both line segments are equal or share the left endpoint
    se1.edgeType = edgeType.NON_CONTRIBUTING;
    if (se1.inOut == se2.inOut) {
    se2.edgeType =  edgeType.SAME_TRANSITION
  }
  se2.edgeType =    edgeType.DIFFERENT_TRANSITION;

    if (leftCoincide && !rightCoincide) {
      // honestly no idea, but changing events selection from [2, 1]
      // to [0, 1] fixes the overlapping self-intersecting polygons issue
      divideSegment(events[0].otherEvent, events[1].point, queue);
    }
    return 2;
  }

  // the line segments share the right endpoint
  if (rightCoincide) {
    divideSegment(events[0], events[1].point, queue);
    return 3;
  }

  // no line segment includes totally the other one
  if (events[0] != events[3].otherEvent) {
    divideSegment(events[0], events[1].point, queue);
    divideSegment(events[1], events[2].point, queue);
    return 3;
  }

  // one line segment includes the other one
  divideSegment(events[0], events[1].point, queue);
  divideSegment(events[3].otherEvent, events[2].point, queue);

  return 3;
}


/**
 * @param  {SweepEvent} se
 * @param  {Array.<Number>} p
 * @param  {Queue} queue
 * @return {Queue}
 */
func divideSegment(se *SweepEvent, p Point, queue *queue.PriorityQueue)  {
  var r = NewSweepEvent(p, false, se,            se.isSubject);
  var l = NewSweepEvent(p, true,  se.otherEvent, se.isSubject);

  if (equals(se.point, se.otherEvent.point)) {
    panic(fmt.Errorf("what is that? %+v", se));
  }

  r.contourId,  l.contourId = se.contourId, se.contourId;

  // avoid a rounding error. The left event would be processed after the right event
  if (compareEvents(l, se.otherEvent) > 0) {
    se.otherEvent.left = true;
    l.left = false;
  }

  // avoid a rounding error. The left event would be processed after the right event
  // if (compareEvents(se, r) > 0) {}

  se.otherEvent.otherEvent = l;
  se.otherEvent = r;

  queue.push(l);
  queue.push(r);

  return queue;
}


func subdivideSegments(eventQueue *queue.PriorityQueue, subject, clipping Polygon, sbbox, cbbox Bounds, op operation) {
  sweepLine := rbtree.NewTree(compareSegments);
  var sortedEvents []*SweepEvent;

  var rightbound = min(sbbox.Max.X, cbbox.Max.X);

  var prev, next *SweepEvent;

  for (eventQueue.Len > 0) {
    var event = eventQueue.Get(1);
    sortedEvents = append(sortedEvents, event);

    // optimization by bboxes for intersection and difference goes here
    if ((op == INTERSECTION && event.point.X > rightbound) ||
        (op == DIFFERENCE   && event.point.X > sbbox.Max.X)) {
      break;
    }

    if (event.left) {
      sweepLine.Insert(event);
      // _renderSweepLine(sweepLine, event.point, event);

      next = sweepLine.Get(event);
      prev = sweepLine.Get(event);
      event.iterator = sweepLine.Get(event);

      // Cannot get out of the tree what we just put there
      if (!prev || !next) {
        panic("brute");
        var iterators = findIterBrute(sweepLine);
        prev = iterators[0];
        next = iterators[1];
      }

      if (!prev.Min() ) {
        prev.Prev();
      } else {
        prev = sweepLine.iterator(); //findIter(sweepLine.max());
        prev.prev();
        prev.next();
      }
      next.next();

      computeFields(event, prev.data(), sweepLine, op);

      if (next.data()) {
        if (possibleIntersection(event, next.data(), eventQueue) == 2) {
          computeFields(event, prev.data(), sweepLine, operation);
          computeFields(event, next.data(), sweepLine, operation);
        }
      }

      if (prev.data()) {
        if (possibleIntersection(prev.data(), event, eventQueue) == 2) {
          var prevprev = sweepLine.findIter(prev.data());
          if (!prevprev.Min()) {
            prevprev.prev();
          } else {
            prevprev = sweepLine.findIter(sweepLine.max());
            prevprev.next();
          }
          computeFields(prev.data(), prevprev.data(), sweepLine, op);
          computeFields(event, prev.data(), sweepLine, op);
        }
      }
    } else {
      event = event.otherEvent;
      next = sweepLine.findIter(event);
      prev = sweepLine.findIter(event);

      // _renderSweepLine(sweepLine, event.otherEvent.point, event);

      if (!(prev && next)) {continue};

      if (!prev.Min() ) {
        prev.prev();
      } else {
        prev = sweepLine.iterator();
        prev.prev(); // sweepLine.findIter(sweepLine.max());
        prev.next();
      }
      next.next();
      sweepLine.remove(event);

      //_renderSweepLine(sweepLine, event.otherEvent.point, event);

      if (next.data() && prev.data()) {
        possibleIntersection(prev.data(), next.data(), eventQueue);
      }
    }
  }
  return sortedEvents;
}

func findIterBrute(sweepLine *rbtree.Tree) {//, q)
  var prev = sweepLine.iterator();
  var next = sweepLine.iterator();
  var it   = sweepLine.iterator()//, data;
  for  {
    data := it.next()
    if data == nil {
      break
    }
    prev.next();
    next.next();
    if (data == event) {
      break;
    }
  }
  return prev, next;
}




func changeOrientation(contour Path) {
  return contour.reverse();
}





/**
 * @param  {Array.<SweepEvent>} sortedEvents
 * @return {Array.<SweepEvent>}
 */
func orderEvents(sortedEvents) {
  var event, i, len int
  var resultEvents  []*SweepEvent;
  for (i = 0, len = sortedEvents.length; i < len; i++) {
    event = sortedEvents[i];
    if ((event.left && event.inResult) ||
      (!event.left && event.otherEvent.inResult)) {
      resultEvents.push(event);
    }
  }

  // Due to overlapping edges the resultEvents array can be not wholly sorted
  var sorted = false;
  while (!sorted) {
    sorted = true;
    for (i = 0, len = resultEvents.length; i < len; i++) {
      if ((i + 1) < len &&
        compareEvents(resultEvents[i], resultEvents[i + 1]) === 1) {
        resultEvents[i], resultEvents[i+1] = resultEvents[i+1], resultEvents[i]
        sorted = false;
      }
    }
  }

  for (i = 0, len = resultEvents.length; i < len; i++) {
    resultEvents[i].pos = i;
  }

  for (i = 0, len = resultEvents.length; i < len; i++) {
    if (!resultEvents[i].left) {
      var temp = resultEvents[i].pos;
      resultEvents[i].pos = resultEvents[i].otherEvent.pos;
      resultEvents[i].otherEvent.pos = temp;
    }
  }

  return resultEvents;
}


/**
 * @param  {Array.<SweepEvent>} sortedEvents
 * @return {Array.<*>} polygons
 */
function connectEdges(sortedEvents) {
  var i, len;
  var resultEvents = orderEvents(sortedEvents);


  // "false"-filled array
  var processed = Array(resultEvents.length);
  var result = [];

  var depth  = [];
  var holeOf = [];
  var isHole = {};

  for (i = 0, len = resultEvents.length; i < len; i++) {
    if (processed[i]) continue;

    var contour = [];
    result.push(contour);

    var ringId = result.length - 1;
    depth.push(0);
    holeOf.push(-1);


    if (resultEvents[i].prevInResult) {
      var lowerContourId = resultEvents[i].prevInResult.contourId;
      if (!resultEvents[i].prevInResult.resultInOut) {
        addHole(result[lowerContourId], ringId);
        holeOf[ringId] = lowerContourId;
        depth[ringId]  = depth[lowerContourId] + 1;
        isHole[ringId] = true;
      } else if (isHole[lowerContourId]) {
        addHole(result[holeOf[lowerContourId]], ringId);
        holeOf[ringId] = holeOf[lowerContourId];
        depth[ringId]  = depth[lowerContourId];
        isHole[ringId] = true;
      }
    }

    var pos = i;
    var initial = resultEvents[i].point;
    contour.push(initial);

    while (pos >= i) {
      processed[pos] = true;

      if (resultEvents[pos].left) {
        resultEvents[pos].resultInOut = false;
        resultEvents[pos].contourId   = ringId;
      } else {
        resultEvents[pos].otherEvent.resultInOut = true;
        resultEvents[pos].otherEvent.contourId   = ringId;
      }

      pos = resultEvents[pos].pos;
      processed[pos] = true;

      contour.push(resultEvents[pos].point);
      pos = nextPos(pos, resultEvents, processed);
    }

    pos = pos === -1 ? i : pos;

    processed[pos] = processed[resultEvents[pos].pos] = true;
    resultEvents[pos].otherEvent.resultInOut = true;
    resultEvents[pos].otherEvent.contourId   = ringId;


    // depth is even
    /* eslint-disable no-bitwise */
    if (depth[ringId] & 1) {
      changeOrientation(contour);
    }
    /* eslint-enable no-bitwise */
  }

  return result;
}


/**
 * @param  {Number} pos
 * @param  {Array.<SweepEvent>} resultEvents
 * @param  {Array.<Boolean>}    processed
 * @return {Number}
 */
function nextPos(pos, resultEvents, processed) {
  var newPos = pos + 1;
  var length = resultEvents.length;
  while (newPos < length &&
         equals(resultEvents[newPos].point, resultEvents[pos].point)) {
    if (!processed[newPos]) {
      return newPos;
    } else {
      newPos = newPos + 1;
    }
  }

  newPos = pos - 1;

  while (processed[newPos]) {
    newPos = newPos - 1;
  }
  return newPos;
}


function trivialOperation(subject, clipping, operation) {
  var result = null;
  if (subject.length * clipping.length === 0) {
    if (operation === INTERSECTION) {
      result = EMPTY;
    } else if (operation === DIFFERENCE) {
      result = subject;
    } else if (operation === UNION || operation === XOR) {
      result = (subject.length === 0) ? clipping : subject;
    }
  }
  return result;
}


function compareBBoxes(subject, clipping, sbbox, cbbox, operation) {
  var result = null;
  if (sbbox[0] > cbbox[2] ||
      cbbox[0] > sbbox[2] ||
      sbbox[1] > cbbox[3] ||
      cbbox[1] > sbbox[3]) {
    if (operation === INTERSECTION) {
      result = EMPTY;
    } else if (operation === DIFFERENCE) {
      result = subject;
    } else if (operation === UNION || operation === XOR) {
      result = subject.concat(clipping);
    }
  }
  return result;
}


function boolean(subject, clipping, operation) {
  var trivial = trivialOperation(subject, clipping, operation);
  if (trivial) {
    return trivial === EMPTY ? null : trivial;
  }
  var sbbox = [Infinity, Infinity, -Infinity, -Infinity];
  var cbbox = [Infinity, Infinity, -Infinity, -Infinity];

  var eventQueue = fillQueue(subject, clipping, sbbox, cbbox);

  trivial = compareBBoxes(subject, clipping, sbbox, cbbox, operation);
  if (trivial) {
    return trivial === EMPTY ? null : trivial;
  }
  var sortedEvents = subdivideSegments(eventQueue, subject, clipping, sbbox, cbbox, operation);
  return connectEdges(sortedEvents);
}


module.exports = boolean;


module.exports.union = function(subject, clipping) {
  return boolean(subject, clipping, UNION);
};


module.exports.diff = function(subject, clipping) {
  return boolean(subject, clipping, DIFFERENCE);
};


module.exports.xor = function(subject, clipping) {
  return boolean(subject, clipping, XOR);
};


module.exports.intersection = function(subject, clipping) {
  return boolean(subject, clipping, INTERSECTION);
};


/**
 * @enum {Number}
 */
module.exports.operations = {
  INTERSECTION: INTERSECTION,
  DIFFERENCE:   DIFFERENCE,
  UNION:        UNION,
  XOR:          XOR
};


// for testing
module.exports.fillQueue            = fillQueue;
module.exports.computeFields        = computeFields;
module.exports.subdivideSegments    = subdivideSegments;
module.exports.divideSegment        = divideSegment;
module.exports.possibleIntersection = possibleIntersection;
