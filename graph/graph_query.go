package graph

import (
	"sync"

	"github.com/jenmud/draft/graph/parser/cypher"
)

// Eval is a evaluator funtion
type Eval func(Node) bool

// reducer is a mapping function for pulling out nodes that evaluate to true.
func reducer(wg *sync.WaitGroup, in <-chan Node, out chan<- Node, evaluator Eval) {
	defer wg.Done()
	for node := range in {
		if evaluator(node) {
			out <- node
		}
	}
}

// mapper takes a nodes iterator and converts the items into nodes
// and pushes then into the out channel for processing.
func mapper(nodes Iterator, out chan<- Node) {
	for nodes.Next() {
		out <- nodes.Value().(Node)
	}
	close(out)
}

// Query takes a query string and returns a subgraph containing
// the query results.
func (g *Graph) Query(query string) (*Graph, error) {
	subg := New()

	queryResult, err := cypher.Parse("", []byte(query))
	if err != nil {
		return nil, err
	}

	plan := queryResult.(cypher.QueryPlan)

	nodesIter := g.Nodes()
	nodes := []Node{}

	// search for nodes
	for nodesIter.Next() {
		for _, rc := range plan.ReadingClause {
			for _, node := range rc.Match.Nodes {
				n := nodesIter.Value().(Node)
				for _, label := range node.Labels {
					if n.Label == label {
						nodes = append(nodes, n)
					}
				}
			}
		}
	}

	for _, node := range nodes {
		if _, err := subg.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
			return subg, err
		}
	}

	return subg, nil
}
