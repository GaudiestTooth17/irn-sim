package network

import (
	"math"

	"github.com/GaudiestTooth17/irn-sim/sets"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/mat"
)

type AdjacencyList struct {
	// nodes keyed by id
	nodes []graph.Node
	// node ID to the node's neighbors
	adjList map[int64][]graph.Node
	// adjacency matrix
	m *mat.Dense
	// distance matrix
	dm *mat.Dense
}

func NewAdjacencyList(nodes []graph.Node, adjList map[int64][]graph.Node) *AdjacencyList {
	return &AdjacencyList{
		nodes:   nodes,
		adjList: adjList,
		m:       nil,
		dm:      nil,
	}
}

// part of the graph.Graph interface
func (g *AdjacencyList) Node(id int64) graph.Node {
	return g.nodes[id]
}

// part of the graph.Graph interface
func (g *AdjacencyList) Nodes() graph.Nodes {
	allNodes := make([]graph.Node, len(g.nodes))
	for _, node := range g.nodes {
		allNodes[node.ID()] = node
	}
	return &SliceIterator{allNodes, 0}
}

// part of the graph.Graph interface
func (g *AdjacencyList) From(id int64) graph.Nodes {
	return &SliceIterator{g.adjList[id], 0}
}

// part of the graph.Graph interface
func (g *AdjacencyList) HasEdgeBetween(xid, yid int64) bool {
	if g.m != nil {
		return g.m.At(int(xid), int(yid)) == 1
	}
	for _, v := range g.adjList[xid] {
		if v.ID() == yid {
			return true
		}
	}
	return false
}

// part of the graph.Graph interface
func (g *AdjacencyList) Edge(uid, vid int64) graph.Edge {
	u := g.nodes[uid]
	for _, v := range g.adjList[uid] {
		if v.ID() == vid {
			return Link{u, v}
		}
	}
	return nil
}

// part of the graph.Undirect interface
func (g *AdjacencyList) EdgeBetween(xid, yid int64) graph.Edge {
	edge := g.Edge(xid, yid)
	if edge == nil {
		edge = g.Edge(yid, xid)
	}
	return edge
}

// Return the adjacency matrix
func (n *AdjacencyList) M() *mat.Dense {
	if n.m == nil {
		N := int64(n.N())
		backingData := make([]float64, N*N)
		// u goes from 0 to N-1
		for uID, neighbors := range n.adjList {
			// v will just be from the set of u's neighbors
			for _, v := range neighbors {
				backingData[uID*N+v.ID()] = 1
				backingData[v.ID()*N+uID] = 1
			}
		}
		n.m = mat.NewDense(int(N), int(N), backingData)
	}
	return n.m
}

// Return the number of nodes in the network
func (n *AdjacencyList) N() int {
	return len(n.adjList)
}

func (net *AdjacencyList) NodesWithin(nodeID int64, distance int) sets.IntSet {
	if net.dm == nil {
		net.initDistMatrix()
	}

	maxDist := float64(distance)
	id := int(nodeID)
	N, _ := net.m.Dims()

	nodes := sets.EmptyIntSet()
	for node := 0; node < N; node++ {
		dist := net.dm.At(id, node)
		if dist < maxDist {
			nodes.Add(node)
		}
	}
	return nodes
}

func (net *AdjacencyList) initDistMatrix() {
	M := net.m
	N, _ := net.m.Dims()
	dm := mat.DenseCopyOf(net.m)
	dm.Apply(func(i, j int, v float64) float64 {
		if v < 1 {
			return math.Inf(1)
		}
		return v
	}, dm)
	x := mat.DenseCopyOf(M)

	for d := 0; d < N; d++ {
		oldX := x
		x = mat.NewDense(N, N, nil)
		x.Mul(oldX, M)
		x.Apply(func(i, j int, v float64) float64 {
			if v > 0 {
				return 1
			}
			return v
		}, x)
		if mat.Equal(oldX, x) {
			break
		}

		// For every new path we know that the distance is d + 2 since
		// d starts at 0 and we already have everything of distance 1
		dm.Apply(func(i, j int, v float64) float64 {
			if v == math.Inf(1) && x.At(i, j) != 0 {
				return float64(d + 2)
			}
			return v
		}, dm)
	}

	// Sets the self distance to 0 before assignment
	for i := 0; i < N; i++ {
		dm.Set(i, i, 0)
	}

	net.dm = dm
}
