package graph

import "fmt"

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
