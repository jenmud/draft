package main

import (
	"context"
	"fmt"
	"io"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
	"google.golang.org/protobuf/proto"
)

func convertServicePropsToGraphKVs(props map[string][]byte) []graph.KV {
	kvs := make([]graph.KV, len(props))

	count := 0
	for k, v := range props {
		kv := graph.KV{Key: k, Value: v}
		kvs[count] = kv
		count++
	}

	return kvs
}

type server struct {
	graph *graph.Graph
}

func (s *server) Stats(ctx context.Context, req *pb.StatsReq, resp *pb.StatsResp) error {
	stats := s.graph.Stats()
	resp.NumCpu = int32(stats.NumCPU)
	resp.NodeCount = int32(stats.NodeCount)
	resp.EdgeCount = int32(stats.EdgeCount)
	resp.StartTime = stats.StartTime.String()
	resp.NumGoroutines = int32(stats.NumGoroutings)
	resp.TotalMemoryAlloc = int32(stats.MemStats.TotalAlloc)
	return nil
}

// Save the current graph.
func (s *server) Save(w io.Writer) error {
	resp := &pb.DumpResp{}
	if err := s.Dump(context.Background(), &pb.DumpReq{}, resp); err != nil {
		return err
	}

	output, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = w.Write(output)
	return err
}

func dump(g *graph.Graph, resp *pb.DumpResp) error {
	// TODO: add in the subgraph and levels
	nodesIter := g.Nodes()
	edgesIter := g.Edges()

	resp.Nodes = make([]*pb.NodeResp, nodesIter.Size())
	resp.Edges = make([]*pb.EdgeResp, edgesIter.Size())

	ncount := 0
	for nodesIter.Next() {
		node := nodesIter.Value().(graph.Node)
		nresp := &pb.NodeResp{
			Uid:        node.UID,
			Label:      node.Label,
			Properties: node.Properties,
			InEdges:    node.InEdges(),
			OutEdges:   node.OutEdges(),
		}
		resp.Nodes[ncount] = nresp
		ncount++
	}

	ecount := 0
	for edgesIter.Next() {
		edge := edgesIter.Value().(graph.Edge)
		eresp := &pb.EdgeResp{
			Uid:        edge.UID,
			SourceUid:  edge.SourceUID,
			Label:      edge.Label,
			TargetUid:  edge.TargetUID,
			Properties: edge.Properties,
		}
		resp.Edges[ecount] = eresp
		ecount++
	}

	return nil
}

func (s *server) Query(ctx context.Context, req *pb.QueryReq, resp *pb.DumpResp) error {
	g, err := s.graph.Query(req.Query)
	if err != nil {
		return fmt.Errorf("[Query] Error trying to execute a query: %v", err)
	}

	if err := dump(g, resp); err != nil {
		return fmt.Errorf("[Query] Error trying to dump query response: %v", err)
	}

	return nil
}

func (s *server) Dump(ctx context.Context, req *pb.DumpReq, resp *pb.DumpResp) error {
	if err := dump(s.graph, resp); err != nil {
		return fmt.Errorf("[Dump] Error trying to dump the graph: %v", err)
	}

	return nil
}

// See server_node.go for node methods
// See server_edge.go for edge methods
