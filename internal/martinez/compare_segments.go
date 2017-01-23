package martinez

// compareSegments compares two segments and returns 1 if le1 > le2,
// -1 if le2 > le1, and 0 if le1 == le2
func compareSegments(le1, le2 SweepEvent) int {
	if le1 == le2 {
		return 0
	}

	// Segments are not collinear
	if signedArea(le1.point, le1.otherEvent.point, le2.point) != 0 ||
		signedArea(le1.point, le1.otherEvent.point, le2.otherEvent.point) != 0 {

		// If they share their left endpoint use the right endpoint to sort
		if equals(le1.point, le2.point) {
			if le1.isBelow(le2.otherEvent.point) {
				return -1
			}
			return 1
		}

		// Different left endpoint: use the left endpoint to sort
		if le1.point.X() == le2.point.X() {
			if le1.point.Y() < le2.point.Y() {
				return -1
			}
			return 1
		}

		// has the line segment associated to e1 been inserted
		// into S after the line segment associated to e2 ?
		if compareEvents(le1, le2) == 1 {
			if le2.isAbove(le1.point) {
				return -1
			}
			return 1
		}

		// The line segment associated to e2 has been inserted
		// into S after the line segment associated to e1
		if le1.isBelow(le2.point) {
			return -1
		}
		return 1
	}

	if le1.isSubject == le2.isSubject { // same polygon
		if equals(le1.point, le2.point) {
			if equals(le1.otherEvent.point, le2.otherEvent.point) {
				return 0
			} else {
				if le1.contourId > le2.contourId {
					return 1
				}
				return -1
			}
		}
	} else { // Segments are collinear, but belong to separate polygons
		if le1.isSubject {
			return -1
		}
		return 1
	}

	if compareEvents(le1, le2) == 1 {
		return 1
	}
	return -1
}
