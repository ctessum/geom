package geom

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBounds_Within(t *testing.T) {
	b := &Bounds{
		Min: Point{0, 0},
		Max: Point{1, 1},
	}
	tests := []struct {
		p  Polygonal
		in WithinStatus
	}{
		{
			p: &Bounds{
				Min: Point{0, 0},
				Max: Point{1, 1},
			},
			in: OnEdge,
		},
		{
			p: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{1, 1},
			},
			in: Outside,
		},
		{
			p: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{0.6, 0.6},
			},
			in: Outside,
		},
		{
			p: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{2, 2},
			},
			in: Outside,
		},
		{
			p: &Bounds{
				Min: Point{1, 1},
				Max: Point{2, 2},
			},
			in: Outside,
		},
		{
			p: &Bounds{
				Min: Point{1.5, 1.5},
				Max: Point{2, 2},
			},
			in: Outside,
		},
		{
			p: &Bounds{
				Min: Point{-1, -1},
				Max: Point{1, 1},
			},
			in: Inside,
		},
		{
			p: &Bounds{
				Min: Point{-1, -1},
				Max: Point{2, 2},
			},
			in: Inside,
		},
		{
			p:  Polygon{{{0, 0}, {1, 1}, {0, 1}}},
			in: Inside,
		},
		{
			p:  Polygon{{{0, 0}, {1, 1}, {-1, 1}}},
			in: Inside,
		},
		{
			p:  Polygon{{{0.2, 0.2}, {0.7, 0.7}, {0.2, 0.7}}},
			in: Outside,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			in := b.Within(test.p)
			if in != test.in {
				t.Errorf("%v != %v", in, test.in)
			}
		})
	}
}

func TestBounds_Polygons(t *testing.T) {
	b := &Bounds{
		Min: Point{0, 0},
		Max: Point{1, 1},
	}
	p := b.Polygons()
	want := []Polygon{{{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}}}}
	if !reflect.DeepEqual(p, want) {
		t.Errorf("%v != %v", p, want)
	}
}

func TestBounds_Intersection(t *testing.T) {
	b := &Bounds{
		Min: Point{0, 0},
		Max: Point{1, 1},
	}
	tests := []struct {
		p Polygonal
		i Polygonal
	}{
		{
			p: &Bounds{
				Min: Point{0, 0},
				Max: Point{1, 1},
			},
			i: b,
		},
		{
			p: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{1, 1},
			},
			i: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{1, 1},
			},
		},
		{
			p: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{0.6, 0.6},
			},
			i: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{0.6, 0.6},
			},
		},
		{
			p: &Bounds{
				Min: Point{0.5, 0.5},
				Max: Point{2, 2},
			},
			i: Polygon{{{0.5, 1}, {0.5, 0.5}, {1, 0.5}, {1, 1}, {0.5, 1}}},
		},
		{
			p: &Bounds{
				Min: Point{1, 1},
				Max: Point{2, 2},
			},
			i: Polygon{},
		},
		{
			p: &Bounds{
				Min: Point{1.5, 1.5},
				Max: Point{2, 2},
			},
			i: Polygon{},
		},
		{
			p: &Bounds{
				Min: Point{-1, -1},
				Max: Point{1, 1},
			},
			i: b,
		},
		{
			p: &Bounds{
				Min: Point{-1, -1},
				Max: Point{2, 2},
			},
			i: b,
		},
		{
			p: Polygon{{{0, 0}, {1, 1}, {0, 1}}},
			i: Polygon{{{0, 0}, {1, 1}, {0, 1}}},
		},
		{
			p: Polygon{{{0, 0}, {1, 1}, {-1, 1}}},
			i: Polygon{{{0, 1}, {0, 0}, {1, 1}, {0, 1}}},
		},
		{
			p: Polygon{{{0.2, 0.2}, {0.7, 0.7}, {0.2, 0.7}}},
			i: Polygon{{{0.2, 0.2}, {0.7, 0.7}, {0.2, 0.7}}},
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			i := b.Intersection(test.p)
			if i == nil && test.i.Len() == 0 {
				return
			}
			if !reflect.DeepEqual(i, test.i) {
				t.Errorf("%v != %v", i, test.i)
			}
		})
	}
}

func TestBounds_Area(t *testing.T) {
	b := &Bounds{
		Min: Point{0, 0},
		Max: Point{1, 1},
	}
	a := b.Area()
	if a != 1 {
		t.Errorf("%v != %v", a, 1)
	}
}

func TestBounds_Centroid(t *testing.T) {
	b := &Bounds{
		Min: Point{0, 0},
		Max: Point{1, 1},
	}
	c := b.Centroid()
	want := Point{0.5, 0.5}
	if !reflect.DeepEqual(c, want) {
		t.Errorf("%v != %v", c, want)
	}
}
