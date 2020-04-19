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
	node, err := s.graph.GetNode(req.Uid)
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
