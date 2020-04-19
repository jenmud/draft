package graph

// Node is a node in the graph.
type Node struct {
	UID string `json:"uid"`
	Label string `json:"label"`
	Properties map[string]Value `json:"properties"`
	inEdges map[string]string
	outEdges map[string]string
}