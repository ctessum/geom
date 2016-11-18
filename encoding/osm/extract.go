// Package osm extracts and manipulates OpenStreetMap (OSM) data. Refer to
// openstreetmap.org for more information about OSM data.
package osm

import (
	"fmt"
	"io"
	"runtime"
	"sort"

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

// ObjectType specifies the valid OpenStreetMap types.
type ObjectType int

const (
	// Node is an OpenStreetMap node.
	Node ObjectType = iota
	// ClosedWay is an OpenStreetMap way that is closed (i.e., a polygon).
	ClosedWay
	// OpenWay is an OpenStreetMap way that is open (i.e., a line string).
	OpenWay
	// Relation is an OpenStreetMap relation.
	Relation
)

// Tags holds information about the tags that are in a database.
type Tags []*TagCount

// Len returns the length of the receiver to implement the sort.Sort interface.
func (t *Tags) Len() int { return len(*t) }

// Less returns whether item i is less than item j
// to implement the sort.Sort interface.
func (t *Tags) Less(i, j int) bool {
	tt := *t
	if tt[i].TotalCount < tt[j].TotalCount {
		return true
	}
	if tt[i].TotalCount > tt[j].TotalCount {
		return false
	}
	if tt[i].Key < tt[j].Key {
		return true
	}
	if tt[i].Key > tt[j].Key {
		return false
	}
	return tt[i].Value < tt[j].Value
}

// Table creates a table of the information held by the receiver.
func (t *Tags) Table() [][]string {
	o := make([][]string, len(*t)+1)
	o[0] = []string{"Key", "Value", "Total", "Node", "Closed way", "Open way", "Relation"}
	for i, tt := range *t {
		o[i+1] = []string{
			tt.Key,
			tt.Value,
			fmt.Sprintf("%d", tt.TotalCount),
			fmt.Sprintf("%d", tt.ObjectCount[Node]),
			fmt.Sprintf("%d", tt.ObjectCount[ClosedWay]),
			fmt.Sprintf("%d", tt.ObjectCount[OpenWay]),
			fmt.Sprintf("%d", tt.ObjectCount[Relation]),
		}
	}
	return o
}

// Filter applies function f to all records in the receiver
// and returns a copy of the receiver that only contains
// the records for which f returns true.
func (t *Tags) Filter(f func(*TagCount) bool) *Tags {
	var o Tags
	for _, tt := range *t {
		if f(tt) {
			o = append(o, tt)
		}
	}
	return &o
}

// Swap swaps elements i and j
// to implement the sort.Sort interface.
func (t *Tags) Swap(i, j int) { (*t)[i], (*t)[j] = (*t)[j], (*t)[i] }

// TagCount hold information about the number of instances of
// the specified tag in a database.
type TagCount struct {
	Key, Value  string
	ObjectCount map[ObjectType]int
	TotalCount  int
}

// DominantType returns the most frequently occuring ObjectType for
// this tag.
func (t *TagCount) DominantType() ObjectType {
	v := 0
	result := ObjectType(-1)
	for typ, vv := range t.ObjectCount {
		if vv > v {
			result = typ
		}
	}
	return result
}

// CountTags returns the different tags in the database and the number of
// instances of each one.
func CountTags(rs io.ReadSeeker) (Tags, error) {
	tags := make(map[string]map[string]*TagCount)

	addTag := func(key, val string, typ ObjectType) {
		if _, ok := tags[key]; !ok {
			tags[key] = make(map[string]*TagCount)
		}
		if _, ok := tags[key][val]; !ok {
			tags[key][val] = &TagCount{
				Key:         key,
				Value:       val,
				ObjectCount: make(map[ObjectType]int),
			}
		}
		t := tags[key][val]
		t.ObjectCount[typ]++
		t.TotalCount++
		tags[key][val] = t
	}

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
			for key, val := range v.(*osmpbf.Node).Tags {
				addTag(key, val, Node)
			}
		case *osmpbf.Way:
			if w := v.(*osmpbf.Way); wayIsClosed(w) {
				for key, val := range w.Tags {
					addTag(key, val, ClosedWay)
				}
			} else {
				for key, val := range w.Tags {
					addTag(key, val, OpenWay)
				}
			}
		case *osmpbf.Relation:
			for key, val := range v.(*osmpbf.Relation).Tags {
				addTag(key, val, Relation)
			}
		default:
			return nil, fmt.Errorf("unknown type %T\n", vtype)
		}
	}
	var tagList Tags
	for _, d := range tags {
		for _, d2 := range d {
			tagList = append(tagList, d2)
		}
	}
	sort.Sort(sort.Reverse(&tagList))
	return tagList, nil
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
