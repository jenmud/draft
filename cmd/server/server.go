package main

import (
	"context"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
)

type server struct {
	graph *graph.Graph
}

func (s *server) AddNode(ctx context.Context, req *pb.NodeReq) (*pb.NodeResp, error) {
	panic("NotImplemented")
}

func (s *server) RemoveNode(ctx context.Context, req *pb.UIDReq) (*pb.RemoveResp, error) {
	panic("NotImplemented")
}

func (s *server) Node(ctx context.Context, req *pb.UIDReq) (*pb.NodeResp, error) {
	panic("NotImplemented")
}

func (s *server) Nodes(req *pb.NodesReq, stream pb.Graph_NodesServer) error {
	panic("NotImplemented")
}

func (s *server) AddEdge(ctx context.Context, req *pb.EdgeReq) (*pb.EdgeResp, error) {
	panic("NotImplemented")
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
