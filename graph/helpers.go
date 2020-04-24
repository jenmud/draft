package graph

import "fmt"

// convertPropertiesToKV converts a Property map to a array of key values.
func convertPropertiesToKV(props map[string][]byte) []KV {
	kvs := make([]KV, len(props))

	count := 0
	for k, v := range props {
		kvs[count] = KV{Key: k, Value: v}
	}

	return kvs
}

func mapperNode(iter Iterator, out chan<- Node) {
	for iter.Next() {
		node := iter.Value().(Node)
		out <- node
	}
	close(out)
}

// reducer is the filter.
func reducerNode(in chan Node, out chan Node, filter FilterType, kv KV) {
	for node := range in {
		switch filter {
		case LABEL:
			if node.Label == kv.Key {
				out <- node
			}
		default:
			panic(fmt.Sprintf("Filter %v is not supported", filter))
		}
	}

	close(out)
}
