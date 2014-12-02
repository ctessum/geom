package geomop

import (
	"fmt"
	"sort"
	. "testing"

	"github.com/twpayne/gogeom/geom"
)

type sorter geom.Polygon

func (s sorter) Len() int      { return len(s.Rings) }
func (s sorter) Swap(i, j int) { s.Rings[i], s.Rings[j] = s.Rings[j], s.Rings[i] }
func (s sorter) Less(i, j int) bool {
	if len(s.Rings[i]) != len(s.Rings[j]) {
		return len(s.Rings[i]) < len(s.Rings[j])
	}
	for k := range s.Rings[i] {
		pi, pj := s.Rings[i][k], s.Rings[j][k]
		if pi.X != pj.X {
			return pi.X < pj.X
		}
		if pi.Y != pj.Y {
			return pi.Y < pj.Y
		}
	}
	return false
}

// basic normalization just for tests; to be improved if needed
func normalize(poly geom.Polygon) geom.Polygon {
	for i, c := range poly.Rings {
		if len(c) == 0 {
			continue
		}

		// find bottom-most of leftmost points, to have fixed anchor
		min := 0
		for j, p := range c {
			if p.X < c[min].X || p.X == c[min].X && p.Y < c[min].Y {
				min = j
			}
		}

		// rotate points to make sure min is first
		poly.Rings[i] = append(c[min:], c[:min]...)
	}

	sort.Sort(sorter(poly))

	var poly2 geom.Polygon
	poly2.Rings = make([][]geom.Point, len(poly.Rings))
	for i, r := range poly.Rings {
		poly2.Rings[i] = make([]geom.Point, 0, len(r))
		for j, p := range r {
			if j == 0 || !PointEquals(p, r[j-1]) {
				poly2.Rings[i] = append(poly2.Rings[i], p)
			}
		}
	}
	return poly2
}

func dump(poly geom.Polygon) string {
	return fmt.Sprintf("%v", normalize(poly))
}

func TestBug3(t *T) {
	cases := []struct{ subject, clipping, result geom.T }{
		// original reported github issue #3
		{
			subject: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 1}, {1, 2}, {2, 2}, {2, 1}}}}),
			clipping: geom.T(geom.Polygon{[][]geom.Point{
				{{2, 1}, {2, 2}, {3, 2}, {3, 1}},
				{{1, 2}, {1, 3}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 3}, {3, 2}}}}),
			result: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 1}, {2, 1}, {3, 1},
					{3, 2}, {3, 3},
					{2, 3}, {1, 3},
					{1, 2}}}}),
		},
		// simplified variant of issue #3, for easier debugging
		{
			subject: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 1}}}}),
			clipping: geom.T(geom.Polygon{[][]geom.Point{
				{{2, 1}, {2, 2}, {3, 2}},
				{{1, 2}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 2}}}}),
			result: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 3}, {3, 2}, {2, 1}}}}),
		},
		{
			subject: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 1}}}}),
			clipping: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 3}, {2, 2}},
				{{2, 2}, {2, 3}, {3, 2}}}}),
			result: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 3}, {3, 2}, {2, 2}, {2, 1}}}}),
		},
		// another variation, now with single degenerated curve
		{
			subject: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 1}}}}),
			clipping: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 3}, {2, 2}, {2, 3}, {3, 2}}}}),
			result: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 3}, {3, 2}, {2, 2}, {2, 1}}}}),
		},
		{
			subject: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 1}}}}),
			clipping: geom.T(geom.Polygon{[][]geom.Point{
				{{2, 1}, {2, 2}, {2, 3}, {3, 2}},
				{{1, 2}, {2, 3}, {2, 2}}}}),
			result: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 3}, {3, 2}, {2, 1}}}}),
		},
		// "union" with effectively empty polygon (wholly self-intersecting)
		{
			subject: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 1}}}}),
			clipping: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 3}, {1, 2}, {2, 2}, {2, 3}}}}),
			result: geom.T(geom.Polygon{[][]geom.Point{
				{{1, 2}, {2, 2}, {2, 1}}}}),
		},
	}
	for _, c := range cases {
		union, err := Construct(c.subject, c.clipping, UNION)
		handle(err)
		result := dump(union.(geom.Polygon))
		err = FixOrientation(c.result)
		handle(err)
		if result != dump(c.result.(geom.Polygon)) {
			t.Errorf("case UNION:\nsubject:  %v\nclipping: %v\nexpected: %v\ngot:      %v",
				c.subject, c.clipping, c.result, result)
		}
	}
}
