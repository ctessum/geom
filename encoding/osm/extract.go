// Package osm extracts and manipulates OpenStreetMap (OSM) data. Refer to
// openstreetmap.org for more information about OSM data.
package osm

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

// Extract extracts OpenStreetMap data from osm.pbf file rs.
// keep determines which records are included in the output.
func Extract(ctx context.Context, rs io.ReadSeeker, keep KeepFunc) (*Data, error) {
	o := &Data{
		Nodes:     make(map[osm.NodeID]*osm.Node),
		Ways:      make(map[osm.WayID]*osm.Way),
		Relations: make(map[osm.RelationID]*osm.Relation),

		dependentNodes:     make(map[osm.NodeID]empty),
		dependentWays:      make(map[osm.WayID]empty),
		dependentRelations: make(map[osm.RelationID]empty),
	}

	needAnotherPass := true
	//passI := 0
	for needAnotherPass {
		needAnotherPass = false
		if _, err := rs.Seek(0, 0); err != nil {
			return nil, err
		}
		scanner := osmpbf.New(ctx, rs, runtime.GOMAXPROCS(-1))
		for scanner.Scan() {
			obj := scanner.Object()
			switch objType := obj.(type) {
			case *osm.Node:
				o.processNode(obj.(*osm.Node), keep)
			case *osm.Way:
				if o.processWay(obj.(*osm.Way), keep) {
					needAnotherPass = true
				}
			case *osm.Relation:
				if o.processRelation(obj.(*osm.Relation), keep) {
					needAnotherPass = true
				}
			default:
				return nil, fmt.Errorf("unknown type %T\n", objType)
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		if err := scanner.Close(); err != nil {
			return nil, err
		}
		//passI++
		//log.Printf("pass %d: %d nodes, %d ways, %d relations, "+
		//	"%d dependent nodes, %d dependent ways, %d dependent relations", passI,
		//	len(o.Nodes), len(o.Ways), len(o.Relations),
		//	len(o.dependentNodes), len(o.dependentWays), len(o.dependentRelations))
	}
	return o, nil
}

// ExtractTag extracts OpenStreetMap data with the given tag set to one of the
// given values.
func ExtractTag(rs io.ReadSeeker, tag string, values ...string) (*Data, error) {
	return Extract(context.Background(), rs, KeepTags(map[string][]string{tag: values}))
}

// ObjectType specifies the valid OpenStreetMap types.
type ObjectType int

const (
	// Node is an OpenStreetMap node.
	Node ObjectType = iota
	// Way can be either open or closed.
	Way
	// ClosedWay is an OpenStreetMap way that is closed (i.e., a polygon).
	ClosedWay
	// OpenWay is an OpenStreetMap way that is open (i.e., a line string).
	OpenWay
	// Relation is an OpenStreetMap relation.
	Relation
)

type empty struct{}

// Data holds OpenStreetMap data and relationships.
type Data struct {
	Nodes     map[osm.NodeID]*osm.Node
	Ways      map[osm.WayID]*osm.Way
	Relations map[osm.RelationID]*osm.Relation

	// These list the objects that are dependent on other objects,
	// and the objects that they are dependent on.
	dependentNodes     map[osm.NodeID]empty
	dependentWays      map[osm.WayID]empty
	dependentRelations map[osm.RelationID]empty
}

// Filter returns a copy of the receiver where only objects
// selected by keep are retained.
func (o *Data) Filter(keep KeepFunc) *Data {
	out := &Data{
		Nodes:     make(map[osm.NodeID]*osm.Node),
		Ways:      make(map[osm.WayID]*osm.Way),
		Relations: make(map[osm.RelationID]*osm.Relation),

		dependentNodes:     make(map[osm.NodeID]empty),
		dependentWays:      make(map[osm.WayID]empty),
		dependentRelations: make(map[osm.RelationID]empty),
	}

	needAnotherPass := true
	for needAnotherPass {
		needAnotherPass = false
		for _, n := range o.Nodes {
			out.processNode(n, keep)
		}
		for _, w := range o.Ways {
			if out.processWay(w, keep) {
				needAnotherPass = true
			}
		}
		for _, r := range o.Relations {
			if out.processRelation(r, keep) {
				needAnotherPass = true
			}
		}
	}
	return out
}

// hasTag checks if tags[t] is one of the values in wantTags. If len(v) == 0,
// the function will return true is tags[t] has any value.
func hasTag(tags map[string][]string, wantTags map[string][]string) bool {
	for wantTagKey, wantTagValues := range wantTags {
		tagValues, ok := tags[wantTagKey]
		if !ok {
			continue
		}
		if len(tagValues) == 0 {
			return true
		}
		for _, wantTagValue := range wantTagValues {
			for _, tagValue := range tagValues {
				if tagValue == wantTagValue {
					return true
				}
			}
		}
	}
	return false
}

func (o *Data) hasNeedNode(id osm.NodeID) (has, need bool) {
	if _, ok := o.Nodes[id]; ok {
		has = true
		return
	}
	if _, ok := o.dependentNodes[id]; ok {
		need = true
	}
	return
}

func (o *Data) hasNeedWay(id osm.WayID) (has, need bool) {
	if _, ok := o.Ways[id]; ok {
		has = true
		return
	}
	if _, ok := o.dependentWays[id]; ok {
		need = true
	}
	return
}

func (o *Data) hasNeedRelation(id osm.RelationID) (has, need bool) {
	if _, ok := o.Relations[id]; ok {
		has = true
		return
	}
	if _, ok := o.dependentRelations[id]; ok {
		need = true
	}
	return
}

// If the node has the tag we want, add it to the list.
func (o *Data) processNode(n *osm.Node, keep KeepFunc) {
	hasNode, needNode := o.hasNeedNode(n.ID)
	if hasNode {
		return
	}
	if keep(o, n) || needNode {
		o.Nodes[n.ID] = n
	}
}

// If the way has the tag we want or if we've determined that it's
// part of a relation that we want, store the way and the IDs of its dependent nodes.
func (o *Data) processWay(w *osm.Way, keep KeepFunc) (anotherPass bool) {
	hasWay, needWay := o.hasNeedWay(w.ID)
	if hasWay {
		return
	}
	if keep(o, w) || needWay {
		o.Ways[w.ID] = w
		for _, n := range w.Nodes {
			if _, needNode := o.hasNeedNode(n.ID); !needNode {
				o.dependentNodes[n.ID] = empty{}
				anotherPass = true
			}
		}
	}
	return
}

// If the relation has the tag we want or if we've determined that it's
// part of a different relation that we want, store the IDs of its
// members and set the flag for another pass through the file to
// get the IDs for the dependent nodes, ways and other relations in the relation.
func (o *Data) processRelation(r *osm.Relation, keep KeepFunc) (anotherPass bool) {
	hasRelation, needRelation := o.hasNeedRelation(r.ID)
	if hasRelation {
		return
	}
	if keep(o, r) || needRelation {
		o.Relations[r.ID] = r
		for _, m := range r.Members {
			switch m.Type {
			case osm.TypeNode:
				if _, needNode := o.hasNeedNode(osm.NodeID(m.Ref)); !needNode {
					o.dependentNodes[osm.NodeID(m.Ref)] = empty{}
					anotherPass = true
				}
			case osm.TypeWay:
				if _, needWay := o.hasNeedWay(osm.WayID(m.Ref)); !needWay {
					o.dependentWays[osm.WayID(m.Ref)] = empty{}
					anotherPass = true
				}
			case osm.TypeRelation:
				if _, needR := o.hasNeedRelation(osm.RelationID(m.Ref)); !needR {
					o.dependentRelations[osm.RelationID(m.Ref)] = empty{}
					anotherPass = true
				}
			default:
				panic(fmt.Errorf("unknown member type %v", m.Type))
			}
		}
	}
	return
}
