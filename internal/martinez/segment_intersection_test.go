package martinez

import "testing"

func TestSegmentIntersection(t *testing.T) {
	tests := []struct {
		p1, p2, p3, p4  point
		noEndpointTouch bool
		n               int
		r1, r2          point
		msg             string
	}{
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 1, y: 0}, p4: point{x: 2, y: 2},
			n:   0,
			msg: "no intersections",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 1, y: 0}, p4: point{x: 10, y: 2},
			n:   0,
			msg: "no intersection",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 1, y: 0}, p4: point{x: 0, y: 1},
			n:   1,
			r1:  point{x: 0.5, y: 0.5},
			msg: "1 intersection",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: 1}, p4: point{x: 0, y: 0},
			n:   1,
			r1:  point{x: 0, y: 0},
			msg: "shared point 1",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: 1}, p4: point{x: 1, y: 1},
			n:   1,
			r1:  point{x: 1, y: 1},
			msg: "shared point 2",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0.5, y: 0.5}, p4: point{x: 1, y: 0},
			n:   1,
			r1:  point{x: 0.5, y: 0.5},
			msg: "T-crossing",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 10, y: 10}, p3: point{x: 1, y: 1}, p4: point{x: 5, y: 5},
			n:  2,
			r1: point{x: 1, y: 1}, r2: point{x: 5, y: 5},
			msg: "full overlap",
		},
		{
			p1: point{x: 1, y: 1}, p2: point{x: 10, y: 10}, p3: point{x: 1, y: 1}, p4: point{x: 5, y: 5},
			n:  2,
			r1: point{x: 1, y: 1}, r2: point{x: 5, y: 5},
			msg: "shared point + overlap",
		},
		{
			p1: point{x: 3, y: 3}, p2: point{x: 10, y: 10}, p3: point{x: 0, y: 0}, p4: point{x: 5, y: 5},
			n:  2,
			r1: point{x: 3, y: 3}, r2: point{x: 5, y: 5},
			msg: "mutual overlap",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: 0}, p4: point{x: 1, y: 1},
			n:  2,
			r1: point{x: 0, y: 0}, r2: point{x: 1, y: 1},
			msg: "full overlap",
		},
		{
			p1: point{x: 1, y: 1}, p2: point{x: 0, y: 0}, p3: point{x: 0, y: 0}, p4: point{x: 1, y: 1},
			n:  2,
			r1: point{x: 1, y: 1}, r2: point{x: 0, y: 0},
			msg: "full overlap, orientation",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 1, y: 1}, p4: point{x: 2, y: 2},
			n:   1,
			r1:  point{x: 1, y: 1},
			msg: "collinear, shared point",
		},
		{
			p1: point{x: 1, y: 1}, p2: point{x: 0, y: 0}, p3: point{x: 1, y: 1}, p4: point{x: 2, y: 2},
			n:   1,
			r1:  point{x: 1, y: 1},
			msg: "collinear, shared other point",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 2, y: 2}, p4: point{x: 4, y: 4},
			n:   0,
			msg: "collinear, no overlap",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: -1}, p4: point{x: 1, y: 0},
			n:   0,
			msg: "parallel",
		},
		{
			p1: point{x: 1, y: 1}, p2: point{x: 0, y: 0}, p3: point{x: 0, y: -1}, p4: point{x: 1, y: 0},
			n:   0,
			msg: "parallel, orientation",
		},
		{
			p1: point{x: 0, y: -1}, p2: point{x: 1, y: 0}, p3: point{x: 0, y: 0}, p4: point{x: 1, y: 1},
			n:   0,
			msg: "parallel, position",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: 1}, p4: point{x: 0, y: 0},
			noEndpointTouch: true,
			n:               0,
			msg:             "shared point 1, skip touches",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: 1}, p4: point{x: 1, y: 1},
			noEndpointTouch: true,
			n:               0,
			msg:             "shared point 2, skip touches",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 1, y: 1}, p4: point{x: 2, y: 2},
			noEndpointTouch: true,
			n:               0,
			msg:             "collinear, shared point, skip touches",
		},
		{
			p1: point{x: 1, y: 1}, p2: point{x: 0, y: 0}, p3: point{x: 1, y: 1}, p4: point{x: 2, y: 2},
			noEndpointTouch: true,
			n:               0,
			msg:             "collinear, shared other point, skip touches",
		},
		{
			p1: point{x: 0, y: 0}, p2: point{x: 1, y: 1}, p3: point{x: 0, y: 0}, p4: point{x: 1, y: 1},
			noEndpointTouch: true,
			n:               0,
			msg:             "full overlap, skip touches",
		},
		{
			p1: point{x: 1, y: 1}, p2: point{x: 0, y: 0}, p3: point{x: 0, y: 0}, p4: point{x: 1, y: 1},
			noEndpointTouch: true,
			n:               0,
			msg:             "full overlap, orientation, skip touches",
		},
	}
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			i1, i2, n := segmentIntersection(test.p1, test.p2, test.p3, test.p4, test.noEndpointTouch)
			if n != test.n {
				t.Errorf("number intersections %d != %d", n, test.n)
			}
			if i1 != test.r1 {
				t.Errorf("i1 %+v != %+v", i1, test.r1)
			}
			if i2 != test.r2 {
				t.Errorf("i2 %+v != %+v", i2, test.r2)
			}
		})
	}
}
