package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/astranet/meshRPC/cluster"
	"github.com/gin-gonic/gin"
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
		Value:  "0.0.0.0:11999",
	})
	httpListenHost = app.String(cli.StringOpt{
		Name:   "H http-host",
		Desc:   "Specify listen HTTP host.",
		EnvVar: "APP_HTTP_HOST",
		Value:  "0.0.0.0",
	})
	httpListenPort = app.String(cli.StringOpt{
		Name:   "P http-port",
		Desc:   "Specify listen HTTP port.",
		EnvVar: "APP_HTTP_PORT",
		Value:  "8282",
	})
)

var app = cli.App("mesh_api", "An example API Gateway for meshRPC cluster.")

func main() {
	app.Action = func() {
		// Init a cluster client
		c := cluster.NewAstraCluster("mesh_api", &cluster.AstraOptions{
			Tags: []string{
				*clusterName,
			},
			Nodes: *clusterNodes,
			// Debug: true,
		})
		// Listen on a TCP address, this address can be used
		// by other peers to discover each other in this cluster.
		if err := c.ListenAndServe(*netAddr); err != nil {
			closer.Fatalln(err)
		}

		// Start a normal Gin HTTP server that will use cluster endpoints.
		httpListenAndServe(c)
	}
	app.Run(os.Args)
}

func httpListenAndServe(c cluster.Cluster) {
	// Init default Gin router
	router := gin.Default()
	// Wait until greeter service dependency is available.
	wait(c, map[string]cluster.HandlerSpec{
		"greeter": greeter.RPCHandlerSpec,
	})
	// Init a new meshRPC client to the greeter service.
	greeterClient := c.NewClient("greeter", greeter.RPCHandlerSpec)
	// A greeter.ServiceClient instance successfully conforms the greeter.Service interface
	// and may be used in place of the local greeter.Service instance.
	var svc greeter.Service = greeter.NewServiceClient(greeterClient, nil)

	// Set an endpoint handler for the Greet function.
	// Example Request:
	// $ curl http://localhost:8282/greeter/greet/Max
	router.GET("/greeter/greet/:name", func(c *gin.Context) {
		// Service call is actually done over meshRPC...
		message, err := svc.Greet(c.Param("name"))
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, message)
		return
	})

	// Set an endpoint handler for the SendPostcard function.
	// Example Request:
	// $ curl -X POST http://localhost:8282/greeter/sendPostcard/Max/World/Hello
	router.POST("/greeter/sendPostcard/:recipient/:address/:message", func(c *gin.Context) {
		// Fill the object fields with supplied params.
		postcard := &greeter.Postcard{
			Recipient: c.Param("recipient"),
			Address:   c.Param("address"),
			Message:   c.Param("message"),
		}
		// Service call is actually done over meshRPC...
		err := svc.SendPostcard(postcard)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
		c.JSON(200, postcard)
		return
	})

	// Bonus! Simply expose you custom HTTP handlers with "Use" function.
	// Note that we init another HTTP client to greeter.HandlerSpec instead
	// of greeter.RPCHandlerSpec, because they might have different HTTP verbs
	// and permissions attached.
	//
	// Example Request:
	// $ curl http://localhost:8282/greeter/check
	greeterClient2 := c.NewClient("greeter", greeter.HandlerSpec)
	router.GET("/greeter/check", gin.WrapH(greeterClient2.Use("Check")))
	// And meshRPC endpoints too!
	// Example Request:
	// $ curl -d'{"name": "Max"}' http://localhost:8282/greeter/greet
	router.POST("/greeter/greet", gin.WrapH(greeterClient.Use("Greet")))

	listenAddr := *httpListenHost + ":" + *httpListenPort
	if err := router.Run(listenAddr); err != nil {
		log.Fatalln(err)
	}
}

func wait(c cluster.Cluster, specs map[string]cluster.HandlerSpec) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	if err := c.Wait(ctx, specs); err != nil {
		log.Println("mesh_rpc: service await failure:", err)
	}
}
