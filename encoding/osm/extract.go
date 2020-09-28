// Package osm extracts and manipulates OpenStreetMap (OSM) data. Refer to
// openstreetmap.org for more information about OSM data.
package osm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/paulmach/osm/osmxml"
	"golang.org/x/sync/errgroup"
)

// ExtractFile extracts OpenStreetMap data from the given file path,
// determining whether it is an XML or PBF file from the extension
// (.osm or .pbf, respectively).
// keep determines which records are included in the output.
// keepTags determines whether the tags and other metadata
// should be removed from each record to reduce memeory use.
func ExtractFile(ctx context.Context, file string, keep KeepFunc, keepTags bool) (*Data, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("osm: %v", err)
	}
	defer f.Close()
	x := filepath.Ext(file)
	switch x {
	case ".pbf":
		return ExtractPBF(ctx, f, keep, keepTags)
	case ".osm":
		return ExtractXML(ctx, f, keep, keepTags)
	default:
		return nil, fmt.Errorf("osm: invalid file extension '%s'", x)
	}
}

// ExtractPBF extracts OpenStreetMap data from osm.pbf file rs.
// keep determines which records are included in the output.
func ExtractPBF(ctx context.Context, rs io.ReadSeeker, keep KeepFunc, keepTags bool) (*Data, error) {
	//scanFunc := func() osm.Scanner { return osmpbf.New(ctx, rs, runtime.GOMAXPROCS(-1)) }
	scanFunc := func() osm.Scanner { return osmpbf.New(ctx, rs, 1) }
	return extract(ctx, rs, scanFunc, keep, keepTags)
}

// ExtractXML extracts OpenStreetMap data from osm file rs.
// keep determines which records are included in the output.
func ExtractXML(ctx context.Context, rs io.ReadSeeker, keep KeepFunc, keepTags bool) (*Data, error) {
	scanFunc := func() osm.Scanner { return osmxml.New(ctx, rs) }
	return extract(ctx, rs, scanFunc, keep, keepTags)
}

// extract extracts OpenStreetMap data from file rs.
// keep determines which records are included in the output.
func extract(ctx context.Context, rs io.ReadSeeker, scanFunc func() osm.Scanner, keep KeepFunc, keepTags bool) (*Data, error) {
	o := &Data{
		Nodes:     make(map[osm.NodeID]*Node),
		Ways:      make(map[osm.WayID]*Way),
		Relations: make(map[osm.RelationID]*Relation),

		dependentNodes:     make(map[osm.NodeID]empty),
		dependentWays:      make(map[osm.WayID]empty),
		dependentRelations: make(map[osm.RelationID]empty),
	}

	nprocs := runtime.GOMAXPROCS(-1)
	var passMX sync.Mutex

	needAnotherPass := true
	for needAnotherPass {
		needAnotherPass = false

		eg := new(errgroup.Group)
		objChan := make(chan osm.Object, nprocs)
		for i := 0; i < nprocs; i++ {
			eg.Go(func() error {
				for obj := range objChan {
					switch objType := obj.(type) {
					case *osm.Node:
						o.processNode(obj.(*osm.Node), keep, keepTags)
					case *osm.Way:
						if o.processWay(obj.(*osm.Way), keep, keepTags) {
							passMX.Lock()
							needAnotherPass = true
							passMX.Unlock()
						}
					case *osm.Relation:
						if o.processRelation(obj.(*osm.Relation), keep, keepTags) {
							passMX.Lock()
							needAnotherPass = true
							passMX.Unlock()
						}
					case *osm.Note, *osm.Bounds, *osm.User:
					default:
						return fmt.Errorf("unknown type %T", objType)
					}
				}
				return nil
			})
		}

		if _, err := rs.Seek(0, 0); err != nil {
			return nil, err
		}
		scanner := scanFunc()
		for scanner.Scan() {
			objChan <- scanner.Object()
		}
		close(objChan)
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		if err := scanner.Close(); err != nil {
			return nil, err
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}
	}
	return o, nil
}

// ExtractTag extracts OpenStreetMap data with the given tag set to one of the
// given values.
// keepTags determines whether the tags and other metadata
// should be removed from each record to reduce memeory use.
func ExtractTag(rs io.ReadSeeker, tag string, keepTags bool, values ...string) (*Data, error) {
	return ExtractPBF(context.Background(), rs, KeepTags(map[string][]string{tag: values}), keepTags)
}

// ObjectType specifies the valid OpenStreetMap types.
type ObjectType int

const (
	// NodeType is an OpenStreetMap node.
	NodeType ObjectType = iota
	// WayType can be either open or closed.
	WayType
	// ClosedWayType is an OpenStreetMap way that is closed (i.e., a polygon).
	ClosedWayType
	// OpenWayType is an OpenStreetMap way that is open (i.e., a line string).
	OpenWayType
	// RelationType is an OpenStreetMap relation.
	RelationType
)

type empty struct{}

// Data holds OpenStreetMap data and relationships.
type Data struct {
	Nodes     map[osm.NodeID]*Node
	Ways      map[osm.WayID]*Way
	Relations map[osm.RelationID]*Relation

	// These list the objects that are dependent on other objects,
	// and the objects that they are dependent on.
	dependentNodes     map[osm.NodeID]empty
	dependentWays      map[osm.WayID]empty
	dependentRelations map[osm.RelationID]empty

	nodeMX              sync.RWMutex
	wayMX               sync.RWMutex
	relationMX          sync.RWMutex
	dependentNodeMX     sync.RWMutex
	dependentWayMX      sync.RWMutex
	dependentRelationMX sync.RWMutex
}

// Filter returns a copy of the receiver where only objects
// selected by keep are retained.
func (o *Data) Filter(keep KeepFunc) *Data {
	out := &Data{
		Nodes:     make(map[osm.NodeID]*Node),
		Ways:      make(map[osm.WayID]*Way),
		Relations: make(map[osm.RelationID]*Relation),

		dependentNodes:     make(map[osm.NodeID]empty),
		dependentWays:      make(map[osm.WayID]empty),
		dependentRelations: make(map[osm.RelationID]empty),
	}

	needAnotherPass := true
	for needAnotherPass {
		needAnotherPass = false
		for _, n := range o.Nodes {
			out.processNodeNoCopy(n, keep, false)
		}
		for _, w := range o.Ways {
			if out.processWayNoCopy(w, keep, false) {
				needAnotherPass = true
			}
		}
		for _, r := range o.Relations {
			if out.processRelationNoCopy(r, keep, false) {
				needAnotherPass = true
			}
		}
	}
	return out
}

// hasTag checks if tags[t] is one of the values in wantTags. If len(v) == 0,
// the function will return true is tags[t] has any value.
func hasTag(tags osm.Tags, wantTags map[string][]string) bool {
	for _, t := range tags {
		wantTagValues, ok := wantTags[t.Key]
		if !ok {
			continue
		}
		if len(wantTagValues) == 0 {
			return true
		}
		for _, wantTagValue := range wantTagValues {
			if t.Value == wantTagValue {
				return true
			}
		}
	}
	return false
}

func (o *Data) hasNeedNode(id osm.NodeID) (has, need bool) {
	o.nodeMX.RLock()
	defer o.nodeMX.RUnlock()
	if _, ok := o.Nodes[id]; ok {
		has = true
		return
	}
	o.dependentNodeMX.RLock()
	defer o.dependentNodeMX.RUnlock()
	if _, ok := o.dependentNodes[id]; ok {
		need = true
	}
	return
}

func (o *Data) hasNeedWay(id osm.WayID) (has, need bool) {
	o.wayMX.RLock()
	defer o.wayMX.RUnlock()
	if _, ok := o.Ways[id]; ok {
		has = true
		return
	}
	o.dependentWayMX.RLock()
	defer o.dependentWayMX.RUnlock()
	if _, ok := o.dependentWays[id]; ok {
		need = true
	}
	return
}

func (o *Data) hasNeedRelation(id osm.RelationID) (has, need bool) {
	o.relationMX.RLock()
	defer o.relationMX.RUnlock()
	if _, ok := o.Relations[id]; ok {
		has = true
		return
	}
	o.dependentRelationMX.RLock()
	defer o.dependentRelationMX.RUnlock()
	if _, ok := o.dependentRelations[id]; ok {
		need = true
	}
	return
}

// If the node has the tag we want, add it to the list.
func (o *Data) processNode(n *osm.Node, keep KeepFunc, keepTags bool) {
	hasNode, needNode := o.hasNeedNode(n.ID)
	if hasNode {
		return
	}
	if keep(o, n) || needNode {
		o.nodeMX.Lock()
		o.Nodes[n.ID] = copyNode(n, keepTags)
		o.nodeMX.Unlock()
	}
}

func (o *Data) processNodeNoCopy(n *Node, keep KeepFunc, keepTags bool) {
	hasNode, needNode := o.hasNeedNode(n.ID)
	if hasNode {
		return
	}
	if keep(o, n) || needNode {
		o.nodeMX.Lock()
		o.Nodes[n.ID] = n
		o.nodeMX.Unlock()
	}
}

// If the way has the tag we want or if we've determined that it's
// part of a relation that we want, store the way and the IDs of its dependent nodes.
func (o *Data) processWay(w *osm.Way, keep KeepFunc, keepTags bool) (anotherPass bool) {
	hasWay, needWay := o.hasNeedWay(w.ID)
	if hasWay {
		return
	}
	if keep(o, w) || needWay {
		o.wayMX.Lock()
		o.Ways[w.ID] = copyWay(w, keepTags)
		o.wayMX.Unlock()
		for _, n := range w.Nodes {
			if _, needNode := o.hasNeedNode(n.ID); !needNode {
				o.dependentNodeMX.Lock()
				o.dependentNodes[n.ID] = empty{}
				o.dependentNodeMX.Unlock()
				anotherPass = true
			}
		}
	}
	return
}

func (o *Data) processWayNoCopy(w *Way, keep KeepFunc, keepTags bool) (anotherPass bool) {
	hasWay, needWay := o.hasNeedWay(w.ID)
	if hasWay {
		return
	}
	if keep(o, w) || needWay {
		o.Ways[w.ID] = w
		for _, n := range w.Nodes {
			if _, needNode := o.hasNeedNode(n); !needNode {
				o.dependentNodeMX.Lock()
				o.dependentNodes[n] = empty{}
				o.dependentNodeMX.Unlock()
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
func (o *Data) processRelation(r *osm.Relation, keep KeepFunc, keepTags bool) (anotherPass bool) {
	hasRelation, needRelation := o.hasNeedRelation(r.ID)
	if hasRelation {
		return
	}
	if keep(o, r) || needRelation {
		o.relationMX.Lock()
		o.Relations[r.ID] = copyRelation(r, keepTags)
		o.relationMX.Unlock()
		for _, m := range r.Members {
			switch m.Type {
			case osm.TypeNode:
				if _, needNode := o.hasNeedNode(osm.NodeID(m.Ref)); !needNode {
					o.dependentNodeMX.Lock()
					o.dependentNodes[osm.NodeID(m.Ref)] = empty{}
					o.dependentNodeMX.Unlock()
					anotherPass = true
				}
			case osm.TypeWay:
				if _, needWay := o.hasNeedWay(osm.WayID(m.Ref)); !needWay {
					o.dependentWayMX.Lock()
					o.dependentWays[osm.WayID(m.Ref)] = empty{}
					o.dependentWayMX.Unlock()
					anotherPass = true
				}
			case osm.TypeRelation:
				if _, needR := o.hasNeedRelation(osm.RelationID(m.Ref)); !needR {
					o.dependentRelationMX.Lock()
					o.dependentRelations[osm.RelationID(m.Ref)] = empty{}
					o.dependentRelationMX.Unlock()
					anotherPass = true
				}
			default:
				panic(fmt.Errorf("unknown member type %v", m.Type))
			}
		}
	}
	return
}

func (o *Data) processRelationNoCopy(r *Relation, keep KeepFunc, keepTags bool) (anotherPass bool) {
	hasRelation, needRelation := o.hasNeedRelation(r.ID)
	if hasRelation {
		return
	}
	if keep(o, r) || needRelation {
		o.relationMX.Lock()
		o.Relations[r.ID] = r
		o.relationMX.Unlock()
		for _, m := range r.Members {
			switch m.Type {
			case osm.TypeNode:
				if _, needNode := o.hasNeedNode(osm.NodeID(m.Ref)); !needNode {
					o.dependentNodeMX.Lock()
					o.dependentNodes[osm.NodeID(m.Ref)] = empty{}
					o.dependentNodeMX.Unlock()
					anotherPass = true
				}
			case osm.TypeWay:
				if _, needWay := o.hasNeedWay(osm.WayID(m.Ref)); !needWay {
					o.dependentWayMX.Lock()
					o.dependentWays[osm.WayID(m.Ref)] = empty{}
					o.dependentWayMX.Unlock()
					anotherPass = true
				}
			case osm.TypeRelation:
				if _, needR := o.hasNeedRelation(osm.RelationID(m.Ref)); !needR {
					o.dependentRelationMX.Lock()
					o.dependentRelations[osm.RelationID(m.Ref)] = empty{}
					o.dependentRelationMX.Unlock()
					anotherPass = true
				}
			default:
				panic(fmt.Errorf("unknown member type %v", m.Type))
			}
		}
	}
	return
}

// Node holds a subset of the information specifying an
// OpenStreetMap node.
type Node struct {
	ID   osm.NodeID
	Lat  float64
	Lon  float64
	Tags osm.Tags
}

func copyNode(n *osm.Node, keepTags bool) *Node {
	o := &Node{
		ID:  n.ID,
		Lat: n.Lat,
		Lon: n.Lon,
	}
	if keepTags {
		o.Tags = n.Tags
	}
	return o
}

// Way holds a subset of the information specifying an
// OpenStreetMap way.
type Way struct {
	ID    osm.WayID
	Nodes []osm.NodeID
	Tags  osm.Tags
}

func copyWay(w *osm.Way, keepTags bool) *Way {
	o := &Way{
		ID: w.ID,
	}
	o.Nodes = make([]osm.NodeID, len(w.Nodes))
	for i, n := range w.Nodes {
		o.Nodes[i] = n.ID
	}
	if keepTags {
		o.Tags = w.Tags
	}
	return o
}

// Relation holds a subset of the information specifying an
// OpenStreetMap way.
type Relation struct {
	ID      osm.RelationID
	Members []Member
	Tags    osm.Tags
}

// Member is a member of a relation.
type Member struct {
	Ref  int64
	Type osm.Type
}

func copyRelation(r *osm.Relation, keepTags bool) *Relation {
	o := &Relation{
		ID: r.ID,
	}
	o.Members = make([]Member, len(r.Members))
	for i, m := range r.Members {
		o.Members[i] = Member{Ref: m.Ref, Type: m.Type}
	}
	if keepTags {
		o.Tags = r.Tags
	}
	return o
}
