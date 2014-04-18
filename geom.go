// Copyright (c) 2011 Mateusz Czapliński (Go port)
// Copyright (c) 2011 Mahir Iqbal (as3 version)
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// based on http://code.google.com/p/as3polyclip/ (MIT licensed)
// and code by Martínez et al: http://wwwdi.ujaen.es/~fmartin/bool_op.html (public domain)

// Package geomop provides implementation of algorithms for geometry operations.
// For further details, consult the description of Polygon.Construct method.
package geomop

import (
	"github.com/twpayne/gogeom/geom"
	"math"
	"reflect"
)

// Equals returns true if both p1 and p2 describe exactly the same point.
func PointEquals(p1, p2 geom.Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

func pointSubtract(p1, p2 geom.Point) geom.Point {
	return geom.Point{p1.X - p2.X, p1.Y - p2.Y}
}

// Length returns distance from p to point (0, 0).
func lengthToOrigin(p geom.Point) float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Used to represent an edge of a polygon.
type segment struct {
	start, end geom.Point
}

// Contour represents a sequence of vertices connected by line segments, forming a closed shape.
type Contour []geom.Point

func (c Contour) segment(index int) segment {
	if index == len(c)-1 {
		return segment{c[len(c)-1], c[0]}
	}
	return segment{c[index], c[index+1]}
	// if out-of-bounds, we expect panic detected by runtime
}

// Checks if a point is inside a contour using the "point in polygon" raycast method.
// This works for all polygons, whether they are clockwise or counter clockwise,
// convex or concave.
// See: http://en.wikipedia.org/wiki/Point_in_polygon#Ray_casting_algorithm
// Returns true if p is inside the polygon defined by contour.
func (c Contour) Contains(p geom.Point) bool {
	// Cast ray from p.x towards the right
	intersections := 0
	for i := range c {
		curr := c[i]
		ii := i + 1
		if ii == len(c) {
			ii = 0
		}
		next := c[ii]

		if (p.Y >= next.Y || p.Y <= curr.Y) &&
			(p.Y >= curr.Y || p.Y <= next.Y) {
			continue
		}
		// Edge is from curr to next.

		if p.X >= math.Max(curr.X, next.X) ||
			next.Y == curr.Y {
			continue
		}

		// Find where the line intersects...
		xint := (p.Y-curr.Y)*(next.X-curr.X)/(next.Y-curr.Y) + curr.X
		if curr.X != next.X && p.X > xint {
			continue
		}
		intersections++
	}
	return intersections%2 != 0
}

// Clone returns a copy of a contour.
func (c Contour) Clone() Contour {
	return append([]geom.Point{}, c...)
}

// NumVertices returns total number of all vertices of all contours of a polygon.
func NumVertices(p geom.Polygon) int {
	num := 0
	for _, c := range p.Rings {
		num += len(c)
	}
	return num
}

// Clone returns a duplicate of a polygon.
func Clone(p geom.Polygon) geom.Polygon {
	var r geom.Polygon
	r.Rings = make([][]geom.Point, len(p.Rings))
	for i, rr := range p.Rings {
		r.Rings[i] = make([]geom.Point, len(rr))
		for j, pp := range p.Rings[i] {
			r.Rings[i][j] = pp
		}
	}
	return r
}

// Op describes an operation which can be performed on two polygons.
type Op int

const (
	UNION Op = iota
	INTERSECTION
	DIFFERENCE
	XOR
)

// Construct computes a 2D polygon, which is a result of performing
// specified Boolean operation on the provided pair of polygons (p <Op> clipping).
// It uses algorithm described by F. Martínez, A. J. Rueda, F. R. Feito
// in "A new algorithm for computing Boolean operations on polygons"
// - see: http://wwwdi.ujaen.es/~fmartin/bool_op.html
// The paper describes the algorithm as performing in time O((n+k) log n),
// where n is number of all edges of all polygons in operation, and
// k is number of intersections of all polygon edges.
// "subject" and "clipping" can both be of type geom.Polygon,
// geom.MultiPolygon, geom.LineString, or geom.MultiLineString.
func Construct(subject, clipping geom.T, operation Op) geom.T {
	// Prepare the input shapes
	var c clipper
	switch clipping.(type) {
	case geom.Polygon, geom.MultiPolygon:
		c.subject = convertToPolygon(subject)
		c.clipping = convertToPolygon(clipping)
		switch subject.(type) {
		case geom.Polygon, geom.MultiPolygon:
			c.outType = outputPolygons
		case geom.LineString, geom.MultiLineString:
			c.outType = outputLines
		}

	case geom.LineString, geom.MultiLineString:
		switch subject.(type) {
		case geom.Polygon, geom.MultiPolygon:
			// swap clipping and subject
			c.subject = convertToPolygon(clipping)
			c.clipping = convertToPolygon(subject)
			c.outType = outputLines
		case geom.LineString, geom.MultiLineString:
			c.subject = convertToPolygon(subject)
			c.clipping = convertToPolygon(clipping)
			c.outType = outputPoints
		}
	}
	// Run the clipper
	return c.compute(operation)
}

// convert input shapes to polygon to make internal processing easier
func convertToPolygon(g geom.T) geom.Polygon {
	switch g.(type) {
	case geom.Polygon:
		return g.(geom.Polygon)
	case geom.MultiPolygon:
		var out geom.Polygon
		out.Rings = make([][]geom.Point, 0)
		for _, p := range g.(geom.MultiPolygon).Polygons {
			for _, r := range p.Rings {
				out.Rings = append(out.Rings, r)
			}
		}
		return out
	case geom.LineString:
		var out geom.Polygon
		g2 := g.(geom.LineString)
		out.Rings = make([][]geom.Point, 1)
		out.Rings[0] = make([]geom.Point, len(g2.Points))
		for j, p := range g2.Points {
			out.Rings[0][j] = p
		}
		return out
	case geom.MultiLineString:
		var out geom.Polygon
		g2 := g.(geom.MultiLineString)
		out.Rings = make([][]geom.Point, len(g2.LineStrings))
		for i, ls := range g2.LineStrings {
			out.Rings[i] = make([]geom.Point, len(ls.Points))
			for j, p := range ls.Points {
				out.Rings[i][j] = p
			}
		}
		return out
	default:
		panic(NewError(g))
	}
}

type UnsupportedGeometryError struct {
	Type reflect.Type
}

func NewError(g geom.T) UnsupportedGeometryError {
	return UnsupportedGeometryError{reflect.TypeOf(g)}
}

func (e UnsupportedGeometryError) Error() string {
	return "Unsupported geometry type: " + e.Type.String()
}
