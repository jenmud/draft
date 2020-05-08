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

	resp.Uid = node.UID
	resp.Label = node.Label
	resp.Properties = node.Properties
	resp.InEdges = node.InEdges()
	resp.OutEdges = node.OutEdges()

	return nil
}

func (s *server) RemoveNode(ctx context.Context, req *pb.UIDReq, resp *pb.RemoveResp) error {
	resp.Uid = req.Uid

	if err := s.graph.RemoveNode(req.Uid); err != nil {
		resp.Error = err.Error()
		resp.Success = false
		return fmt.Errorf("[RemoveNode] Error removing node: %v", err)
	}

	resp.Success = true // if we got to this point then everything was successful
	return nil
}

func (s *server) Node(ctx context.Context, req *pb.NodeReq, resp *pb.NodeResp) error {
	// if we have a uid in the node request, then just use that.
	if req.Uid != "" {
		node, err := s.graph.Node(req.Uid)
		if err != nil {
			return fmt.Errorf("[Node] Error fetching node: %v", err)
		}

		resp.Uid = node.UID
		resp.Label = node.Label
		resp.Properties = node.Properties
		resp.InEdges = node.InEdges()
		resp.OutEdges = node.OutEdges()
		return nil
	}

	// if we don't have a Uid do a filter for labels and properties.
	iter := s.graph.NodesBy([]string{req.Label}, req.Properties)
	if iter.Size() != 1 {
		return fmt.Errorf("[Node] Error fetching node, expected 1 but found %d", iter.Size())
	}

	node := iter.Value().(graph.Node)
	resp.Uid = node.UID
	resp.Label = node.Label
	resp.Properties = node.Properties
	resp.InEdges = node.InEdges()
	resp.OutEdges = node.OutEdges()

	return nil
}

func (s *server) Nodes(ctx context.Context, req *pb.NodesReq, stream pb.Graph_NodesStream) error {
	iter := s.graph.NodesBy(req.Label, req.Properties)
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
