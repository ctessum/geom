package osm

import (
	"fmt"

	"github.com/ctessum/geom"
	"github.com/ctessum/geom/op"
	"github.com/qedus/osmpbf"
)

// GeomTags holds a geometry object and the tags that apply to it.
type GeomTags struct {
	geom.Geom
	Tags map[string]string
}

// Geom converts the OSM data to geometry objects, keeping the tag information.
func (o *Data) Geom() ([]*GeomTags, error) {

	items := make([]*GeomTags, 0, len(o.Ways))
	for _, r := range o.Relations {
		if _, ok := o.dependentRelations[r.ID]; !ok {
			g, err := relationToGeom(r, o.Relations, o.Ways, o.Nodes)
			if err != nil {
				return nil, err
			}
			items = append(items, &GeomTags{
				Geom: g,
				Tags: r.Tags,
			})
		}
	}
	for _, w := range o.Ways {
		if _, ok := o.dependentWays[w.ID]; !ok {
			items = append(items, &GeomTags{
				Geom: wayToGeom(w, o.Nodes),
				Tags: w.Tags,
			})
		}
	}
	for _, n := range o.Nodes {
		if _, ok := o.dependentNodes[n.ID]; !ok {
			items = append(items, &GeomTags{
				Geom: nodeToPoint(n),
				Tags: n.Tags,
			})
		}
	}
	return items, nil
}

func nodeToPoint(n *osmpbf.Node) geom.Point {
	return geom.Point{X: n.Lon, Y: n.Lat}
}

func wayToGeom(way *osmpbf.Way, nodes map[int64]*osmpbf.Node) geom.Geom {
	if wayIsClosed(way) {
		return wayToPolygon(way, nodes)
	}
	return wayToLineString(way, nodes)
}

// wayIsClosed determines whether a way represents a polygon.
func wayIsClosed(way *osmpbf.Way) bool {
	return way.NodeIDs[0] == way.NodeIDs[len(way.NodeIDs)-1]
}

// wayToPolygon converts a way to a polygon
func wayToPolygon(way *osmpbf.Way, nodes map[int64]*osmpbf.Node) geom.Polygon {
	p := make(geom.Polygon, 1)
	for _, nid := range way.NodeIDs {
		p[0] = append(p[0], nodeToPoint(nodes[nid]))
	}
	return p
}

// wayToLineString converts a way to a LineString
func wayToLineString(way *osmpbf.Way, nodes map[int64]*osmpbf.Node) geom.LineString {
	var p geom.LineString
	for _, nid := range way.NodeIDs {
		p = append(p, nodeToPoint(nodes[nid]))
	}
	return p
}

// relationToGeom converts a relation to a geometry object.
func relationToGeom(relation *osmpbf.Relation,
	relations map[int64]*osmpbf.Relation, ways map[int64]*osmpbf.Way,
	nodes map[int64]*osmpbf.Node) (geom.Geom, error) {

	var nNodes, nLines, nPolygons int
	for _, m := range relation.Members {
		switch m.Type {
		case osmpbf.WayType:
			if wayIsClosed(ways[m.ID]) {
				nPolygons++
			} else {
				nLines++
			}
		case osmpbf.NodeType:
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
	return relationToGeometryCollection(relation, relations, ways, nodes)
}

// relationToMultiPoint converts a relation to a MultiPoint
func relationToMultiPoint(relation *osmpbf.Relation,
	nodes map[int64]*osmpbf.Node) geom.MultiPoint {

	p := make(geom.MultiPoint, len(relation.Members))
	for i, m := range relation.Members {
		switch m.Type {
		case osmpbf.NodeType:
			p[i] = nodeToPoint(nodes[m.ID])
		default:
			panic(fmt.Errorf("unsupported relation type %T", m.Type))
		}
	}
	return p
}

// relationToPolygon converts a relation to a polygon
func relationToPolygon(relation *osmpbf.Relation, ways map[int64]*osmpbf.Way,
	nodes map[int64]*osmpbf.Node) (geom.Polygon, error) {
	var p geom.Polygon
	for _, m := range relation.Members {
		switch m.Type {
		case osmpbf.WayType:
			p = append(p, wayToPolygon(ways[m.ID], nodes)[0])
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
func relationToMultiLineString(relation *osmpbf.Relation, ways map[int64]*osmpbf.Way,
	nodes map[int64]*osmpbf.Node) geom.MultiLineString {
	var p geom.MultiLineString
	for _, m := range relation.Members {
		switch m.Type {
		case osmpbf.WayType:
			p = append(p, wayToLineString(ways[m.ID], nodes))
		default:
			panic(fmt.Errorf("unsupported relation type %T", m.Type))
		}
	}
	return p
}

func relationToGeometryCollection(relation *osmpbf.Relation,
	relations map[int64]*osmpbf.Relation, ways map[int64]*osmpbf.Way,
	nodes map[int64]*osmpbf.Node) (geom.Geom, error) {

	p := make(geom.GeometryCollection, len(relation.Members))
	for i, m := range relation.Members {
		switch m.Type {
		case osmpbf.WayType:
			p[i] = wayToGeom(ways[m.ID], nodes)
		case osmpbf.NodeType:
			p[i] = nodeToPoint(nodes[m.ID])
		case osmpbf.RelationType:
			var err error
			p[i], err = relationToGeom(relations[m.ID], relations, ways, nodes)
			if err != nil {
				return nil, err
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
	Collection
)

// DominantType returns the most frequently occurring type among the
// given features.
func DominantType(gt []*GeomTags) (GeomType, error) {
	var points, lines, polys, collections int
	for _, g := range gt {
		switch g.Geom.(type) {
		case geom.Point:
			points++
		case geom.MultiPoint:
			points++
		case geom.LineString, geom.MultiLineString:
			lines++
		case geom.Polygon:
			polys++
		case geom.GeometryCollection:
			collections++
			continue
		default:
			return -1, fmt.Errorf("invalid geometry type %T", g.Geom)
		}
	}
	if points >= lines && points >= polys && points >= collections {
		return Point, nil
	}
	if lines > points && lines >= polys && lines >= collections {
		return Line, nil
	}
	if polys > points && polys > lines && polys >= collections {
		return Poly, nil
	}
	return Collection, nil
}
