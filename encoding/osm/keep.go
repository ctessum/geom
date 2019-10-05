package osm

import (
	"fmt"

	"github.com/ctessum/geom"
	"github.com/qedus/osmpbf"
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
		case *osmpbf.Node:
			return hasTag(object.(*osmpbf.Node).Tags, tags)
		case *osmpbf.Way:
			return hasTag(object.(*osmpbf.Way).Tags, tags)
		case *osmpbf.Relation:
			return hasTag(object.(*osmpbf.Relation).Tags, tags)
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
		case *osmpbf.Node:
			// For nodes, keep anything that is within b.
			n := object.(*osmpbf.Node)
			return b.Overlaps(geom.Point{X: n.Lon, Y: n.Lat}.Bounds())
		case *osmpbf.Way:
			// For ways, keep anything that requires a node that we're already keeping.
			w := object.(*osmpbf.Way)
			for _, n := range w.NodeIDs {
				if has, _ := o.hasNeedNode(n); has {
					return true
				}
			}
		case *osmpbf.Relation:
			// For relations, keep anything that requires a node, way or relation that
			// we're already keeping.
			r := object.(*osmpbf.Relation)
			for _, m := range r.Members {
				switch m.Type {
				case osmpbf.NodeType:
					if has, _ := o.hasNeedNode(m.ID); has {
						return true
					}
				case osmpbf.WayType:
					if has, _ := o.hasNeedWay(m.ID); has {
						return true
					}
				case osmpbf.RelationType:
					if has, _ := o.hasNeedRelation(m.ID); has {
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
