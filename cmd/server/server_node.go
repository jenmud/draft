package main

import (
	"context"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func (s *server) AddNode(ctx context.Context, req *pb.NodeReq) (*pb.NodeResp, error) {
	kvs := make([]graph.KV, len(req.Properties))

	count := 0
	for k, v := range req.Properties {
		kv := graph.KV{
			Key: k,
			Value: graph.Value{
				Type:  v.Type,
				Value: v.Value,
			},
		}

		kvs[count] = kv
		count++
	}

	node, err := s.graph.AddNode(req.Uid, req.Label, kvs...)
	if err != nil {
		return nil, err
	}

	resp := pb.NodeResp{
		Uid:        node.UID,
		Label:      node.Label,
		Properties: convertGraphPropsToServiceProps(node.Properties),
	}

	return &resp, nil
}

func (s *server) RemoveNode(ctx context.Context, req *pb.UIDReq) (*pb.RemoveResp, error) {
	s.graph.RemoveNode(req.Uid)
	return &pb.RemoveResp{Uid: req.Uid, Success: true}, nil
}

func (s *server) Node(ctx context.Context, req *pb.UIDReq) (*pb.NodeResp, error) {
	node, err := s.graph.Node(req.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.NodeResp{Uid: node.UID, Label: node.Label, Properties: convertGraphPropsToServiceProps(node.Properties)}, nil
}

func (s *server) Nodes(req *pb.NodesReq, stream pb.Graph_NodesServer) error {
	iter := s.graph.Nodes()
	for iter.Next() {
		node := iter.Value().(graph.Node)
		resp := pb.NodeResp{Uid: node.UID, Label: node.Label, Properties: convertGraphPropsToServiceProps(node.Properties)}
		if err := stream.Send(&resp); err != nil {
			return nil
		}
	}

	return nil
}
