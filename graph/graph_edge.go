package graph

import (
	"fmt"

	"github.com/jenmud/draft/graph/iterator"
)

// HasEdge returns true if the graph has a edge with the provided uid.
func (g *Graph) HasEdge(uid string) bool {
	g.lock.RLock()
	defer g.lock.RUnlock()

	_, ok := g.edges[uid]
	return ok
}

// AddEdge adds a new edge to the graph.
func (g *Graph) AddEdge(uid, sourceUID, label, targetUID string, kv ...KV) (Edge, error) {
	if !g.HasNode(sourceUID) {
		return Edge{}, fmt.Errorf("[AddEdge] No such not with UID %s", sourceUID)
	}

	if !g.HasNode(targetUID) {
		return Edge{}, fmt.Errorf("[AddEdge] No such not with UID %s", targetUID)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.edges[uid]; ok {
		return Edge{}, fmt.Errorf("[AddEdge] Edge UID %s already exists", uid)
	}

	edge := NewEdge(uid, sourceUID, label, targetUID, kv...)
	g.edges[edge.UID] = edge
	return edge, nil
}

// RemoveEdge removes the edge from the graph.
func (g *Graph) RemoveEdge(uid string) {
	g.lock.Lock()
	defer g.lock.Unlock()
	delete(g.edges, uid)
}

// Edge returns the edge with the provided uid.
func (g *Graph) Edge(uid string) (Edge, error) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	edge, ok := g.edges[uid]
	if !ok {
		return Edge{}, fmt.Errorf("[GetEdge] No such edge with UID %s found", uid)
	}

	return edge, nil
}

// Edges returns a edge iterator.
func (g *Graph) Edges() Iterator {
	g.lock.RLock()
	defer g.lock.RUnlock()

	edges := make([]interface{}, len(g.edges))
	count := 0
	for _, edge := range g.edges {
		edges[count] = edge
		count++
	}

	return iterator.New(edges)
}
