package main

import (
	"context"
	"fmt"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func (s *server) AddNode(ctx context.Context, req *pb.NodeReq, resp *pb.NodeResp) error {
	kvs := make([]graph.KV, len(req.Properties))

	count := 0
	for k, v := range req.Properties {
		kvs[count] = graph.KV{Key: k, Value: v}
		count++
	}

	node, err := s.graph.AddNode(req.Uid, req.Label, kvs...)
	if err != nil {
		return fmt.Errorf("[AddNode] Error adding node: %v", err)
	}

	resp = &pb.NodeResp{
		Uid:        node.UID,
		Label:      node.Label,
		Properties: node.Properties,
		InEdges:    node.InEdges(),
		OutEdges:   node.OutEdges(),
	}

	return nil
}

func (s *server) RemoveNode(ctx context.Context, req *pb.UIDReq, resp *pb.RemoveResp) error {
	resp = &pb.RemoveResp{Uid: req.Uid}

	if err := s.graph.RemoveNode(req.Uid); err != nil {
		resp.Error = err.Error()
		resp.Success = false
		return fmt.Errorf("[RemoveNode] Error removing node: %v", err)
	}

	resp.Success = true // if we got to this point then everything was successful
	return nil
}

func (s *server) Node(ctx context.Context, req *pb.UIDReq, resp *pb.NodeResp) error {
	node, err := s.graph.Node(req.Uid)
	if err != nil {
		return fmt.Errorf("[Node] Error fetching node: %v", err)
	}

	resp = &pb.NodeResp{
		Uid:        node.UID,
		Label:      node.Label,
		Properties: node.Properties,
		InEdges:    node.InEdges(),
		OutEdges:   node.OutEdges(),
	}

	return nil
}

func (s *server) Nodes(req *pb.NodesReq, stream pb.Graph_NodesStream) error {
	iter := s.graph.Nodes()
	for iter.Next() {
		node := iter.Value().(graph.Node)

		resp := pb.NodeResp{
			Uid:        node.UID,
			Label:      node.Label,
			Properties: node.Properties,
			InEdges:    node.InEdges(),
			OutEdges:   node.OutEdges(),
		}

		if err := stream.Send(&resp); err != nil {
			return fmt.Errorf("[Nodes] Error fetching and streaming nodes: %v", err)
		}
	}

	return nil
}
