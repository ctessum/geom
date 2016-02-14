package geom

import "testing"

func TestLength(t *testing.T) {
	tests := []struct {
		test     Linear
		expected float64
	}{
		{
			test: MultiLineString{
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
			expected: 12.,
		},
	}
	for i, test := range tests {
		result := test.test.Length()
		if result != test.expected {
			t.Errorf("%d: expected %g, got %g", i, test.expected, result)
		}
	}
}
