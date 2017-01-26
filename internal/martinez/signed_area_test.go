package martinez

import "testing"

func TestSignedArea(t *testing.T) {
	tests := []struct {
		points []point
		result float64
		msg    string
	}{
		{
			points: []point{point{x: 0, y: 0}, point{x: 0, y: 1}, point{x: 1, y: 1}},
			result: -1,
			msg:    "negative area",
		},
		{
			points: []point{point{x: 0, y: 1}, point{x: 0, y: 0}, point{x: 1, y: 0}},
			result: 1,
			msg:    "positive area",
		},
		{
			points: []point{point{x: 0, y: 0}, point{x: 1, y: 1}, point{x: 2, y: 2}},
			result: 0,
			msg:    "collinear, 0 area",
		},
		{
			points: []point{point{x: -1, y: 0}, point{x: 2, y: 3}, point{x: 0, y: 1}},
			result: 0,
			msg:    "point on segment",
		},
		{
			points: []point{point{x: 2, y: 3}, point{x: -1, y: 0}, point{x: 0, y: 1}},
			result: 0,
			msg:    "point on segment",
		},
	}
	for _, test := range tests {
		if signedArea(test.points[0], test.points[1], test.points[2]) != test.result {
			t.Error(test.msg)
		}
	}
}
