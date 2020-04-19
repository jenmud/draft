package graph

// Graph is a graph store.
type Graph struct {
	nodes map[string]Node
	edges map[string]Edge
}