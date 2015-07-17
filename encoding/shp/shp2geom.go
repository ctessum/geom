package shp

import (
	"fmt"
	"math"
	"reflect"

	"github.com/ctessum/geom"
	"github.com/ctessum/geom/op"
	"github.com/jonas-p/go-shp"
)

// Shp2Geom converts a shapefile shape to a geometry
// object that can be used with other packages.
// This function can be used to wrap the go-shp "Shape()" method.
func shp2Geom(n int, s shp.Shape) (int, geom.T, error) {
	switch t := reflect.TypeOf(s); {
	case t == reflect.TypeOf(&shp.Point{}):
		return n, point2geom(*s.(*shp.Point)), nil
	case t == reflect.TypeOf(&shp.PointM{}):
		return n, pointM2geom(*s.(*shp.PointM)), nil
	case t == reflect.TypeOf(&shp.PointZ{}):
		return n, pointZ2geom(*s.(*shp.PointZ)), nil
	case t == reflect.TypeOf(&shp.Polygon{}):
		return n, polygon2geom(*s.(*shp.Polygon)), nil
	case t == reflect.TypeOf(&shp.PolygonM{}):
		return n, polygonM2geom(*s.(*shp.PolygonM)), nil
	case t == reflect.TypeOf(&shp.PolygonZ{}):
		return n, polygonZ2geom(*s.(*shp.PolygonZ)), nil
	case t == reflect.TypeOf(&shp.PolyLine{}):
		return n, polyLine2geom(*s.(*shp.PolyLine)), nil
	case t == reflect.TypeOf(&shp.PolyLineM{}):
		return n, polyLineM2geom(*s.(*shp.PolyLineM)), nil
	case t == reflect.TypeOf(&shp.PolyLineZ{}):
		return n, polyLineZ2geom(*s.(*shp.PolyLineZ)), nil
	//case t == "MultiPatch": // not yet supported
	case t == reflect.TypeOf(&shp.MultiPoint{}):
		return n, multiPoint2geom(*s.(*shp.MultiPoint)), nil
	case t == reflect.TypeOf(&shp.MultiPointM{}):
		return n, multiPointM2geom(*s.(*shp.MultiPointM)), nil
	case t == reflect.TypeOf(&shp.MultiPointZ{}):
		return n, multiPointZ2geom(*s.(*shp.MultiPointZ)), nil
	case t == reflect.TypeOf(&shp.Null{}):
		return n, nil, nil
	default:
		return n, nil, fmt.Errorf("Unsupported shape type: %v", t)
	}
}

// Functions for converting shp to geom

func point2geom(s shp.Point) geom.T {
	return geom.Point(s)
}
func pointM2geom(s shp.PointM) geom.T {
	return geom.PointM(s)
}
func pointZ2geom(s shp.PointZ) geom.T {
	return geom.PointZM(s)
}
func getStartEnd(parts []int32, points []shp.Point, i int) (start, end int) {
	start = int(parts[i])
	if i == len(parts)-1 {
		end = len(points)
	} else {
		end = int(parts[i+1])
	}
	return
}
func polygon2geom(s shp.Polygon) geom.T {
	var pg geom.Polygon = make([][]geom.Point, len(s.Parts))
	for i := 0; i < len(s.Parts); i++ {
		start, end := getStartEnd(s.Parts, s.Points, i)
		pg[i] = make([]geom.Point, end-start)
		// Go backwards around the rings to switch to OGC format
		for j := end - 1; j >= start; j-- {
			pg[i][j-start] = geom.Point(s.Points[j])
		}
	}
	// Make sure the winding direction is correct
	op.FixOrientation(pg)
	return pg
}
func polygonM2geom(s shp.PolygonM) geom.T {
	var pg geom.PolygonM = make([][]geom.PointM, len(s.Parts))
	jj := 0
	for i := 0; i < len(s.Parts); i++ {
		start, end := getStartEnd(s.Parts, s.Points, i)
		jj += end - start
		pg[i] = make([]geom.PointM, end-start)
		// Go backwards around the rings to switch to OGC format
		for j := end - 1; j >= start; j-- {
			ss := s.Points[j]
			pg[i][j-start] = geom.PointM{ss.X, ss.Y, s.MArray[jj]}
			jj--
		}
	}
	// Make sure the winding direction is correct
	op.FixOrientation(pg)
	return pg
}

func polygonZ2geom(s shp.PolygonZ) geom.T {
	var pg geom.PolygonZM = make([][]geom.PointZM, len(s.Parts))
	jj := -1
	for i := 0; i < len(s.Parts); i++ {
		start, end := getStartEnd(s.Parts, s.Points, i)
		jj += end - start
		pg[i] = make([]geom.PointZM, end-start)
		// Go backwards around the rings to switch to OGC format
		for j := end - 1; j >= start; j-- {
			ss := s.Points[j]
			pg[i][j-start] = geom.PointZM{ss.X, ss.Y, s.ZArray[jj],
				s.MArray[jj]}
			jj--
		}
	}
	// Make sure the winding direction is correct
	op.FixOrientation(pg)
	return pg
}
func polyLine2geom(s shp.PolyLine) geom.T {
	var pl geom.MultiLineString = make([]geom.LineString, len(s.Parts))
	for i := 0; i < len(s.Parts); i++ {
		start, end := getStartEnd(s.Parts, s.Points, i)
		pl[i] = make([]geom.Point, end-start)
		for j := start; j < end; j++ {
			pl[i][j-start] = geom.Point(s.Points[j])
		}
	}
	return pl
}
func polyLineM2geom(s shp.PolyLineM) geom.T {
	var pl geom.MultiLineStringM = make([]geom.LineStringM, len(s.Parts))
	jj := 0
	for i := 0; i < len(s.Parts); i++ {
		start, end := getStartEnd(s.Parts, s.Points, i)
		pl[i] = make([]geom.PointM, end-start)
		for j := start; j < end; j++ {
			ss := s.Points[j]
			pl[i][j-start] =
				geom.PointM{ss.X, ss.Y, s.MArray[jj]}
			jj++
		}
	}
	return pl
}
func polyLineZ2geom(s shp.PolyLineZ) geom.T {
	var pl geom.MultiLineStringZM = make([]geom.LineStringZM, len(s.Parts))
	jj := 0
	for i := 0; i < len(s.Parts); i++ {
		start, end := getStartEnd(s.Parts, s.Points, i)
		pl[i] = make([]geom.PointZM, end-start)
		for j := start; j < end; j++ {
			ss := s.Points[j]
			pl[i][j-start] =
				geom.PointZM{ss.X, ss.Y, s.ZArray[jj], s.MArray[jj]}
			jj++
		}
	}
	return pl
}
func multiPoint2geom(s shp.MultiPoint) geom.T {
	var mp geom.MultiPoint = make([]geom.Point, len(s.Points))
	for i, p := range s.Points {
		mp[i] = geom.Point(p)
	}
	return mp
}
func multiPointM2geom(s shp.MultiPointM) geom.T {
	var mp geom.MultiPointM = make([]geom.PointM, len(s.Points))
	for i, p := range s.Points {
		mp[i] = geom.PointM{p.X, p.Y, s.MArray[i]}
	}
	return mp
}
func multiPointZ2geom(s shp.MultiPointZ) geom.T {
	var mp geom.MultiPointZM = make([]geom.PointZM, len(s.Points))
	for i, p := range s.Points {
		mp[i] = geom.PointZM{p.X, p.Y, s.ZArray[i], s.MArray[i]}
	}
	return mp
}

// Geom2Shp converts a geometry object to a shapefile shape.
func geom2Shp(g geom.T) (shp.Shape, error) {
	if g == nil {
		return &shp.Null{}, nil
	}
	switch t := g.(type) {
	case geom.Point:
		return geom2point(g.(geom.Point)), nil
	case geom.PointM:
		return geom2pointM(g.(geom.PointM)), nil
	case geom.PointZM:
		return geom2pointZ(g.(geom.PointZM)), nil
	case geom.Polygon:
		return geom2polygon(g.(geom.Polygon)), nil
	case geom.PolygonM:
		return geom2polygonM(g.(geom.PolygonM)), nil
	case geom.PolygonZM:
		return geom2polygonZ(g.(geom.PolygonZM)), nil
	case geom.MultiLineString:
		return geom2polyLine(g.(geom.MultiLineString)), nil
	case geom.MultiLineStringM:
		return geom2polyLineM(g.(geom.MultiLineStringM)), nil
	case geom.MultiLineStringZM:
		return geom2polyLineZ(g.(geom.MultiLineStringZM)), nil
	//case t == "MultiPatch": // not yet supported
	case geom.MultiPoint:
		return geom2multiPoint(g.(geom.MultiPoint)), nil
	case geom.MultiPointM:
		return geom2multiPointM(g.(geom.MultiPointM)), nil
	case geom.MultiPointZM:
		return geom2multiPointZ(g.(geom.MultiPointZM)), nil
	default:
		return nil, fmt.Errorf("Unsupported geom type: %v", t)
	}
}

// Functions for converting geom to shp

func geom2point(g geom.Point) shp.Shape {
	p := shp.Point(g)
	return &p
}
func geom2pointM(g geom.PointM) shp.Shape {
	p := shp.PointM(g)
	return &p
}
func geom2pointZ(g geom.PointZM) shp.Shape {
	p := shp.PointZ(g)
	return &p
}
func geom2polygon(g geom.Polygon) shp.Shape {
	parts := make([][]shp.Point, len(g))
	for i, r := range g {
		parts[i] = make([]shp.Point, len(r))
		// switch the winding direction
		for j := len(r) - 1; j >= 0; j-- {
			parts[i][j] = shp.Point(r[j])
		}
	}
	p := shp.Polygon(*shp.NewPolyLine(parts))
	return &p
}
func valrange(a []float64) [2]float64 {
	out := [2]float64{math.Inf(1), math.Inf(-1)}
	for _, val := range a {
		if val < out[0] {
			out[0] = val
		}
		if val > out[1] {
			out[1] = val
		}
	}
	return out
}
func geom2polygonM(g geom.PolygonM) shp.Shape {
	parts := make([][]shp.Point, len(g))
	m := make([]float64, 0)
	jj := 0
	for i, r := range g {
		m = append(m, make([]float64, len(r))...)
		parts[i] = make([]shp.Point, len(r))
		// switch the winding direction
		for j := len(r) - 1; j >= 0; j-- {
			p := r[j]
			parts[i][j] = shp.Point{p.X, p.Y}
			m[jj] = p.M
			jj++
		}
	}
	l := shp.NewPolyLine(parts)
	p := shp.PolygonM{l.Box, l.NumParts, l.NumPoints,
		l.Parts, l.Points, [2]float64{}, nil, valrange(m), m}
	return &p
}
func geom2polygonZ(g geom.PolygonZM) shp.Shape {
	parts := make([][]shp.Point, len(g))
	m := make([]float64, 0)
	z := make([]float64, 0)
	jj := 0
	for i, r := range g {
		m = append(m, make([]float64, len(r))...)
		z = append(z, make([]float64, len(r))...)
		parts[i] = make([]shp.Point, len(r))
		// switch the winding direction
		for j := len(r) - 1; j >= 0; j-- {
			p := r[j]
			parts[i][j] = shp.Point{p.X, p.Y}
			m[jj] = p.M
			z[jj] = p.Z
			jj++
		}
	}
	l := shp.NewPolyLine(parts)
	p := shp.PolygonZ{l.Box, l.NumParts, l.NumPoints,
		l.Parts, l.Points, valrange(z), z, valrange(m), m}
	return &p
}
func geom2polyLine(g geom.MultiLineString) shp.Shape {
	parts := make([][]shp.Point, len(g))
	for i, r := range g {
		parts[i] = make([]shp.Point, len(r))
		for j, l := range r {
			parts[i][j] = shp.Point(l)
		}
	}
	return shp.NewPolyLine(parts)
}
func geom2polyLineM(g geom.MultiLineStringM) shp.Shape {
	parts := make([][]shp.Point, len(g))
	m := make([]float64, 0)
	jj := 0
	for i, r := range g {
		m = append(m, make([]float64, len(r))...)
		parts[i] = make([]shp.Point, len(r))
		for j, l := range r {
			parts[i][j] = shp.Point{l.X, l.Y}
			m[jj] = l.M
			jj++
		}
	}
	l := shp.NewPolyLine(parts)
	p := shp.PolyLineM{l.Box, l.NumParts, l.NumPoints,
		l.Parts, l.Points, valrange(m), m}
	return &p
}
func geom2polyLineZ(g geom.MultiLineStringZM) shp.Shape {
	parts := make([][]shp.Point, len(g))
	m := make([]float64, 0)
	z := make([]float64, 0)
	jj := 0
	for i, r := range g {
		m = append(m, make([]float64, len(r))...)
		z = append(z, make([]float64, len(r))...)
		parts[i] = make([]shp.Point, len(r))
		for j, l := range r {
			parts[i][j] = shp.Point{l.X, l.Y}
			m[jj] = l.M
			z[jj] = l.Z
			jj++
		}
	}
	l := shp.NewPolyLine(parts)
	p := shp.PolyLineZ{l.Box, l.NumParts, l.NumPoints,
		l.Parts, l.Points, valrange(z), z, valrange(m), m}
	return &p
}
func bounds2box(g geom.T) shp.Box {
	b := g.Bounds(nil)
	return shp.Box{b.Min.X, b.Min.Y, b.Max.X, b.Max.Y}
}
func geom2multiPoint(g geom.MultiPoint) shp.Shape {
	mp := new(shp.MultiPoint)
	mp.Box = bounds2box(g)
	mp.NumPoints = int32(len(g))
	mp.Points = make([]shp.Point, len(g))
	for i, p := range g {
		mp.Points[i] = shp.Point(p)
	}
	return mp
}
func geom2multiPointM(g geom.MultiPointM) shp.Shape {
	mp := new(shp.MultiPointM)
	mp.Box = bounds2box(g)
	mp.NumPoints = int32(len(g))
	m := make([]float64, len(g))
	mp.Points = make([]shp.Point, len(g))
	for i, p := range g {
		mp.Points[i] = shp.Point{p.X, p.Y}
		m[i] = p.M
	}
	mp.MArray = m
	mp.MRange = valrange(m)
	return mp
}
func geom2multiPointZ(g geom.MultiPointZM) shp.Shape {
	mp := new(shp.MultiPointZ)
	mp.Box = bounds2box(g)
	mp.NumPoints = int32(len(g))
	m := make([]float64, len(g))
	z := make([]float64, len(g))
	mp.Points = make([]shp.Point, len(g))
	for i, p := range g {
		mp.Points[i] = shp.Point{p.X, p.Y}
		m[i] = p.M
		z[i] = p.Z
	}
	mp.MArray = m
	mp.MRange = valrange(m)
	mp.ZArray = z
	mp.ZRange = valrange(z)
	return mp
}
