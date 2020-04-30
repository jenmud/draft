package graph

import (
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
func filterByLabels(labels []string, nodes <-chan Node, out chan<- Node) {
	for node := range nodes {
		if len(labels) == 0 {
			out <- node
			continue
		}

	labelLoop:
		for _, label := range labels {
			if node.Label == label {
				out <- node
				break labelLoop
			}
		}
	}
}

// filterByProperties filters nodes which contain the given labels.
//func filterByLabels(labels []string, nodes <-chan Node, out chan<- Node) {
//	for node := range nodes {
//		if len(labels) == 0 {
//			out <- node
//			continue
//		}
//
//	labelLoop:
//		for _, label := range labels {
//			if node.Label == label {
//				out <- node
//				break labelLoop
//			}
//		}
//	}
//}

// nodeMapper maps iterItems to graph nodes.
func nodeMapper(nodes <-chan interface{}, out chan<- Node) {
	for node := range nodes {
		out <- node.(Node)
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
	final := make(chan Node, g.NodeCount())

	// search for nodes
	for _, rc := range plan.ReadingClause {
		for _, match := range rc.Matches {
			nodes := make(chan Node, g.NodeCount())
			nodeMapper(g.Nodes().Channel(), nodes)
			for _, node := range match.Nodes {
				filterByLabels(node.Labels, nodes, final)
			}
		}
	}

	close(final)
	for node := range final {
		if _, err := subg.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
			return subg, err
		}
	}

	return subg, nil
}
