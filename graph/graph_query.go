package graph

import (
	"log"

	"github.com/jenmud/draft/graph/parser/cypher"
)

// Eval is a evaluator funtion
type Eval func(Node) bool

// reducer is a mapping function for pulling out nodes that evaluate to true.
func reducer(in <-chan Node, out chan<- Node, evaluator Eval) {
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

// filterByLables filters nodes which contain the given labels.
func filterByLabels(labels []string, g *Graph, nodes Iterator) {
	if len(labels) == 0 {
		for nodes.Next() {
			node := nodes.Value().(Node)
			if _, err := g.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
				log.Printf("[filterByLabels] %s", err)
			}
		}
	}

	for nodes.Next() {
		node := nodes.Value().(Node)
	labelLoop:
		for _, label := range labels {
			if node.Label == label {
				if _, err := g.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
					log.Printf("[filterByLabels] %s", err)
				}
				break labelLoop
			}
		}
	}
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

	// search for nodes
	for _, rc := range plan.ReadingClause {
		for _, match := range rc.Matches {
			nodes := g.Nodes()
			for _, node := range match.Nodes {
				filterByLabels(node.Labels, subg, nodes)
			}
		}
	}

	return subg, nil
}
