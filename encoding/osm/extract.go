// Package osm extracts and manipulates OpenStreetMap (OSM) data. Refer to
// openstreetmap.org for more information about OSM data.
package osm

import (
	"fmt"
	"io"
	"runtime"

	"github.com/qedus/osmpbf"
)

// ExtractTag extracts OpenStreetMap data with the given tag set to one of the
// given values.
func ExtractTag(rs io.ReadSeeker, tag string, values ...string) (*Data, error) {

	o := &Data{
		Nodes:     make(map[int64]*osmpbf.Node),
		Ways:      make(map[int64]*osmpbf.Way),
		Relations: make(map[int64]*osmpbf.Relation),

		dependentNodes:     make(map[int64]empty),
		dependentWays:      make(map[int64]empty),
		dependentRelations: make(map[int64]empty),
	}

	needAnotherPass := true
	//passI := 0
	for needAnotherPass {
		needAnotherPass = false
		if _, err := rs.Seek(0, 0); err != nil {
			return nil, err
		}
		data := osmpbf.NewDecoder(rs)
		if err := data.Start(runtime.GOMAXPROCS(-1)); err != nil {
			return nil, err
		}
		for {
			var v interface{}
			var err error
			if v, err = data.Decode(); err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			switch vtype := v.(type) {
			case *osmpbf.Node:
				o.processNode(v.(*osmpbf.Node), tag, values)
			case *osmpbf.Way:
				if o.processWay(v.(*osmpbf.Way), tag, values) {
					needAnotherPass = true
				}
			case *osmpbf.Relation:
				if o.processRelation(v.(*osmpbf.Relation), tag, values) {
					needAnotherPass = true
				}
			default:
				return nil, fmt.Errorf("unknown type %T\n", vtype)
			}
		}
		//passI++
		//log.Printf("pass %d: %d nodes, %d ways, %d relations, "+
		//	"%d dependent nodes, %d dependent ways, %d dependent relations", passI,
		//	len(o.Nodes), len(o.Ways), len(o.Relations),
		//	len(o.dependentNodes), len(o.dependentWays), len(o.dependentRelations))
	}
	return o, nil
}

type empty struct{}

// Data holds OpenStreetMap data and relationships.
type Data struct {
	Nodes     map[int64]*osmpbf.Node
	Ways      map[int64]*osmpbf.Way
	Relations map[int64]*osmpbf.Relation

	dependentNodes     map[int64]empty
	dependentWays      map[int64]empty
	dependentRelations map[int64]empty
}

// hasTag checks if tags[t] is one of the values in v. If len(v) == 0,
// the function will return true is tags[t] has any value.
func hasTag(tags map[string]string, t string, v []string) bool {
	vv, ok := tags[t]
	if !ok {
		return false
	}
	if len(v) == 0 {
		return true
	}
	for _, vvv := range v {
		if vv == vvv {
			return true
		}
	}
	return false
}

func (o *Data) hasNeedNode(id int64) (has, need bool) {
	if _, ok := o.Nodes[id]; ok {
		has = true
		return
	}
	if _, ok := o.dependentNodes[id]; ok {
		need = true
	}
	return
}

func (o *Data) hasNeedWay(id int64) (has, need bool) {
	if _, ok := o.Ways[id]; ok {
		has = true
		return
	}
	if _, ok := o.dependentWays[id]; ok {
		need = true
	}
	return
}

func (o *Data) hasNeedRelation(id int64) (has, need bool) {
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
func (o *Data) processNode(n *osmpbf.Node, tag string, vals []string) {
	hasNode, needNode := o.hasNeedNode(n.ID)
	if hasNode {
		return
	}
	if hasTag(n.Tags, tag, vals) || needNode {
		o.Nodes[n.ID] = n
	}
}

// If the way has the tag we want or if we've determined that it's
// part of a relation that we want, store the way and the IDs of its dependent nodes.
func (o *Data) processWay(w *osmpbf.Way, tag string, vals []string) (anotherPass bool) {
	hasWay, needWay := o.hasNeedWay(w.ID)
	if hasWay {
		return
	}
	if hasTag(w.Tags, tag, vals) || needWay {
		o.Ways[w.ID] = w
		for _, n := range w.NodeIDs {
			if _, needNode := o.hasNeedNode(n); !needNode {
				o.dependentNodes[n] = empty{}
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
func (o *Data) processRelation(r *osmpbf.Relation,
	tag string, vals []string) (anotherPass bool) {

	hasRelation, needRelation := o.hasNeedRelation(r.ID)
	if hasRelation {
		return
	}
	if hasTag(r.Tags, tag, vals) || needRelation {
		o.Relations[r.ID] = r
		for _, m := range r.Members {
			switch m.Type {
			case osmpbf.NodeType:
				if _, needNode := o.hasNeedNode(m.ID); !needNode {
					o.dependentNodes[m.ID] = empty{}
					anotherPass = true
				}
			case osmpbf.WayType:
				if _, needWay := o.hasNeedWay(m.ID); !needWay {
					o.dependentWays[m.ID] = empty{}
					anotherPass = true
				}
			case osmpbf.RelationType:
				if _, needR := o.hasNeedRelation(m.ID); !needR {
					o.dependentRelations[m.ID] = empty{}
					anotherPass = true
				}
			default:
				panic(fmt.Errorf("unknown member type %v", m.Type))
			}
		}
	}
	return
}
