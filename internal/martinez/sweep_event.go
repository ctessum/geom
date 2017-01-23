package martinez

// SweepEvent is a sweep line event wraps around a vertex point.
type SweepEvent struct {
	left         bool //Is left endpoint?
	point        Point
	otherEvent   *SweepEvent // other edge reference
	isSubject    bool        // Belongs to source or clipping polygon
	edgetype     edgeType    // Edge contribution type (default is normal).
	inOut        bool        // In-out transition for the sweepline crossing polygon
	otherInOut   bool
	prevInResult *SweepEvent // Previous event in result?
	inResult     bool        // Does event belong to result?
	resultInOut  bool
	contourId    int
}

// NewSweepEvent creates a new sweepline event for point p.
func NewSweepEvent(point Point, left bool, otherEvent *SweepEvent, isSubject bool, edgeType edgeType) *SweepEvent {
	return &SweepEvent{
		left:       left, //Is left endpoint?
		point:      point,
		otherEvent: otherEvent, // other edge reference
		isSubject:  isSubject,  // Belongs to source or clipping polygon
		edgetype:   edgeType,   // Edge contribution type (default is normal).
		inOut:      false,      // In-out transition for the sweepline crossing polygon
		otherInOut: false,
	}
}

func (s *SweepEvent) isBelow(p Point) bool {
	if s.left {
		return signedArea(s.point, s.otherEvent.point, p) > 0
	}
	return signedArea(s.otherEvent.point, s.point, p) > 0
}

func (s *SweepEvent) isAbove(p Point) bool {
	return !s.isBelow(p)
}
func (s *SweepEvent) isVertical() bool {
	return s.point.X() == s.otherEvent.point.X()
}
