package main

import (
	"log"
	"os"

	"github.com/astranet/meshRPC/cluster"
	cli "github.com/jawher/mow.cli"
	"github.com/xlab/closer"

	greeter "github.com/astranet/meshRPC/example/greeter/service"
)

var (
	clusterNodes = app.Strings(cli.StringsOpt{
		Name:   "N nodes",
		Desc:   "A list of cluster nodes to join for service discovery.",
		EnvVar: "MESHRPC_CLUSTER_NODES",
		Value:  []string{},
	})
	clusterName = app.String(cli.StringOpt{
		Name:   "T tag",
		Desc:   "Cluster tag name.",
		EnvVar: "MESHRPC_CLUSTER_TAGNAME",
		Value:  "example",
	})
	netAddr = app.String(cli.StringOpt{
		Name:   "listen-addr",
		Desc:   "Listen address for cluster discovery and private networking.",
		EnvVar: "MESHRPC_LISTEN_ADDR",
		Value:  "0.0.0.0:0",
	})
)

var app = cli.App("greeter", "A Greeter service server for meshRPC cluster.")

func main() {
	app.Action = func() {
		defer closer.Close()

		// Init a new cluster connection
		c := cluster.NewAstraCluster("greeter", &cluster.AstraOptions{
			Tags: []string{
				*clusterName,
			},
			Nodes: *clusterNodes,
			// Debug: true,
		})
		// Init a new service instance (that is your code)
		service := greeter.NewService()
		// Init an RPC handler (that is the generated meshRPC code)
		meshRPC := greeter.NewRPCHandler(service, nil)
		// Publish RPC handlers on the cluster.
		c.Publish(meshRPC)

		// Bonus: publish your own HTTP handler too:
		handler := greeter.NewHandler()
		c.Publish(handler)
		// why not?

		// Start cluster connection and server
		if err := c.ListenAndServe(*netAddr); err != nil {
			log.Fatalln(err)
		}

		closer.Hold()
	}
	app.Run(os.Args)
}
