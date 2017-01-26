package martinez

import "github.com/Workiva/go-datastructures/queue"

// Compare compares two SweepEvents, returning
// 1 if e1 > e2 and -1 if e2 > e1.
func (e1 *SweepEvent) Compare(e2I queue.Item) int {
	e2 := e2I.(*SweepEvent)
	var p1 = e1.point
	var p2 = e2.point

	// Different x-coordinate
	if p1.X() > p2.X() {
		return 1
	}
	if p1.X() < p2.X() {
		return -1
	}

	// Different points, but same x-coordinate
	// Event with lower y-coordinate is processed first
	if p1.Y() != p2.Y() {
		if p1.Y() > p2.Y() {
			return 1
		}
		return -1
	}

	return specialCases(e1, e2, p1, p2)
}

func specialCases(e1, e2 *SweepEvent, p1, p2 Point) int {
	// Same coordinates, but one is a left endpoint and the other is
	// a right endpoint. The right endpoint is processed first
	if e1.left != e2.left {
		if e1.left {
			return 1
		}
		return -1
	}

	// Same coordinates, both events
	// are left endpoints or right endpoints.
	// not collinear
	if e1.otherEvent != nil && e2.otherEvent != nil {
		if signedArea(p1, e1.otherEvent.point, e2.otherEvent.point) != 0 {
			// the event associate to the bottom segment is processed first
			if !e1.isBelow(e2.otherEvent.point) {
				return 1
			}
			return -1
		}
	}

	// uncomment this if you want to play with multipolygons
	// if (e1.isSubject === e2.isSubject) {
	//   if(equals(e1.point, e2.point) && e1.contourId === e2.contourId) {
	//     return 0;
	//   } else {
	//     return e1.contourId > e2.contourId ? 1 : -1;
	//   }
	// }

	if e1.isSubject != e2.isSubject {
		if !e1.isSubject && e2.isSubject {
			return 1
		}
		return -1
	}
	return 0
}
