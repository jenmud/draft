package graph

import (
	"fmt"
	"log"

	"github.com/jenmud/draft/graph/parser/cypher"
)

// Query takes a query string and returns a subgraph containing
// the query results.
func (g *Graph) Query(query string) (*Graph, error) {
	subg := New()

	queryResult, err := cypher.Parse("", []byte(query))
	if err != nil {
		return nil, err
	}

	// search for nodes
	for _, rc := range queryResult.(cypher.QueryPlan).ReadingClause {
		for _, match := range rc.Matches {
			for _, node := range match.Nodes {
				nodes := g.NodesBy(node.Labels, node.Properties)
				for nodes.Next() {
					node := nodes.Value().(Node)
					if _, err := subg.AddNode(node.UID, node.Label, convertPropertiesToKV(node.Properties)...); err != nil {
						log.Printf("[Query] %v", err)
					}

					// query the edges and add the attached in and out bound nodes.
					// ()-->(node)
					for _, edgeUID := range node.InEdges() {
						edge, err := g.Edge(edgeUID)
						if err != nil {
							return subg, fmt.Errorf("[Query] Error populating edges: %v", err)
						}

						sourceNode, err := g.Node(edge.SourceUID)
						if err != nil {
							return subg, fmt.Errorf("[Query] Error fetching inbound node: %v", err)
						}

						if _, err := subg.AddNode(sourceNode.UID, sourceNode.Label, convertPropertiesToKV(sourceNode.Properties)...); err != nil {
							log.Printf("[Query] Error inserting source node: %v", err)
						}

						if _, err := subg.AddEdge(edge.UID, sourceNode.UID, edge.Label, node.UID, convertPropertiesToKV(edge.Properties)...); err != nil {
							log.Printf("[Query] Error inbound edge: %v", err)
						}
					}

					// (node)-->()
					for _, edgeUID := range node.OutEdges() {
						edge, err := g.Edge(edgeUID)
						if err != nil {
							return subg, fmt.Errorf("[Query] Error populating edges: %v", err)
						}

						targetNode, err := g.Node(edge.TargetUID)
						if err != nil {
							return subg, fmt.Errorf("[Query] Error fetching outbound node: %v", err)
						}

						if _, err := subg.AddNode(targetNode.UID, targetNode.Label, convertPropertiesToKV(targetNode.Properties)...); err != nil {
							log.Printf("[Query] Error inserting target node: %v", err)
						}

						if _, err := subg.AddEdge(edge.UID, node.UID, edge.Label, targetNode.UID, convertPropertiesToKV(edge.Properties)...); err != nil {
							log.Printf("[Query] Error outbound edge: %v", err)
						}
					}
				}
			}
		}
	}

	return subg, nil
}
