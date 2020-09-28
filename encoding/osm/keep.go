package osm

import (
	"fmt"

	"github.com/ctessum/geom"
	"github.com/paulmach/osm"
)

// KeepFunc is a function that determines whether an OSM object should
// be included in the output. The object may be either *osmpbf.Node,
// *osmpbf.Way, or *osmpbf.Relation
type KeepFunc func(d *Data, object interface{}) bool

// KeepTags keeps OSM objects that contain the given tag key
// with at least one of the given tag values, where
// the keys and values correspond to the keys and values
// of the 'tags' input. If no
// tag valuse are given, then all object with the given
// key will be kept.
func KeepTags(tags map[string][]string) KeepFunc {
	return func(_ *Data, object interface{}) bool {
		switch object.(type) {
		case *osm.Node:
			return hasTag(object.(*osm.Node).Tags, tags)
		case *Node:
			return hasTag(object.(*Node).Tags, tags)
		case *osm.Way:
			return hasTag(object.(*osm.Way).Tags, tags)
		case *Way:
			return hasTag(object.(*Way).Tags, tags)
		case *osm.Relation:
			return hasTag(object.(*osm.Relation).Tags, tags)
		case *Relation:
			return hasTag(object.(*Relation).Tags, tags)
		default:
			panic(fmt.Errorf("osm: invalid object type %T", object))
		}
	}
}

// KeepBounds keeps OSM objects that overlap with b.
// Using KeepBounds in combination with other KeepFuncs may
// result in unexpected results.
func KeepBounds(b *geom.Bounds) KeepFunc {
	return func(o *Data, object interface{}) bool {
		switch object.(type) {
		case *osm.Node:
			// For nodes, keep anything that is within b.
			n := object.(*osm.Node)
			return b.Overlaps(geom.Point{X: n.Lon, Y: n.Lat}.Bounds())
		case *osm.Way:
			// For ways, keep anything that requires a node that we're already keeping.
			w := object.(*osm.Way)
			for _, n := range w.Nodes {
				if has, _ := o.hasNeedNode(n.ID); has {
					return true
				}
			}
		case *osm.Relation:
			// For relations, keep anything that requires a node, way or relation that
			// we're already keeping.
			r := object.(*osm.Relation)
			for _, m := range r.Members {
				switch m.Type {
				case osm.TypeNode:
					if has, _ := o.hasNeedNode(osm.NodeID(m.Ref)); has {
						return true
					}
				case osm.TypeWay:
					if has, _ := o.hasNeedWay(osm.WayID(m.Ref)); has {
						return true
					}
				case osm.TypeRelation:
					if has, _ := o.hasNeedRelation(osm.RelationID(m.Ref)); has {
						return true
					}
				default:
					panic(fmt.Errorf("unknown member type %v", m.Type))
				}
			}
		default:
			panic(fmt.Errorf("osm: invalid object type %T", object))
		}
		return false
	}
}

// KeepAll specifies that all objects should be kept.
func KeepAll() KeepFunc {
	return func(_ *Data, _ interface{}) bool { return true }
}
