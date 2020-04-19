package graph

// NewNode returns a new node instance.
func NewNode(uid, label string, kv ...KV) Node {
	return Node{
		UID:        uid,
		Label:      label,
		Properties: NewProperties(kv...),
		inEdges:    make(map[string]string),
		outEdges:   make(map[string]string),
	}
}

// Node is a node in the graph.
type Node struct {
	UID        string           `json:"uid"`
	Label      string           `json:"label"`
	Properties map[string]Value `json:"properties"`
	inEdges    map[string]string
	outEdges   map[string]string
}
