package graph

// Edge is a edge in the graph.
type Edge struct {
	UID string `json:"uid"`
	SourceUID string `json:"source_uid"`
	Label string `json:"label"`
	TargetUID string `json:"target_uid"`
	Properties map[string]Value `json:"properties"`
}
