package graph

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	type TestCase struct {
		g           *Graph
		Query       string
		Expected    []Node
		Name        string
		ShouldError bool
	}

	reader := bytes.NewReader(readTestData(t, "simple-graph.json"))

	g, err := NewFromJSON(reader)
	if err != nil {
		t.Fatal(err)
	}

	// Add in some extra nodes to make the testing below easier
	g.AddNode("node-flower", "flower")
	g.AddNode("node-flower-rose", "flower", KV{Key: "name", Value: []byte("rose")})
	g.AddNode("node-car", "car")

	tests := []TestCase{
		TestCase{
			g:     g,
			Name:  "NoLabelsExpectAllNodes",
			Query: `MATCH (n) RETURN n`,
			Expected: []Node{
				Node{
					UID:        "node-car",
					Label:      "car",
					Properties: map[string][]byte{},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-flower",
					Label:      "flower",
					Properties: map[string][]byte{},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-flower-rose",
					Label:      "flower",
					Properties: map[string][]byte{"name": []byte("rose")},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-foo",
					Label:      "person",
					Properties: map[string][]byte{"name": []byte("foo")},
					inEdges:    map[string]struct{}{},
					outEdges: map[string]struct{}{
						"edge-like":  {},
						"edge-knows": {},
						"edge-owns":  {},
					},
				},
				Node{
					UID:        "node-bar",
					Label:      "person",
					Properties: map[string][]byte{"name": []byte("bar")},
					inEdges: map[string]struct{}{
						"edge-like":  {},
						"edge-knows": {},
					},
					outEdges: map[string]struct{}{
						"edge-dislike": {},
					},
				},
				Node{
					UID:        "node-dog",
					Label:      "animal",
					Properties: map[string][]byte{"name": []byte("socks")},
					inEdges: map[string]struct{}{
						"edge-owns":    {},
						"edge-dislike": {},
					},
					outEdges: map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:     g,
			Name:  "EdgesAndEdgeNodesIncluded",
			Query: `MATCH (n {name: "socks"}) RETURN n`,
			Expected: []Node{
				Node{
					UID:        "node-foo",
					Label:      "person",
					Properties: map[string][]byte{"name": []byte("foo")},
					inEdges:    map[string]struct{}{},
					outEdges: map[string]struct{}{
						"edge-owns": {},
					},
				},
				Node{
					UID:        "node-bar",
					Label:      "person",
					Properties: map[string][]byte{"name": []byte("bar")},
					inEdges:    map[string]struct{}{},
					outEdges: map[string]struct{}{
						"edge-dislike": {},
					},
				},
				Node{
					UID:        "node-dog",
					Label:      "animal",
					Properties: map[string][]byte{"name": []byte("socks")},
					inEdges: map[string]struct{}{
						"edge-owns":    {},
						"edge-dislike": {},
					},
					outEdges: map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:     g,
			Name:  "OnlyShouldContainFlowerNodes",
			Query: `MATCH (n:flower) RETURN n`,
			Expected: []Node{
				Node{
					UID:        "node-flower",
					Label:      "flower",
					Properties: map[string][]byte{},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-flower-rose",
					Label:      "flower",
					Properties: map[string][]byte{"name": []byte("rose")},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:    g,
			Name: "MultiMatchByLabel",
			Query: `
			MATCH (n:car)
			MATCH (m:flower)
			RETURN n, m`,
			Expected: []Node{
				Node{
					UID:        "node-flower",
					Label:      "flower",
					Properties: map[string][]byte{},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-flower-rose",
					Label:      "flower",
					Properties: map[string][]byte{"name": []byte("rose")},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-car",
					Label:      "car",
					Properties: map[string][]byte{},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:    g,
			Name: "MultiMatchByLabelAndProperty",
			Query: `
			MATCH (n:flower {name: "rose"})
			MATCH (m:car)
			RETURN n, m`,
			Expected: []Node{
				Node{
					UID:        "node-flower-rose",
					Label:      "flower",
					Properties: map[string][]byte{"name": []byte("rose")},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
				Node{
					UID:        "node-car",
					Label:      "car",
					Properties: map[string][]byte{},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:           g,
			Name:        "MultipleLablesNotSupported",
			Query:       `MATCH (n:animal:person) RETURN n`,
			Expected:    []Node{},
			ShouldError: true,
		},
	}

	for _, test := range tests {
		subg, err := test.g.Query(test.Query)
		if test.ShouldError {
			assert.NotNil(t, err, "%s query expected to fail: %s", test.Name)
			continue
		} else {
			assert.Nil(t, err, "%s did not expect a error but got: %s", test.Name, err)
		}

		nodes := subg.Nodes()
		actual := make([]Node, nodes.Size())

		count := 0
		for nodes.Next() {
			actual[count] = nodes.Value().(Node)
			count++
		}

		assert.ElementsMatch(t, test.Expected, actual, "%s expected %#v but got %#v", test.Name, test.Expected, actual)
	}

}
