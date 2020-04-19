package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	expected := NewEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})
	actual, err := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestAddEdge_missing_source(t *testing.T) {
	g := New()

	n2, _ := g.AddNode("node-2", "person")
	actual, err := g.AddEdge("edge-1234", "nissing", "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestAddEdge_missing_target(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	actual, err := g.AddEdge("edge-1234", n1.UID, "knows", "missing", KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}
