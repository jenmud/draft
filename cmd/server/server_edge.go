package main

import (
	"context"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

func (s *server) AddEdge(ctx context.Context, req *pb.EdgeReq) (*pb.EdgeResp, error) {
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

	edge, err := s.graph.AddEdge(req.Uid, req.SourceUid, req.Label, req.TargetUid, kvs...)
	if err != nil {
		return nil, err
	}

	resp := pb.EdgeResp{
		Uid:        edge.UID,
		SourceUid:  edge.SourceUID,
		Label:      edge.Label,
		TargetUid:  edge.TargetUID,
		Properties: convertGraphPropsToServiceProps(edge.Properties),
	}

	return &resp, nil
}

func (s *server) RemoveEdge(ctx context.Context, req *pb.UIDReq) (*pb.RemoveResp, error) {
	panic("NotImplemented")
}

func (s *server) Edge(ctx context.Context, req *pb.UIDReq) (*pb.EdgeResp, error) {
	panic("NotImplemented")
}

func (s *server) Edges(req *pb.EdgesReq, stream pb.Graph_EdgesServer) error {
	panic("NotImplemented")
}

func (s *server) Dump(ctx context.Context, req *pb.DumpReq) (*pb.DumpResp, error) {
	panic("NotImplemented")
}
