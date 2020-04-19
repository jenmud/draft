package graph

import (
	"sync"
)

// New returns a new empty graph.
func New() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string]Edge),
	}
}

// Graph is a graph store.
type Graph struct {
	lock  sync.RWMutex
	nodes map[string]Node
	edges map[string]Edge
}

// See graph_node.go for all the node related methods
// See graph_edges.go for all the edge related methods
