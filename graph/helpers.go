package graph

import "github.com/jenmud/draft/graph/iterator"

// convertPropertiesToKV converts a Property map to a array of key values.
func convertPropertiesToKV(props map[string]Value) []KV {
	kvs := make([]KV, len(props))

	count := 0
	for k, v := range props {
		kvs[count] = KV{Key: k, Value: v}
	}

	return kvs
}

// mapper is used to run throug the items in the iterator
// sorting items into sorted piles.
// TODO: this should take some sort of interface
func mapper(iter Iterator, itemType ItemType, filter ...Filter) Iterator {
	mapped := []interface{}{}

	for iter.Next() {
		item := iter.Value().(Node)
		for _, f := range filter {
			switch f.Type {
			case LABEL:
				if item.Label == string(f.Value.Value) {
					mapped = append(mapped, iter.Value())
				}
			}
		}
	}

	return iterator.New(mapped)
}
