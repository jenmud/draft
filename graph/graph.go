package graph

import (
	"encoding/json"
	"sync"
)

// New returns a new empty graph.
func New() *Graph {
	return &Graph{
		nodes: make(map[string]Node),
		edges: make(map[string]Edge),
	}
}

// NewFromJSON takes a JSON formatted output and returns a new Graph.
func NewFromJSON(r io.Reader) (*Graph, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	graph := New()
	err = json.Unmarshal(data, graph)
	return graph, err
}

// Graph is a graph store.
type Graph struct {
	lock  sync.RWMutex
	nodes map[string]Node
	edges map[string]Edge
}

// See graph_node.go for all the node related methods
// See graph_edges.go for all the edge related methods

// MarshalJSON marchals the graph into a JSON format.
func (g *Graph) MarshalJSON() ([]byte, error) {
	type G struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}

	nodes := g.Nodes()
	edges := g.Edges()

	graph := G{
		Nodes: make([]Node, nodes.Size()),
		Edges: make([]Edge, edges.Size()),
	}

	ncount := 0
	for nodes.Next() {
		node := nodes.Value().(Node)
		graph.Nodes[ncount] = node
		ncount++
	}

	ecount := 0
	for edges.Next() {
		edge := edges.Value().(Edge)
		graph.Edges[ecount] = edge
		ecount++
	}

	return json.Marshal(graph)
}

// UnmarshalJSON unmarshals JSON data into the graph.
func (g *Graph) UnmarshalJSON(b []byte) error {
	type G struct {
		Nodes []Node `json:"nodes"`
		Edges []Edge `json:"edges"`
	}

	graph := G{}
	if err := json.Unmarshal(b, &graph); err != nil {
		return err
	}

	for _, node := range graph.Nodes {
		if _, err := g.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
			return err
		}
	}

	for _, edge := range graph.Edges {
		if _, err := g.AddEdge(edge.UID, edge.SourceUID, edge.Label, edge.TargetUID, convertPropertiesToKV(edge.Properties)...); err != nil {
			return err
		}
	}

	return nil
}
