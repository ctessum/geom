package osm

import (
	"fmt"
	"io"
	"runtime"
	"sort"

	"github.com/qedus/osmpbf"
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
			v = vv
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

// CountTags returns the different tags in the receiver and the number of
// instances of each one.
func (o *Data) CountTags() Tags {
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
	for _, n := range o.Nodes {
		for key, val := range n.Tags {
			addTag(key, val, Node)
		}
	}
	for _, w := range o.Ways {
		if wayIsClosed(w) {
			for key, val := range w.Tags {
				addTag(key, val, ClosedWay)
			}
		} else {
			for key, val := range w.Tags {
				addTag(key, val, OpenWay)
			}
		}
	}
	for _, r := range o.Relations {
		for key, val := range r.Tags {
			addTag(key, val, Relation)
		}
	}

	var tagList Tags
	for _, d := range tags {
		for _, d2 := range d {
			tagList = append(tagList, d2)
		}
	}
	sort.Sort(sort.Reverse(&tagList))
	return tagList
}
