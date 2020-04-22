package main

import (
	"context"
	"io"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
	"google.golang.org/protobuf/proto"
)

func convertGraphPropsToServiceProps(props map[string]graph.Value) map[string]*pb.Value {
	p := make(map[string]*pb.Value)

	for k, v := range props {
		p[k] = &pb.Value{Type: v.Type, Value: v.Value}
	}

	return p
}

func convertServicePropsToGraphKVs(props map[string]*pb.Value) []graph.KV {
	kvs := make([]graph.KV, len(props))

	count := 0
	for k, v := range props {
		kv := graph.KV{Key: k, Value: graph.Value{Type: v.Type, Value: v.Value}}
		kvs[count] = kv
		count++
	}

	return kvs
}

type server struct {
	graph *graph.Graph
}

func (s *server) Stats(ctx context.Context, req *pb.StatsReq) (*pb.StatsResp, error) {
	stats := s.graph.Stats()
	resp := &pb.StatsResp{
		NumCpu:           int32(stats.NumCPU),
		NodeCount:        int32(stats.NodeCount),
		EdgeCount:        int32(stats.EdgeCount),
		StartTime:        stats.StartTime.String(),
		NumGoroutines:    int32(stats.NumGoroutings),
		TotalMemoryAlloc: int32(stats.MemStats.TotalAlloc),
	}

	return resp, nil
}

// Save the current graph.
func (s *server) Save(w io.Writer) error {
	dump, err := s.Dump(context.Background(), &pb.DumpReq{})
	if err != nil {
		return err
	}

	output, err := proto.Marshal(dump)
	if err != nil {
		return err
	}

	_, err = w.Write(output)
	return err
}

func (s *server) Dump(ctx context.Context, req *pb.DumpReq) (*pb.DumpResp, error) {
	// TODO: add in the subgraph and levels
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
