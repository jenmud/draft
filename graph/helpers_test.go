package graph

import (
	"testing"

	"github.com/jenmud/draft/graph/iterator"
	"github.com/stretchr/testify/assert"
)

func TestMapperNode(t *testing.T) {
	n1 := NewNode("node-1", "person")
	n2 := NewNode("node-2", "person")
	n3 := NewNode("node-3", "animal")
	n4 := NewNode("node-3", "pet")

	iter := iterator.New([]interface{}{n1, n2, n3, n4})
	out := make(chan Node, iter.Size())

	mapperNode(iter, out)

	expected := []Node{n1, n2, n3, n4}
	actual := make([]Node, iter.Size())

	count := 0
	for n := range out {
		actual[count] = n
		count++
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestReduce__label(t *testing.T) {
	n1 := NewNode("node-1", "person")
	n2 := NewNode("node-2", "person")
	n3 := NewNode("node-3", "animal")
	n4 := NewNode("node-3", "pet")

	iter := iterator.New([]interface{}{n1, n2, n3, n4})
	in := make(chan Node, iter.Size())
	out := make(chan Node, iter.Size())

	mapperNode(iter, in)
	reducerNode(in, out, LABEL, KV{Key: "person"})

	expected := []Node{n1, n2}
	actual := make([]Node, 2)

	count := 0
	for n := range out {
		actual[count] = n
		count++
	}

	assert.ElementsMatch(t, expected, actual)
}
