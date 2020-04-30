package graph

import (
	"bytes"
	"testing"

	"github.com/jenmud/draft/graph/iterator"
	"github.com/stretchr/testify/assert"
)

func TestFilterByLabels(t *testing.T) {
	type TestCase struct {
		Labels   []string
		Nodes    []Node
		Expected []Node
		Name     string
	}

	n1 := NewNode("node-n1", "person")
	n2 := NewNode("node-n2", "car")
	n3 := NewNode("node-n3", "person")
	n4 := NewNode("node-n4", "bike")

	tests := []TestCase{
		TestCase{
			Labels:   []string{"person"},
			Nodes:    []Node{n1, n2, n3, n4},
			Expected: []Node{n1, n3},
			Name:     "SingleLabel",
		},
		TestCase{
			Labels:   []string{"person", "bike"},
			Nodes:    []Node{n1, n2, n3, n4},
			Expected: []Node{n1, n3, n4},
			Name:     "MultipleLabels",
		},
		TestCase{
			Labels:   []string{},
			Nodes:    []Node{n1, n2, n3, n4},
			Expected: []Node{n1, n2, n3, n4},
			Name:     "NoLabels",
		},
	}

	for _, test := range tests {
		nodes := make(chan Node, len(test.Nodes))
		out := make(chan Node, len(test.Nodes))

		for _, node := range test.Nodes {
			nodes <- node
		}

		close(nodes)

		filterByLabels(test.Labels, nodes, out)
		close(out)

		actual := []Node{}
		for node := range out {
			actual = append(actual, node)
		}

		assert.ElementsMatch(t, test.Expected, actual, "%s expected %v but got %v", test.Name, test.Expected, actual)
	}
}

func TestNodeMapper(t *testing.T) {
	n1 := NewNode("node-n1", "person")
	n2 := NewNode("node-n2", "car")
	iter := iterator.New([]interface{}{n1, n2})
	out := make(chan Node, 2)
	nodeMapper(iter.Channel(), out)
	assert.Equal(t, n1, <-out)
	assert.Equal(t, n2, <-out)
}

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

	tests := []TestCase{
		TestCase{
			g:     g,
			Name:  "NoLabelsExpectAllNodes",
			Query: `MATCH (n) RETURN n`,
			Expected: []Node{
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
				Node{
					UID:        "node-dog",
					Label:      "animal",
					Properties: map[string][]byte{"name": []byte("socks")},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:     g,
			Name:  "OnlyShouldContainAnimalNodes",
			Query: `MATCH (n:animal) RETURN n`,
			Expected: []Node{
				Node{
					UID:        "node-dog",
					Label:      "animal",
					Properties: map[string][]byte{"name": []byte("socks")},
					inEdges:    map[string]struct{}{},
					outEdges:   map[string]struct{}{},
				},
			},
		},
		TestCase{
			g:    g,
			Name: "MultiMatch",
			Query: `
			MATCH (n:animal)
			MATCH (m:person)
			RETURN n, m`,
			Expected: []Node{
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
				Node{
					UID:        "node-dog",
					Label:      "animal",
					Properties: map[string][]byte{"name": []byte("socks")},
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

		assert.ElementsMatch(t, test.Expected, actual, "%s expected %v but got %v", test.Name, test.Expected, actual)
	}

}
