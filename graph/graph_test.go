package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNode(t *testing.T) {
	g := New()
	expected := NewNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	actual, err := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddNode_Duplicate(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	actual, err := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	assert.NotNil(t, err)
	assert.Equal(t, Node{}, actual)
}

func TestRemoveNode(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	g.RemoveNode("abcd-1234")
	assert.Equal(t, false, g.HasNode("abcd-1234"))
}

func TestHasNode(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	assert.Equal(t, true, g.HasNode("abcd-1234"))
}

func TestHasNode_not_found(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	assert.Equal(t, false, g.HasNode("missing"))
}

func TestGetNode(t *testing.T) {
	g := New()
	expected, _ := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	actual, err := g.GetNode("abcd-1234")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestGetNode_not_found(t *testing.T) {
	g := New()
	g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	actual, err := g.GetNode("abcd-1234-missing")
	assert.NotNil(t, err)
	assert.Equal(t, Node{}, actual)
}

func TestNodes(t *testing.T) {
	g := New()
	expected1, _ := g.AddNode("abcd-1234", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("foo")}})
	expected2, _ := g.AddNode("abcd-4321", "person", KV{Key: "name", Value: Value{Type: "string", Value: []byte("bar")}})

	expected := []Node{expected1, expected2}
	actual := []Node{}

	iter := g.Nodes()
	for iter.Next() {
		actual = append(actual, iter.Value().(Node))
	}

	assert.ElementsMatch(t, expected, actual)
}
