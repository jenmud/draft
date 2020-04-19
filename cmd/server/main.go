package main

import (
	"flag"
	"log"
	"net"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
	"google.golang.org/grpc"
)

var (
	addr = ":8000"
)

func init() {
	flag.StringVar(&addr, "addr", addr, "Address and port to service and accept client connections.")
	flag.Parse()
}

// run start the RPC service.
func run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterGraphServer(s, &server{graph: graph.New()})
	return s.Serve(listener)
}

// main is the main entrypoint.
func main() {
	log.Fatal(run(addr))
}
