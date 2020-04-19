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

func TestAddEdge_duplicate_uid(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})
	actual, err := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestHasEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	assert.Equal(t, true, g.HasEdge("edge-1234"))
}

func TestHasEdge_missing(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	assert.Equal(t, false, g.HasEdge("missing"))
}

func TestRemoveEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})

	g.RemoveEdge("edge-1234")
	assert.Equal(t, false, g.HasEdge("edge-1234"))
}

func TestEdge(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")
	expected, _ := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID, KV{Key: "since", Value: Value{Type: "string", Value: []byte("school")}})
	actual, err := g.Edge("edge-1234")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEdge_no_such_edge(t *testing.T) {
	g := New()

	actual, err := g.Edge("edge-missing")

	assert.NotNil(t, err)
	assert.Equal(t, Edge{}, actual)
}

func TestEdges(t *testing.T) {
	g := New()

	n1, _ := g.AddNode("node-1", "person")
	n2, _ := g.AddNode("node-2", "person")

	e1, _ := g.AddEdge("edge-1234", n1.UID, "knows", n2.UID)
	e2, _ := g.AddEdge("edge-2345", n1.UID, "knows", n2.UID)
	e3, _ := g.AddEdge("edge-3456", n1.UID, "knows", n2.UID)

	expected := []Edge{e1, e2, e3}
	actual := []Edge{}

	iter := g.Edges()
	for iter.Next() {
		edge := iter.Value().(Edge)
		actual = append(actual, edge)
	}

	assert.ElementsMatch(t, expected, actual)
}
