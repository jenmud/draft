package main

import (
	"context"
	"fmt"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func (s *server) AddEdge(ctx context.Context, req *pb.EdgeReq, resp *pb.EdgeResp) error {
	kvs := make([]graph.KV, len(req.Properties))

	count := 0
	for k, v := range req.Properties {
		kvs[count] = graph.KV{Key: k, Value: v}
		count++
	}

	edge, err := s.graph.AddEdge(req.Uid, req.SourceUid, req.Label, req.TargetUid, kvs...)
	if err != nil {
		return fmt.Errorf("[AddEdge] Error adding edge: %v", err)
	}

	resp = &pb.EdgeResp{
		Uid:        edge.UID,
		SourceUid:  edge.SourceUID,
		Label:      edge.Label,
		TargetUid:  edge.TargetUID,
		Properties: edge.Properties,
	}

	return nil
}

func (s *server) RemoveEdge(ctx context.Context, req *pb.UIDReq, resp *pb.RemoveResp) error {
	resp = &pb.RemoveResp{Uid: req.Uid}

	if err := s.graph.RemoveEdge(req.Uid); err != nil {
		resp.Error = err.Error()
		return fmt.Errorf("[RemoveEdge] Error removing edge: %v", err)
	}

	resp.Success = true // at this point every was successful
	return nil
}

func (s *server) Edge(ctx context.Context, req *pb.UIDReq, resp *pb.EdgeResp) error {
	edge, err := s.graph.Edge(req.Uid)
	if err != nil {
		return fmt.Errorf("[Edge] Error fetching edge: %v", err)
	}

	resp = &pb.EdgeResp{
		Uid:        edge.UID,
		SourceUid:  edge.SourceUID,
		Label:      edge.Label,
		TargetUid:  edge.TargetUID,
		Properties: edge.Properties,
	}

	return nil
}

func (s *server) Edges(ctx context.Context, req *pb.EdgesReq, stream pb.Graph_EdgesStream) error {
	iter := s.graph.Edges()
	for iter.Next() {
		edge := iter.Value().(graph.Edge)

		resp := pb.EdgeResp{
			Uid:        edge.UID,
			SourceUid:  edge.SourceUID,
			Label:      edge.Label,
			TargetUid:  edge.TargetUID,
			Properties: edge.Properties,
		}

		if err := stream.Send(&resp); err != nil {
			return fmt.Errorf("[Edges] Error fetching and streaming edges: %v", err)
		}
	}

	return nil
}
