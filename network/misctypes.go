package network

import (
	"gonum.org/v1/gonum/graph"
)

type Link struct {
	src  graph.Node
	dest graph.Node
}

func NewLink(src, dest graph.Node) Link {
	return Link{src, dest}
}

// part of the graph.Edge interface
func (l Link) To() graph.Node {
	return l.dest
}

// part of the graph.Edge interface
func (l Link) From() graph.Node {
	return l.src
}

// part of the graph.Edge interface
func (l Link) ReversedEdge() graph.Edge {
	return Link{
		src:  l.dest,
		dest: l.src,
	}
}

type Vertex int

func NewVertex(id int64) Vertex {
	return Vertex(id)
}

// part of the graph.Node interface
func (u Vertex) ID() int64 {
	return int64(u)
}

// Adheres to the graph.Nodes interface which is an interface for interating across nodes.
type SliceIterator struct {
	nodes        []graph.Node
	currentIndex int
}

func (v *SliceIterator) Node() graph.Node {
	if v.currentIndex >= len(v.nodes) {
		return nil
	}
	v.currentIndex++
	return v.nodes[v.currentIndex-1]
}

func (v *SliceIterator) Next() bool {
	return v.currentIndex < len(v.nodes)
}

func (v *SliceIterator) Len() int {
	return len(v.nodes) - v.currentIndex
}

func (v *SliceIterator) Reset() {
	v.currentIndex = 0
}
