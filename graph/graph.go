package graph

// New returns a new empty graph.
func New() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string]Edge),
	}
}

// Graph is a graph store.
type Graph struct {
	nodes map[string]Node
	edges map[string]Edge
}
