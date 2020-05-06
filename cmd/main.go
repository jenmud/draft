package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jenmud/draft/graph"
	pb "github.com/jenmud/draft/service"
	micro "github.com/micro/go-micro/v2"
	microConfig "github.com/micro/go-micro/v2/config"
	microEnv "github.com/micro/go-micro/v2/config/source/env"
	microFlag "github.com/micro/go-micro/v2/config/source/flag"
	"google.golang.org/protobuf/proto"
)

var (
	version = micro.Version("v0.0.0")
	store   *graph.Graph
	config  microConfig.Config
)

func parseArgs() {
	var err error

	config, err = microConfig.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	flag.String("addr", ":", "Address to accept client connections on")
	flag.String("name", "draft.srv", "Service name")
	flag.String("dump", "", "Load a dump (.draft) file")
	flag.Parse()

	err = config.Load(
		microEnv.NewSource(microEnv.WithPrefix("DRAFT")),
		microFlag.NewSource(microFlag.IncludeUnset(true)),
	)

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	store = graph.New()
	parseArgs()

	dump := config.Get("dump").String("")
	if dump != "" {
		log.Printf("Loading from %s", dump)
		data, err := ioutil.ReadFile(dump)
		if err != nil {
			log.Fatal(err)
		}

		dump := pb.DumpResp{}
		if err := proto.Unmarshal(data, &dump); err != nil {
			log.Fatal(err)
		}

		if err := load(store, dump); err != nil {
			log.Fatal(err)
		}
	}
}

// load a dump into the graph.
func load(g *graph.Graph, dump pb.DumpResp) error {
	start := time.Now()

	for _, node := range dump.Nodes {
		if _, err := g.AddNode(node.Uid, node.Label, convertServicePropsToGraphKVs(node.Properties)...); err != nil {
			return fmt.Errorf("[load] %s", err)
		}
	}

	for _, edge := range dump.Edges {
		if _, err := g.AddEdge(edge.Uid, edge.SourceUid, edge.Label, edge.TargetUid, convertServicePropsToGraphKVs(edge.Properties)...); err != nil {
			return fmt.Errorf("[load] %s", err)
		}
	}

	log.Printf("Loaded %d nodes and %d edges in %s", g.NodeCount(), g.EdgeCount(), time.Now().Sub(start))
	return nil
}

// run start the RPC service.
func run() error {
	mservice := micro.NewService(
		version,
		micro.Name(config.Get("name").String("draft.srv")),
		micro.Address(config.Get("addr").String("0.0.0.0:")),
	)

	pb.RegisterGraphHandler(mservice.Server(), &server{graph: store})
	return mservice.Run()

	// c := make(chan os.Signal, 1)

	// go func() {
	// 	<-c
	// 	var b bytes.Buffer
	// 	server.Save(&b)
	// 	ioutil.WriteFile("../web/example/dump.draft", b.Bytes(), 0644)
	// }()

	// signal.Notify(c, os.Interrupt)
}

// main is the main entrypoint.
func main() {
	log.Fatal(run())
}
