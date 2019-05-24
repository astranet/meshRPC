# Mesh RPC

Automatic Service Mesh generator for pure Go micro services, a humble alternative to gRPC! A service mesh is a dedicated infrastructure layer for managing service-to-service communication, including RPC over HTTP. The `meshRPC` tool paired with `cluster` package is able to transform any type of legacy Go service into a new stack ops dream.

Even for such legacy services that contain many layers inside a single process, this framework can be used to decouple things,
using interface substitution. Consider an interface A, then use this tool to generate a microservice that implements interface A, but once called, instead of interface A invocation, there will be an RPC call over network to the corresponding microservice.
In this way you can separate a big project by small pieces without hurting integrity (just adding a bit of network latency).

All generated microservices require zero-configuration and are load-balanced (round robin with sticky sessions) out of the box!

### Install

```
$ go get -u github.com/astranet/meshRPC
```

### Usage in Go

Create a service file like this one:

```go
package greeter

type Postcard struct {
    PictureURL string
    Address    string
    Recipient  string
    Message    string
}

//go:generate meshRPC expose -P greeter -y

type Service interface {
    Greet(name string) (string, error)
    SendPostcard(card *Postcard) error
}

// Then implement methods for the service.
// See https://github.com/astranet/meshRPC/tree/master/example/greeter/service/service.go

func NewService() Service {
    return &service{}
}

type service struct{}
```

Make sure to implement the methods, they will be exposed to other microservices in cluster soon. See an example [service.go](https://github.com/astranet/meshRPC/tree/master/example/greeter/service/service.go).

Now run `meshRPC expose` or using `go generate`, please note that when running manually, you must specify project dir and the target sources path as arguments.

```
$ meshRPC -R . expose -P greeter service/

Actions to be committed
├── [1]  dir [project]/service must exist
├── [2]  overwrite file [project]/service/handler_gen.go with 152 lines of content
└── [3]  overwrite file [project]/service/client_gen.go with 123 lines of content

Are you sure to apply these changes?: y
queue.go:44: Action#1: dir [project]/service must exist
queue.go:44: Action#2: overwrite file [project]/service/handler_gen.go with 152 lines of content
queue.go:44: Action#3: overwrite file [project]/service/client_gen.go with 123 lines of content
```

Service is complete! Let's create a simple server that will handle cluster connections.

### Example

```go
func() {
    // Init a new cluster connection
    c := cluster.NewAstraCluster("greeter", &cluster.AstraOptions{
        Tags: []string{
            *clusterName,
        },
        Nodes: *clusterNodes,
        Debug: true,
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
}
```

See the full example [main.go](https://github.com/astranet/meshRPC/tree/master/example/greeter/main.go) that just populates config values. It took just a few lines to expose the service over network, as well as your custom HTTP endpoints. But, the question is, how to access those endpoints? For debugging purposes it's easy to also use `c.ListenAndServeHTTP` that will start an usual HTTP server in the same process, so you could connect to `greeter` instance directly. But that's not robust, because you want to have an API Gateway with load balancing and authorization, am I right? :)

At this point users that would use gRPC usually put an [Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/grpc) instance that will automatically convert HTTP/2 gRPC protobufs into HTTP/1.1 JSONs. You're expected to write an Envoy config and an xDS discovery service or install [Istio](https://istio.io/docs/concepts/security/) with automatic sidecar injection, on your Kubernetes cluster. I wish you luck.

#### API Gateway for meshRPC

We'll create a simple server that will look like an HTTP server, but it will also be a cluster discovery endpoint. All service nodes you will start will simply connect to it, and register within the cluster. Actually, any service node may connect to any other service node to get into the mesh. MeshRPC initiates persistent TCP connections between peers, and all virtual connections and streams are multiplexed on these persistent conns. This is how [AstraNet](https://github.com/astranet/astranet) cluster works.

```go
func() {
    // Init a cluster client
    c := cluster.NewAstraCluster("mesh_api", &cluster.AstraOptions{
        Tags: []string{
            *clusterName,
        },
        Nodes: *clusterNodes,
        Debug: true,
    })
    go func() {
        // Listen on a TCP address, this address can be used
        // by other peers to discover each other in this cluster.
        if err := c.ListenAndServe(*netAddr); err != nil {
            closer.Fatalln(err)
        }
    }()

    // Start a normal Gin HTTP server that will use cluster endpoints.
    httpListenAndServe(c)
}

func httpListenAndServe(c cluster.Cluster) {
    // Init default Gin router
    router := gin.Default()
    // Wait until greeter service dependency is available.
    wait(c, greeter.RPCHandlerSpec)
    // Init a new meshRPC client to the greeter service.
    greeterClient := c.NewClient(greeter.RPCHandlerSpec)
    // A greeter.ServiceClient instance sucessfully conforms the greeter.Service interface
    // and may be used in place of the local greeter.Service instance.
    var svc greeter.Service = greeter.NewServiceClient(greeterClient, nil)

    // Set an endpoint handler for the Greet function.
    // Example Request:
    router.GET("/greeter/greet/:name", func(c *gin.Context) {
        // Service call is actually done over meshRPC...
        message, err := svc.Greet(c.Param("name"))

// See the rest in:
// https://github.com/astranet/meshRPC/tree/master/example/mesh_api/main.go
```

See the full example [main.go](https://github.com/astranet/meshRPC/tree/master/example/mesh_api/main.go) to get the idea.
It takes just a few lines more to create an API Gateway for the whole cluster. Make sure that `:11999` port is closed from ouside connections, but `:8282` is okay to be exposed, or connected to classic reverse proxies such as Nginx and Caddy for TLS.

### Example Run

At this moment we have:
* `github.com/astranet/meshRPC/example/greeter` that is a server executable exposing service to meshRPC cluster.
* `github.com/astranet/meshRPC/example/mesh_api` that is an API Gateway for all meshRPC services in this cluster.

Let's run the API Gateway we just created:
```
go install github.com/astranet/meshRPC/example/mesh_api
$ mesh_api

[GIN-debug] GET    /ping                     --> github.com/astranet/meshRPC/cluster.okLoopback.func1 (2 handlers)
[GIN-debug] GET    /__heartbeat__            --> github.com/astranet/meshRPC/cluster.okLoopback.func1 (2 handlers)
[GIN-debug] POST   /__error__                --> github.com/astranet/meshRPC/cluster.errLoopback.func1 (2 handlers)
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

# ...waits until all depencenies respond to ping...

# Then starts an HTTP web server that is our API Gateway.
# Make sure you run greeter in another tab immediately after starting mesh_api instance.

[GIN-debug] GET    /greeter/greet/:name      --> main.httpListenAndServe.func1 (3 handlers)
[GIN-debug] POST   /greeter/sendPostcard/:recipient/:address/:message --> main.httpListenAndServe.func2 (3 handlers)
[GIN-debug] GET    /greeter/check            --> github.com/gin-gonic/gin.WrapH.func1 (3 handlers)
[GIN-debug] POST   /greeter/greet            --> github.com/gin-gonic/gin.WrapH.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8282
```

Meanwhile in a separate tab:

```
$ go install github.com/astranet/meshRPC/example/greeter
$ greeter -N localhost:11999

[GIN-debug] GET    /ping                     --> github.com/astranet/meshRPC/cluster.okLoopback.func1 (2 handlers)
[GIN-debug] GET    /__heartbeat__            --> github.com/astranet/meshRPC/cluster.okLoopback.func1 (2 handlers)
[GIN-debug] POST   /__error__                --> github.com/astranet/meshRPC/cluster.errLoopback.func1 (2 handlers)
[GIN-debug] POST   /rpcHandler/Greet         --> github.com/astranet/meshRPC/cluster.(*astraCluster).Publish.func2 (2 handlers)
[GIN-debug] POST   /rpcHandler/SendPostcard  --> github.com/astranet/meshRPC/cluster.(*astraCluster).Publish.func2 (2 handlers)
[GIN-debug] GET    /handler/Check            --> github.com/astranet/meshRPC/cluster.(*astraCluster).Publish.func2 (2 handlers)

[GIN] 2019/05/24 - 16:03:10 | 200 |         2.1µs |    StVzkmMmyzEA | GET      /ping
```

Yay! There is a ping from the gateway that awaited our service.. You can revise your own methods exposed.

#### Using the API

```
$ curl http://localhost:8282/greeter/greet/Max
"Hello, Max"
```

I called this method 5 times, in `mesh_rpc` logs I got this:
```
[GIN] 2019/05/24 - 16:07:51 | 200 |     2.24422ms |             ::1 | GET      /greeter/greet/Max
[GIN] 2019/05/24 - 16:07:51 | 200 |     1.67816ms |             ::1 | GET      /greeter/greet/Max
[GIN] 2019/05/24 - 16:07:51 | 200 |     909.254µs |             ::1 | GET      /greeter/greet/Max
[GIN] 2019/05/24 - 16:07:52 | 200 |     1.29748ms |             ::1 | GET      /greeter/greet/Max
[GIN] 2019/05/24 - 16:07:52 | 200 |     910.618µs |             ::1 | GET      /greeter/greet/Max
```

In `greeter` logs I got this, correspondingly:

```
[GIN] 2019/05/24 - 16:07:51 | 200 |     101.876µs |    NWx4LgFlyRUA | POST     /rpcHandler/Greet
[GIN] 2019/05/24 - 16:07:51 | 200 |      62.109µs |    NWx4LgFlyRUA | POST     /rpcHandler/Greet
[GIN] 2019/05/24 - 16:07:52 | 200 |     147.524µs |    NWx4LgFlyRUA | POST     /rpcHandler/Greet
[GIN] 2019/05/24 - 16:07:52 | 200 |      61.845µs |    NWx4LgFlyRUA | POST     /rpcHandler/Greet
[GIN] 2019/05/24 - 16:07:53 | 200 |      63.341µs |    NWx4LgFlyRUA | POST     /rpcHandler/Greet
```

And this is without any performance optimization in RPC's part — just plain JSON over HTTP.

Let's run other example endpoints as well:

```
$ curl -d'{"name": "Max"}' http://localhost:8282/greeter/greet
{"error":"","message":"Hello, Max"}
```

This has been an internal RPC endpoint, but it is exposed as part of API surface: you're talking directly to a `greeter` node directly, but using our `mesh_api` gateway as an entrypoint. Please note, that if RPC protocol's data serialization method is other than JSON, for example Protobuf or Cap'n'proto, that means you will deal with binary data.

```
$ curl http://localhost:8282/greeter/check
All ok! 2019-05-24T16:29:04+03:00
```

This has been a "legacy" HTTP endpoint that greeter service exposed before. But now it is accessed
using unified API surface of `mesh_api`, with RPC telemetry, logging and (possible) security and
other stuff attached to it. See [greeter/service/handlers.go](https://github.com/astranet/meshRPC/tree/master/example/greeter/service/handlers.go) example on what has been added before it could be exposed (spoiler: almost nothing at all).

A final example of complex data transfer between two services:

```
$ curl -X POST http://localhost:8282/greeter/sendPostcard/Max/World/Hello
{"PictureURL":"","Address":"World","Recipient":"Max","Message":"Hello"}
```

In `greeter` logs:

```
sending greeter.Postcard{PictureURL:"", Address:"World", Recipient:"Max", Message:"Hello"} to Max
[GIN] 2019/05/24 - 16:35:16 | 200 |      295.19µs |    Y7F51Lq8xN0A | POST     /rpcHandler/SendPostcard
```

Great! Next we need to do is orchestration and scaling.

### Example: Orchestration, Scaling

```
kek
``` 

At this point out tutorial and example section is over. We kindly forwarding you to [meshRPC/example](https://github.com/astranet/meshRPC/tree/master/example) for reference implementation and a playground for starting your cluster.



### Benchmarks


### Fixing templates

1) Edit `templates/XXX_go.tpl`
2) Run `go generate`
3) `go install`

### License

MIT
