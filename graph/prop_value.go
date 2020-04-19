package graph

// Value is a property value.
type Value struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}

// KV is a property key and value pair.
type KV struct {
	Key   string `json:"key"`
	Value Value
}

// NewProperties takes one or more key value pairs and returns a property map.
func NewProperties(kv ...KV) map[string]Value {
	props := make(map[string]Value)

	for _, pair := range kv {
		props[pair.Key] = pair.Value
	}

	return props
}
