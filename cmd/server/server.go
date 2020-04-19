package main

import (
	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func convertGraphPropsToServiceProps(props map[string]graph.Value) map[string]*pb.Value {
	p := make(map[string]*pb.Value)

	for k, v := range props {
		p[k] = &pb.Value{Type: v.Type, Value: v.Value}
	}

	return p
}

type server struct {
	graph *graph.Graph
}


// See server_node.go for node methods
// See server_edge.go for edge methods