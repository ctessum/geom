package osm

import (
	"fmt"

	"github.com/qedus/osmpbf"
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
		for _, n := range w.NodeIDs {
			if x, ok := o.Nodes[n]; !ok {
				return fmt.Errorf("node %d is referenced by way %d but does not exist",
					n, w.ID)
			} else if x == nil {
				return fmt.Errorf("node %d is referenced by way %d but is nil",
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
			case osmpbf.NodeType:
				if x, ok := o.Nodes[m.ID]; !ok {
					return fmt.Errorf("node %d is referenced by relation %d but does not exist",
						m.ID, r.ID)
				} else if x == nil {
					return fmt.Errorf("node %d is referenced by relation %d but is nil",
						m.ID, r.ID)
				}
			case osmpbf.WayType:
				if x, ok := o.Ways[m.ID]; !ok {
					return fmt.Errorf("way %d is referenced by relation %d but does not exist",
						m.ID, r.ID)
				} else if x == nil {
					return fmt.Errorf("way %d is referenced by relation %d but is nil",
						m.ID, r.ID)
				}
			case osmpbf.RelationType:
				if x, ok := o.Relations[m.ID]; !ok {
					return fmt.Errorf("relation %d is referenced by relation %d but does not exist",
						m.ID, r.ID)
				} else if x == nil {
					return fmt.Errorf("relation %d is referenced by relation %d but is nil",
						m.ID, r.ID)
				}
			default:
				return fmt.Errorf("unknown member type %v in relation %d", m.Type, i)
			}
		}
	}
	return nil
}
