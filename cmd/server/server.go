package main

import (
	"context"

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

func (s *server) Dump(ctx context.Context, req *pb.DumpReq) (*pb.DumpResp, error) {
	nodesIter := s.graph.Nodes()
	edgesIter := s.graph.Edges()

	dump := &pb.DumpResp{
		Nodes: make([]*pb.NodeResp, nodesIter.Size()),
		Edges: make([]*pb.EdgeResp, edgesIter.Size()),
	}

	ncount := 0
	for nodesIter.Next() {
		node := nodesIter.Value().(graph.Node)
		resp := &pb.NodeResp{
			Uid:        node.UID,
			Label:      node.Label,
			Properties: convertGraphPropsToServiceProps(node.Properties),
		}
		dump.Nodes[ncount] = resp
		ncount++
	}

	ecount := 0
	for edgesIter.Next() {
		edge := edgesIter.Value().(graph.Edge)
		resp := &pb.EdgeResp{
			Uid:        edge.UID,
			SourceUid:  edge.SourceUID,
			Label:      edge.Label,
			TargetUid:  edge.TargetUID,
			Properties: convertGraphPropsToServiceProps(edge.Properties),
		}
		dump.Edges[ecount] = resp
		ecount++
	}

	return dump, nil
}

// See server_node.go for node methods
// See server_edge.go for edge methods
