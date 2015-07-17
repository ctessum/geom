// Package route finds the shortest route between two points along a geometrical
// network (e.g., a road network). For now, all network links are assumed to be
// bi-directional (e.g., all roads are two-way).
package route

import (
	"fmt"
	"math"

	"github.com/ctessum/geom"
	"github.com/ctessum/geom/index/rtree"
	"github.com/ctessum/geom/op"
	"github.com/gonum/graph"
	"github.com/gonum/graph/path"
)

// A Network is a holder for network data (e.g., a road network)
type Network struct {
	nodes, edges   *rtree.Rtree
	neighbors      map[int]map[int]*edge
	nodeMap        map[int]*node
	maxID          int
	freeMap        map[int]struct{}
	minimizeOption MinimizeOption
	minimumSpeed   float64 // The minimum speed traveled on any link in the network.
}

// NewNetwork initializes a new Network where m determines how to choose
// the shortest route (either by Distance or Time).
func NewNetwork(m MinimizeOption) *Network {
	return &Network{
		neighbors:      make(map[int]map[int]*edge),
		nodeMap:        make(map[int]*node),
		maxID:          0,
		nodes:          rtree.NewTree(25, 50),
		edges:          rtree.NewTree(25, 50),
		minimizeOption: m,
		minimumSpeed:   math.Inf(1),
	}
}

// Has returns whether the node exists within the graph.
// It is not intended for direct use in this package.
func (net Network) Has(n graph.Node) bool {
	_, ok := net.nodeMap[n.ID()]
	return ok
}

// Nodes returns all the nodes in the graph.
// It is not intended for direct use in this package.
func (net Network) Nodes() []graph.Node {
	nodes := make([]graph.Node, len(net.nodeMap))
	i := 0
	for _, n := range net.nodeMap {
		nodes[i] = n
		i++
	}
	return nodes
}

// From returns all nodes that can be reached directly
// from the given node.
// It is not intended for direct use in this package.
func (net Network) From(n graph.Node) []graph.Node {
	if !net.Has(n) {
		return nil
	}
	neighbors := make([]graph.Node, len(net.neighbors[n.ID()]))
	i := 0
	for id := range net.neighbors[n.ID()] {
		neighbors[i] = net.nodeMap[id]
		i++
	}
	return neighbors
}

// HasEdge returns whether an edge exists between
// nodes x and y without considering direction.
// It is not intended for direct use in this package.
func (net Network) HasEdge(x, y graph.Node) bool {
	_, ok := net.neighbors[x.ID()][y.ID()]
	return ok
}

// Edge returns the edge from u to v if such an edge
// exists and nil otherwise. The node v must be directly
// reachable from u as defined by the From method.
// It is not intended for direct use in this package.
func (net Network) Edge(u, v graph.Node) graph.Edge {
	// We don't need to check if neigh exists because
	// it's implicit in the neighbors access.
	if !net.Has(u) {
		return nil
	}
	return net.neighbors[u.ID()][v.ID()]
}

// The math package only provides explicitly sized max
// values. This ensures we get the max for the actual
// type int.
const maxInt int = int(^uint(0) >> 1)

func (net *Network) newNodeID() int {
	if net.maxID != maxInt {
		net.maxID++
		return net.maxID
	}
	// I cannot foresee this ever happening, but just in case, we check.
	if len(net.nodeMap) == maxInt {
		panic("cannot allocate node: graph too large")
	}
	// Should not happen.
	panic("cannot allocate node id: no free id found")
}

// Check if there is already a node at this location, and if there
// is return that one, otherwise create a new node.
func (net *Network) newNode(p geom.Point) *node {
	nearest := net.nodes.NearestNeighbor(p)
	if nearest != nil && op.PointEquals(p, nearest.(*node).Point) {
		return nearest.(*node)
	}
	return &node{
		Point: p,
		id:    net.newNodeID(),
	}
}

func (net *Network) addNode(n graph.Node) {
	if _, exists := net.nodeMap[n.ID()]; exists {
		panic(fmt.Sprintf("route: node ID collision: %d", n.ID()))
	}
	net.nodeMap[n.ID()] = n.(*node)
	net.neighbors[n.ID()] = make(map[int]*edge)
	net.nodes.Insert(n.(*node))
}

// AddLink adds a new link l (which is a line string) to the Network, where
// speed is the speed traveled along the link and should have units that
// are compatible with the units of l (for instance, if l is in units of
// meters, and speed is in units of m/s, the time results will be in units
// of seconds).
func (net *Network) AddLink(l geom.LineString, speed float64) {
	from := net.newNode(l[0])
	to := net.newNode(l[len(l)-1])

	length := op.Length(l)
	e := &edge{
		LineString: l,
		start:      from,
		end:        to,
		length:     length,
		speed:      speed,
		time:       length / speed,
	}
	if e.speed < net.minimumSpeed {
		net.minimumSpeed = e.speed
	}
	fid := from.ID()
	tid := to.ID()
	if fid == tid {
		panic("concrete: adding self edge")
	}
	if !net.Has(from) {
		net.addNode(from)
	}
	if !net.Has(to) {
		net.addNode(to)
	}
	net.edges.Insert(e)
	net.neighbors[fid][tid] = e
	net.neighbors[tid][fid] = e
}

// Weight returns the weight associated with this edge.
// It is not intended for direct use in this package.
func (net *Network) Weight(e graph.Edge) float64 {
	if n, ok := net.neighbors[e.From().ID()]; ok {
		if we, ok := n[e.To().ID()]; ok {
			switch net.minimizeOption {
			// If we're optimizing by time, return use the minimum speed to
			// calculate the time to ensure the heuristic is less than the actual
			// value
			case Time:
				return we.time
			case Distance:
				// If we're optimizing by distance, just return the distance.
				return we.length
			default:
				panic(fmt.Errorf("Invalid MinimizeOption %v", net.minimizeOption))
			}
		}
	}
	panic("route: attempting to find an edge that is not in the graph")
}

type edge struct {
	geom.LineString
	start, end          *node
	length, speed, time float64
}

// From gives the beginning point of this edge.
func (e edge) From() graph.Node {
	return e.start
}

// To gives the final point of this edge.
func (e edge) To() graph.Node {
	return e.end
}

type node struct {
	geom.Point
	id int
}

func (n node) ID() int {
	return n.id
}

// MinimizeOption specifies how the shortest route should be chosen.
type MinimizeOption float64

const (
	// Distance specifies that we are looking to travel the minimum distance.
	Distance MinimizeOption = iota
	// Time specifies that we are looking to travel the minimum time.
	Time
)

// ShortestRoute calculates the shortest route along the network between the
// from and to points. It returns the route ("path"), the distance traveled
// along the route ("distance"), the time it would take travel along the route
// ("time"; this does not count time spent getting to and from the network),
// the distance traveled from the starting
// point to get to the nearest node (e.g., intersection) along the route
// ("startDistance") and the distance traveled to the ending point from
// the nearest node along the route ("endDistance"). This function does
// not change the Network, so multiple function calls can be run concurrently.
func (net Network) ShortestRoute(from, to geom.Point) (
	route geom.MultiLineString, distance, time, startDistance, endDistance float64) {
	startNode := net.nodes.NearestNeighbor(from).(*node)
	endNode := net.nodes.NearestNeighbor(to).(*node)
	startDistance = op.Distance(from, startNode.Point)
	endDistance = op.Distance(to, endNode.Point)
	shortest, _ := path.AStar(startNode, endNode, net, net.costHeuristic)
	nodes, _ := shortest.To(endNode)
	for i := 0; i < len(nodes)-1; i++ {
		e, ok := net.neighbors[nodes[i].ID()][nodes[i+1].ID()]
		if !ok {
			panic("route: missing edge; this shouldn't happen")
		}
		route = append(route, e.LineString)
		distance += e.length
		time += e.time
	}
	return
}

// costHeuristic provides a cost estimate this is guaranteed to be equal
// to or less than the actual cost.
func (net *Network) costHeuristic(x, y graph.Node) float64 {
	distance := op.Distance(x.(*node).Point, y.(*node).Point)
	switch net.minimizeOption {
	// If we're optimizing by time, return use the minimum speed to
	// calculate the time to ensure the heuristic is less than the actual
	// value
	case Time:
		return distance / net.minimumSpeed
	case Distance:
		// If we're optimizing by distance, just return the distance.
		return distance
	default:
		panic(fmt.Errorf("Invalid MinimizeOption %v", net.minimizeOption))
	}
}
