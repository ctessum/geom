package geom

import "testing"

func TestArea(t *testing.T) {
	tests := []struct {
		test     Polygonal
		expected float64
	}{
		{
			test: Polygon{{
				Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
				Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
				Point{X: 0, Y: 0}}},
			expected: 4.,
		},
		{
			test: Polygon{
				{
					Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
					Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
					Point{X: 0, Y: 0},
				},
				{
					Point{X: 0.5, Y: 0.5}, Point{X: 0.5, Y: 1.5},
					Point{X: 1.5, Y: 1.5}, Point{X: 1.5, Y: 0.5},
					Point{X: 0.5, Y: 0.5},
				},
			},
			expected: 3.,
		},
		{
			test: MultiPolygon{
				{{
					Point{X: 0, Y: 0}, Point{X: 2, Y: 0},
					Point{X: 2, Y: 2}, Point{X: 0, Y: 2},
					Point{X: 0, Y: 0},
				}},
				{{
					Point{X: 2, Y: 2}, Point{X: 4, Y: 2},
					Point{X: 4, Y: 4}, Point{X: 2, Y: 4},
					Point{X: 2, Y: 2},
				}},
			},
			expected: 8.,
		},
	}
	for i, test := range tests {
		result := test.test.Area()
		if result != test.expected {
			t.Errorf("%d: expected %g, got %g", i, test.expected, result)
		}
	}
}
