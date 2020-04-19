package graph

import (
	"fmt"
	"sync"

	"github.com/jenmud/draft/graph/iterator"
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

// HasNode returns true if the graph has a node with the provided uid.
func (g *Graph) HasNode(uid string) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	_, ok := g.nodes[uid]
	return ok
}

// AddNode adds a new node to the graph.
func (g *Graph) AddNode(uid, label string, kv ...KV) (Node, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.nodes[uid]; ok {
		return Node{}, fmt.Errorf("[AddNode] Node UID %s already exists", uid)
	}

	node := NewNode(uid, label, kv...)
	g.nodes[node.UID] = node
	return node, nil
}

// GetNode returns the node with the provided uid.
func (g *Graph) GetNode(uid string) (Node, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	node, ok := g.nodes[uid]
	if !ok {
		return Node{}, fmt.Errorf("[GetNode] No such node with UID %s found", uid)
	}

	return node, nil
}

// Nodes returns a node iterator.
func (g *Graph) Nodes() Iter {
	g.lock.RLock()
	defer g.lock.RUnlock()

	nodes := make([]Node, len(g.nodes))
	count := 0
	for _, node := range g.nodes {
		nodes[count] = node
		count++
	}

	return iterator.New(nodes)
}
