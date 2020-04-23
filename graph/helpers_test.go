package graph

import (
	"testing"

	"github.com/jenmud/draft/graph/iterator"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {
	n1 := NewNode("node-1", "person")
	n2 := NewNode("node-2", "person")
	n3 := NewNode("node-3", "animal")
	n4 := NewNode("node-3", "pet")

	iter := iterator.New([]interface{}{n1, n2, n3, n4})
	filter := Filter{Type: LABEL, Value: Value{Type: "string", Value: []byte("person")}}

	mappedIter := mapper(iter, NODE, filter)
	actual := []interface{}{}
	for mappedIter.Next() {
		actual = append(actual, mappedIter.Value())
	}

	expected := []interface{}{n1, n2}
	assert.ElementsMatch(t, expected, actual)
}
