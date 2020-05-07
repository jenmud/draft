package graph

import (
	"bytes"
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

// UpdateEdge updates the graph edge with the new edge.
func (g *Graph) UpdateEdge(edge Edge) (Edge, error) {
	if !g.HasEdge(edge.UID) {
		return edge, fmt.Errorf("[UpdateEdge] Edge does not exists, can not update edge %s", edge)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	g.edges[edge.UID] = edge
	return edge, nil
}

// AddEdge adds a new edge to the graph.
func (g *Graph) AddEdge(uid, sourceUID, label, targetUID string, kv ...KV) (Edge, error) {
	if !g.HasNode(sourceUID) {
		return Edge{}, fmt.Errorf("[AddEdge] No such not with UID %s", sourceUID)
	}

	if !g.HasNode(targetUID) {
		return Edge{}, fmt.Errorf("[AddEdge] No such not with UID %s", targetUID)
	}

	source, err := g.Node(sourceUID)
	if err != nil {
		return Edge{}, fmt.Errorf("[AddEdge] %s", err)
	}

	target, err := g.Node(targetUID)
	if err != nil {
		return Edge{}, fmt.Errorf("[AddEdge] %s", err)
	}

	g.lock.Lock()
	defer g.lock.Unlock()

	if _, ok := g.edges[uid]; ok {
		return Edge{}, fmt.Errorf("[AddEdge] Edge UID %s already exists", uid)
	}

	edge := NewEdge(uid, sourceUID, label, targetUID, kv...)
	g.edges[edge.UID] = edge

	// (source)->(target)
	source.outEdges[edge.UID] = struct{}{}
	target.inEdges[edge.UID] = struct{}{}

	return edge, nil
}

// RemoveEdge removes the edge from the graph.
func (g *Graph) RemoveEdge(uid string) error {
	edge, err := g.Edge(uid)
	if err != nil {
		return fmt.Errorf("[RemoveEdge] %s", err)
	}

	// (source)->(target)
	source, err := g.Node(edge.SourceUID)
	if err != nil {
		// this is only here for safty, but we shoud not
		// get into a situation where this error is returned.
		return fmt.Errorf("[RemoveEdge] %s", err)
	}
	delete(source.outEdges, uid)

	target, err := g.Node(edge.TargetUID)
	if err != nil {
		// this is only here for safty, but we shoud not
		// get into a situation where this error is returned.
		return fmt.Errorf("[RemoveEdge] %s", err)
	}
	delete(target.inEdges, uid)

	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.edges, uid)
	return nil
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

// edgeMapper converts interfaces into Edges
func edgeMapper(in <-chan interface{}, out chan<- Edge) {
	for item := range in {
		out <- item.(Edge)
	}
	close(out)
}

// edgeLabelReducer filters for edges that have the given labels.
func edgeLabelReducer(labels []string, in <-chan Edge, out chan<- Edge) {
	for edge := range in {
		if len(labels) == 0 {
			out <- edge
			continue
		}

		for _, label := range labels {
			if edge.Label == label {
				out <- edge
				continue
			}
		}
	}

	close(out)
}

// edgeSourceTargetReducer filters for edges that have the given source or target uids.
// If source or target is empty, then all edges will be returned.
func edgeSourceTargetReducer(source, target string, in <-chan Edge, out chan<- Edge) {
	for edge := range in {
		// if there is no source or target add the edge
		if source == "" && target == "" {
			out <- edge
			continue
		}

		// if there is a source and not target, then only filter on source
		if target == "" && source == edge.SourceUID {
			out <- edge
			continue
		}

		// if there is a target and not source, then only filter on target
		if source == "" && target == edge.TargetUID {
			out <- edge
			continue
		}
	}

	close(out)
}

// edgePropReducer filters for edges that have the given properties.
func edgePropReducer(props map[string][]byte, in <-chan Edge, out chan<- Edge) {
	for edge := range in {
		if len(props) == 0 {
			out <- edge
			continue
		}

		var allMatched = false
		for key, value := range props {
			nvalue, ok := edge.Properties[key]
			if !ok {
				allMatched = false
				break
			}

			if !bytes.Equal(value, nvalue) {
				allMatched = false
				break
			}

			allMatched = true
		}

		if allMatched {
			out <- edge
		}
	}

	close(out)
}

// EdgesBy returns a edge iterator with filtered edges.
// If labels is an empty list, then any label will be used.
// If props is an empty map, no properties will be used for filtering.
func (g *Graph) EdgesBy(source string, labels []string, target string, props map[string][]byte) Iterator {
	g.lock.RLock()
	defer g.lock.RUnlock()

	in := make(chan Edge, len(g.edges))
	labelFiltered := make(chan Edge, len(g.edges))
	sourceTargetFiltered := make(chan Edge, len(g.edges))
	final := make(chan Edge)

	go edgeMapper(g.Edges().Channel(), in)
	go edgeLabelReducer(labels, in, labelFiltered)
	go edgeSourceTargetReducer(source, target, labelFiltered, sourceTargetFiltered)
	go edgePropReducer(props, sourceTargetFiltered, final)

	edges := []interface{}{}
	for edge := range final {
		edges = append(edges, edge)
	}

	return iterator.New(edges)
}

// EdgeCount returns the total number of edges in the graph.
func (g *Graph) EdgeCount() int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.edges)
}
