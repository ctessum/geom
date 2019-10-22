package osm

import (
	"fmt"
	"math"

	"github.com/ctessum/geom"
	"github.com/ctessum/geom/op"
	"github.com/paulmach/osm"
)

// GeomTags holds a geometry object and the tags that apply to it.
type GeomTags struct {
	geom.Geom
	Tags map[string][]string
}

func tagsToMap(tags osm.Tags) map[string][]string {
	o := make(map[string][]string)
	for _, t := range tags {
		o[t.Key] = append(o[t.Key], t.Value)
	}
	return o
}

// Geom converts the OSM data to geometry objects, keeping the tag information.
func (o *Data) Geom() ([]*GeomTags, error) {
	items := make([]*GeomTags, 0, len(o.Ways))
	for _, r := range o.Relations {
		if _, ok := o.dependentRelations[r.ID]; !ok {
			if r != nil {
				g, err := relationToGeom(r, o.Relations, o.Ways, o.Nodes, make(map[osm.RelationID]struct{}))
				if err != nil {
					return nil, err
				}
				items = append(items, &GeomTags{
					Geom: g,
					Tags: tagsToMap(r.Tags),
				})
			}
		}
	}
	for _, w := range o.Ways {
		if _, ok := o.dependentWays[w.ID]; !ok {
			if w != nil && len(w.Nodes) > 0 {
				items = append(items, &GeomTags{
					Geom: wayToGeom(w, o.Nodes),
					Tags: tagsToMap(w.Tags),
				})
			}
		}
	}
	for _, n := range o.Nodes {
		if _, ok := o.dependentNodes[n.ID]; !ok {
			p, ok := nodeToPoint(n)
			if ok {
				items = append(items, &GeomTags{
					Geom: p,
					Tags: tagsToMap(n.Tags),
				})
			}
		}
	}
	return items, nil
}

func nodeToPoint(n *osm.Node) (geom.Point, bool) {
	if n == nil {
		return geom.Point{X: math.NaN(), Y: math.NaN()}, false
	}
	return geom.Point{X: n.Lon, Y: n.Lat}, true
}

func wayToGeom(way *osm.Way, nodes map[osm.NodeID]*osm.Node) geom.Geom {
	if wayIsClosed(way) {
		return wayToPolygon(way, nodes)
	}
	return wayToLineString(way, nodes)
}

// wayIsClosed determines whether a way represents a polygon.
func wayIsClosed(way *osm.Way) bool {
	return way.Nodes[0].ID == way.Nodes[len(way.Nodes)-1].ID
}

// wayToPolygon converts a way to a polygon
func wayToPolygon(way *osm.Way, nodes map[osm.NodeID]*osm.Node) geom.Polygon {
	p := make(geom.Polygon, 1)
	for _, n := range way.Nodes {
		point, ok := nodeToPoint(nodes[n.ID])
		if ok {
			p[0] = append(p[0], point)
		}
	}
	return p
}

// wayToLineString converts a way to a LineString
func wayToLineString(way *osm.Way, nodes map[osm.NodeID]*osm.Node) geom.LineString {
	var p geom.LineString
	for _, n := range way.Nodes {
		point, ok := nodeToPoint(nodes[n.ID])
		if ok {
			p = append(p, point)
		}
	}
	return p
}

// relationToGeom converts a relation to a geometry object.
func relationToGeom(relation *osm.Relation,
	relations map[osm.RelationID]*osm.Relation, ways map[osm.WayID]*osm.Way,
	nodes map[osm.NodeID]*osm.Node, idStack map[osm.RelationID]struct{}) (geom.Geom, error) {

	var nNodes, nLines, nPolygons int
	for _, m := range relation.Members {
		switch m.Type {
		case osm.TypeWay:
			if w, ok := ways[osm.WayID(m.Ref)]; ok && len(w.Nodes) > 0 && wayIsClosed(w) {
				nPolygons++
			} else {
				nLines++
			}
		case osm.TypeNode:
			nNodes++
		}
	}
	if nPolygons == len(relation.Members) {
		return relationToPolygon(relation, ways, nodes)
	}
	if nLines == len(relation.Members) {
		return relationToMultiLineString(relation, ways, nodes), nil
	}
	if nNodes == len(relation.Members) {
		return relationToMultiPoint(relation, nodes), nil
	}
	return relationToGeometryCollection(relation, relations, ways, nodes, idStack)
}

// relationToMultiPoint converts a relation to a MultiPoint
func relationToMultiPoint(relation *osm.Relation,
	nodes map[osm.NodeID]*osm.Node) geom.MultiPoint {

	p := make(geom.MultiPoint, 0, len(relation.Members))
	for _, m := range relation.Members {
		switch m.Type {
		case osm.TypeNode:
			point, ok := nodeToPoint(nodes[osm.NodeID(m.Ref)])
			if ok {
				p = append(p, point)
			}
		default:
			panic(fmt.Errorf("unsupported relation type %T", m.Type))
		}
	}
	return p
}

// relationToPolygon converts a relation to a polygon
func relationToPolygon(relation *osm.Relation, ways map[osm.WayID]*osm.Way,
	nodes map[osm.NodeID]*osm.Node) (geom.Polygon, error) {
	var p geom.Polygon
	for _, m := range relation.Members {
		switch m.Type {
		case osm.TypeWay:
			if w := ways[osm.WayID(m.Ref)]; w != nil {
				p = append(p, wayToPolygon(w, nodes)[0])
			}
		default:
			panic(fmt.Errorf("unsupported relation type %T", m.Type))
		}
	}
	if err := op.FixOrientation(p); err != nil {
		return nil, err
	}
	return p, nil
}

// relationToMultiLineString converts a relation to a MultiLineString,
// deleting its contained elements from 'ways' and 'nodes'.
func relationToMultiLineString(relation *osm.Relation, ways map[osm.WayID]*osm.Way,
	nodes map[osm.NodeID]*osm.Node) geom.MultiLineString {
	var p geom.MultiLineString
	for _, m := range relation.Members {
		switch m.Type {
		case osm.TypeWay:
			if w := ways[osm.WayID(m.Ref)]; w != nil {
				p = append(p, wayToLineString(w, nodes))
			}
		default:
			panic(fmt.Errorf("unsupported relation type %T", m.Type))
		}
	}
	return p
}

func relationToGeometryCollection(relation *osm.Relation,
	relations map[osm.RelationID]*osm.Relation, ways map[osm.WayID]*osm.Way,
	nodes map[osm.NodeID]*osm.Node, idStack map[osm.RelationID]struct{}) (geom.Geom, error) {

	p := make(geom.GeometryCollection, 0, len(relation.Members))
	for _, m := range relation.Members {
		switch m.Type {
		case osm.TypeWay:
			way := ways[osm.WayID(m.Ref)]
			if way != nil && len(way.Nodes) > 0 {
				p = append(p, wayToGeom(way, nodes))
			}
		case osm.TypeNode:
			point, ok := nodeToPoint(nodes[osm.NodeID(m.Ref)])
			if ok {
				p = append(p, point)
			}
		case osm.TypeRelation:
			id := osm.RelationID(m.Ref)
			if _, ok := idStack[id]; ok {
				continue // Skip self-references, which cause infinite recursion.
			}
			idStack[id] = struct{}{}
			if r := relations[id]; r != nil {
				g, err := relationToGeom(r, relations, ways, nodes, idStack)
				if err != nil {
					return nil, err
				}
				p = append(p, g)
			}
		default:
			panic(fmt.Errorf("unsupported relation type %T", m.Type))
		}
	}
	return p, nil
}

type GeomType int

const (
	Point GeomType = iota
	Line
	Poly
)

func countGeomTagsTypes(gt []*GeomTags) (points, lines, polys int, err error) {
	for _, g := range gt {
		pt, l, pl, err := countGeometryTypes(g.Geom)
		if err != nil {
			return -1, -1, -1, err
		}
		points += pt
		lines += l
		polys += pl
	}
	return
}

func countGeometryTypes(g geom.Geom) (points, lines, polys int, err error) {
	switch g.(type) {
	case geom.Point:
		points++
	case geom.MultiPoint:
		points++
	case geom.LineString, geom.MultiLineString:
		lines++
	case geom.Polygon:
		polys++
	case geom.GeometryCollection:
		for _, f := range g.(geom.GeometryCollection) {
			pt, l, pl, err := countGeometryTypes(f)
			if err != nil {
				return -1, -1, -1, err
			}
			points += pt
			lines += l
			polys += pl
		}
	default:
		return -1, -1, -1, fmt.Errorf("invalid geometry type %T", g)
	}
	return
}

// DominantType returns the most frequently occurring type among the
// given features.
func DominantType(gt []*GeomTags) (GeomType, error) {
	points, lines, polys, err := countGeomTagsTypes(gt)
	if err != nil {
		return -1, err
	}
	if points >= lines && points >= polys {
		return Point, nil
	}
	if lines > points && lines >= polys {
		return Line, nil
	}
	if polys > points && polys > lines {
		return Poly, nil
	}
	panic(fmt.Errorf("no dominant type: %d, %d, %d", points, lines, polys))
}
