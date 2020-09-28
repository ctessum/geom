package osm

import (
	"fmt"

	"github.com/paulmach/osm"
)

// Check checks OSM data to ensure that all necessary components
// are present.
func (o *Data) Check() error {
	for i, n := range o.Nodes {
		if n == nil {
			return fmt.Errorf("node %d is nil", i)
		}
	}
	for i, w := range o.Ways {
		if w == nil {
			return fmt.Errorf("way %d is nil", i)
		}
		for _, n := range w.Nodes {
			if x, ok := o.Nodes[n]; !ok {
				return fmt.Errorf("node %v is referenced by way %v but does not exist",
					n, w.ID)
			} else if x == nil {
				return fmt.Errorf("node %v is referenced by way %v but is nil",
					n, w.ID)
			}
		}
	}
	for i, r := range o.Relations {
		if r == nil {
			return fmt.Errorf("relation %d is nil", i)
		}
		for _, m := range r.Members {
			switch m.Type {
			case osm.TypeNode:
				if x, ok := o.Nodes[osm.NodeID(m.Ref)]; !ok {
					return fmt.Errorf("node %d is referenced by relation %d but does not exist",
						m.Ref, r.ID)
				} else if x == nil {
					return fmt.Errorf("node %d is referenced by relation %d but is nil",
						m.Ref, r.ID)
				}
			case osm.TypeWay:
				if x, ok := o.Ways[osm.WayID(m.Ref)]; !ok {
					return fmt.Errorf("way %d is referenced by relation %d but does not exist",
						m.Ref, r.ID)
				} else if x == nil {
					return fmt.Errorf("way %d is referenced by relation %d but is nil",
						m.Ref, r.ID)
				}
			case osm.TypeRelation:
				if x, ok := o.Relations[osm.RelationID(m.Ref)]; !ok {
					return fmt.Errorf("relation %d is referenced by relation %d but does not exist",
						m.Ref, r.ID)
				} else if x == nil {
					return fmt.Errorf("relation %d is referenced by relation %d but is nil",
						m.Ref, r.ID)
				}
			default:
				return fmt.Errorf("unknown member type %v in relation %d", m.Type, i)
			}
		}
	}
	return nil
}
