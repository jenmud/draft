package graph

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	reader := bytes.NewReader(readTestData(t, "simple-graph.json"))

	g, err := NewFromJSON(reader)
	if err != nil {
		t.Fatal(err)
	}

	subg, err := g.Query(`MATCH (n:person) RETURN n`)
	assert.Nil(t, err)

	expected := []Node{
		Node{
			UID:        "node-foo",
			Label:      "person",
			Properties: map[string][]byte{"name": []byte("foo")},
			inEdges:    map[string]struct{}{},
			outEdges:   map[string]struct{}{},
		},
		Node{
			UID:        "node-bar",
			Label:      "person",
			Properties: map[string][]byte{"name": []byte("bar")},
			inEdges:    map[string]struct{}{},
			outEdges:   map[string]struct{}{},
		},
	}

	nodes := subg.Nodes()
	actual := make([]Node, nodes.Size())

	count := 0
	for nodes.Next() {
		actual[count] = nodes.Value().(Node)
		count++
	}

	assert.ElementsMatch(t, expected, actual)
}
