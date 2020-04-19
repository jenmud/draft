package graph

// Value is a property value.
type Value struct {
	Type string `json:"type"`
	Value []byte `json:"value"`
}