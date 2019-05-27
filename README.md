# Mesh RPC [![Go Report Card](https://goreportcard.com/badge/github.com/astranet/meshRPC)](https://goreportcard.com/report/github.com/astranet/meshRPC) ![Version Beta](https://img.shields.io/badge/version-beta-blue.svg)

<img alt="meshRPC logo, image author - MariaLetta/free-gophers-pack" src="https://cl.ly/0f19c50a98df/51.png" width="250px" />

MeshRPC is an automatic Service Mesh generator for pure Go micro services, it's a humble alternative to gRPC! In a nutshell, a Service Mesh is an inter-service communication infrastructure, including RPC over HTTP.

_With a service mesh,_

* A given Microservice won’t directly communicate with the other microservices.
Rather all service-to-service communications will take places on-top of a software component called service mesh.

* Service Mesh provides built-in support for some network functions such as resiliency, tracing, service discovery, etc.

* Therefore, service developers can focus more on the business logic while most of the work related to network communication is offloaded to the service mesh.

See [this article](https://medium.com/microservices-in-practice/service-mesh-for-microservices-2953109a3c9a) for a good explanation about this emerging concept, in terms of gRPC and Kubernetes.

The `meshRPC` tool paired with `cluster` package is able to transform any type of legacy Go service into a "new stack" developer's dream, without adding too much cost to infrastructure.

All meshRPC-infused microservices require zero-configuration and are load-balanced (round robin with sticky sessions) out of the box!

#### What about legacy monoliths? Should I rewrite 100%?

Of course not! Even for legacy Go monoliths that contain many layers inside a single process, `meshRPC` framework can be used to decouple things, using interface substitution. Consider an `interface A`, then use this tool to generate a microservice that implements `interface A`, but once called, instead of `interface A` invocation, there will be an RPC call over network to the corresponding microservice. In this way you can separate a big project by small pieces without hurting integrity (just adding a bit of network latency).

And you don't need to write RPC models and new data structures, we will generate them for you.

### Install

```
$ go get -u github.com/astranet/meshRPC
```

### Usage in Go

Create a service file like this one:

```go
package greeter

//go:generate meshRPC expose -P greeter -y

type Service interface {
    Greet(name string) (string, error)
    SendPostcard(card *Postcard) error
}

func NewService() Service {
    return &service{}
}

type service struct{}

func (s *service) Greet(name string) (string, error) {
    // TODO
}

func (s *service) SendPostcard(card *Postcard) error {
    // TODO
}

```

Make sure to implement the methods, they will be exposed to other microservices in cluster soon. See the full [example/greeter/service/service.go](https://github.com/astranet/meshRPC/tree/master/example/greeter/service/service.go) for the reference.

Now run `meshRPC expose` or using `go generate`, please note that when running manually, you must specify project dir and the target sources path as arguments. Also, if you have multiple service interfaces in the same package, called for example `FooService` and `BarService`, then `Foo` and `Bar` are module prefixes and should be provided using an additional flag `-M` on each expose call.

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

#### Connect to cluster

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

See the full [example/greeter/main.go](https://github.com/astranet/meshRPC/tree/master/example/greeter/main.go) that populates config values and starts a server. It took just a few lines to expose the service over network, as well as your custom HTTP endpoints. But, the question is, how to access those endpoints? For debugging purposes it's easy to also use `c.ListenAndServeHTTP` that will start an usual HTTP server in the same process, so you could connect to `greeter` instance directly. But that's not robust, because you want to have an API Gateway with load balancing and authorization, am I right? :)

At this point users that would use gRPC usually put an [Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/grpc) instance that will automatically convert HTTP/2 gRPC protobufs into HTTP/1.1 JSONs. You're expected to write an Envoy config and an xDS discovery service or install [Istio](https://istio.io/docs/concepts/security/) with automatic sidecar injection, on your Kubernetes cluster. It might be time consuming.

#### API Gateway for meshRPC

We'll create a simple server that will act as an HTTP server, but it will also be a cluster discovery endpoint. All service nodes you will start will simply connect to it in order to discover each other. Actually, any service node may connect to any other service node to get into the mesh. MeshRPC initiates persistent TCP connections between peers, and all virtual connections and streams are multiplexed on these persistent conns. This is how [AstraNet](https://github.com/astranet/astranet) cluster works.

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
    // Listen on a TCP address, this address can be used
    // by other peers to discover each other in this cluster.
    if err := c.ListenAndServe(*netAddr); err != nil {
        closer.Fatalln(err)
    }

    // Start a normal Gin HTTP server that will use cluster endpoints.
    httpListenAndServe(c)
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
    // A greeter.ServiceClient instance sucessfully conforms the greeter.Service interface
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

// See the rest in:
// https://github.com/astranet/meshRPC/tree/master/example/mesh_api/main.go
```

See the full [example/mesh_api/main.go](https://github.com/astranet/meshRPC/tree/master/example/mesh_api/main.go) to get the idea of a minimal API gateway. It took just a few lines more to create an API gateway for the whole cluster. Make sure that `:11999` port is closed from ouside connections, but `:8282` is okay to be exposed, or connected to classic reverse proxies such as Nginx and Caddy for TLS.

### Example Run

At this moment we have:
* `github.com/astranet/meshRPC/example/greeter` that is a server executable exposing service to meshRPC cluster.
* `github.com/astranet/meshRPC/example/mesh_api` that is an API Gateway for all meshRPC services in this cluster.

Let's run our API Gateway we just created:
```
$ go install github.com/astranet/meshRPC/example/mesh_api
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

This is an internal RPC endpoint, but it is exposed as part of API surface: you're talking directly to a `greeter` node, but using your `mesh_api` gateway as an entrypoint. Please note, that if data serialization protocol is other than JSON, Protobuf for example, it means you will deal with binary data.

```
$ curl http://localhost:8282/greeter/check
All ok! 2019-05-24T16:29:04+03:00
```

This is the "legacy" HTTP endpoint that `greeter` service always had. But it is exposed
using unified API surface of `mesh_api` now, with RPC telemetry, logging and (possibly some) security and
other stuff attached to it. See an [example/greeter/service/handlers.go](https://github.com/astranet/meshRPC/tree/master/example/greeter/service/handlers.go) for a demo on how to expose legacy HTTP endpoints, almost no changes are required.

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

Great! Next we do orchestration and scaling.

### Dockerization and scaling

Let's create a simple `Dockerfile` that creates minimalistic Alpine containers. This is up to you which method to use in practice, for example in some project I'd use Go's base image because I need a lot of dependencies that don't exist in Alpine.

```
FROM golang:1.12-alpine as builder

RUN apk add --no-cache git

ENV GOPATH=/gopath
RUN go get github.com/astranet/meshRPC/example/greeter

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /gopath/bin/greeter /usr/local/bin/

EXPOSE 11999
ENTRYPOINT ["greeter"]
``` 

A similar [Dockerfile](https://github.com/astranet/meshRPC/tree/master/example/mesh_api/Dockerfile) has been made for `mesh_api` too. Note: for real apps you should use vendoring or Go modules instead of simple `go get` in Dockefile.

To get both containers that we've built, use

```
$ docker pull docker.direct/meshrpc/example/greeter
$ docker pull docker.direct/meshrpc/example/mesh_api
```

#### Docker Stack (ex docker-compose)

First, we need `docker-compose.yml` that contains a simple definition. Here we use docker images generated from above, and we set cluster node list using an environment variable `MESHRPC_CLUSTER_NODES` that will be read by our example applications. We'll expose only `8282` port for the API Gateway access. Virtual net allows to reference nodes by their service names only, however, if you want to use host net just make sure to avoid port collision for `mesh_api` and other nodes, and use full addresses.

```yaml
version: "3"
services:
  mesh-api:
    image: docker.direct/meshrpc/example/mesh_api:latest
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
    environment:
      - MESHRPC_CLUSTER_NODES=mesh-api,greeter
    ports:
      - "8282:8282"
    depends_on:
      - greeter
    networks:
      - meshnet
  greeter:
    image: docker.direct/meshrpc/example/greeter:latest
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
    environment:
      - MESHRPC_CLUSTER_NODES=mesh-api,greeter
    networks:
      - meshnet
networks:
  meshnet:
```

Let's run this stack, you can use good old `docker-compose` but for the future's sake I'll use `docker stack` here.

```
$ docker stack deploy -c docker-compose.yml meshrpc-example
Creating network meshrpc-example_meshnet
Creating service meshrpc-example_mesh-api
Creating service meshrpc-example_greeter

$ docker stack ls
NAME                SERVICES            ORCHESTRATOR
meshrpc-example     2                   Swarm
```

Check if everything has started correctly:

```
CONTAINER ID   IMAGE                                           NAMES
e0cde9c1ce9d   docker.direct/meshrpc/example/greeter:latest    meshrpc-example_greeter.1.olhmh4ho2kvqck37zm7qm3zbt
767aee3bd128   docker.direct/meshrpc/example/mesh_api:latest   meshrpc-example_mesh-api.1.dzee3mvff6a3duzejeb1gg0q1
```

Run the usual API query:

```
$ curl http://localhost:8282/greeter/check
All ok! 2019-05-24T19:09:42Z
```

Works flawlessly! Check the docker logs (the `greeter` container): 

```
$ docker logs -f e0cde9c1ce9d
[GIN] 2019/05/24 - 19:06:58 | 200 |        15.6µs |    EvbuE1rKRY0A | GET      /ping
[GIN] 2019/05/24 - 19:09:41 | 200 |       106.9µs |      10.255.0.2 | GET      /handler/Check
[GIN] 2019/05/24 - 19:09:42 | 200 |        62.5µs |      10.255.0.2 | GET      /handler/Check
[GIN] 2019/05/24 - 19:09:42 | 200 |       350.6µs |      10.255.0.2 | GET      /handler/Check
```

Maybe scale the service a little?

```
$ docker service scale meshrpc-example_greeter=5
meshrpc-example_greeter scaled to 5
overall progress: 5 out of 5 tasks
1/5: running   [==================================================>]
2/5: running   [==================================================>]
3/5: running   [==================================================>]
4/5: running   [==================================================>]
5/5: running   [==================================================>]
verify: Service converged
```

And check the load balancing after querying the API Gateway. Nodes become available almost immediately, and require no warmup. Nodes can be taken down without hanging connections.

```
meshrpc-example_greeter.1.olhmh4ho2kvq    | [GIN] | 200 |       287.5µs |      10.255.0.2 | GET      /handler/Check
meshrpc-example_greeter.1.olhmh4ho2kvq    | [GIN] | 200 |       183.9µs |      10.255.0.2 | GET      /handler/Check
meshrpc-example_greeter.1.olhmh4ho2kvq    | [GIN] | 200 |        74.2µs |      10.255.0.2 | GET      /handler/Check
meshrpc-example_greeter.2.y6cen26m1wcu    | [GIN] | 200 |       146.2µs |      10.255.0.2 | GET      /handler/Check
meshrpc-example_greeter.2.y6cen26m1wcu    | [GIN] | 200 |       212.4µs |      10.255.0.2 | GET      /handler/Check
meshrpc-example_greeter.5.t900isxh3568    | [GIN] | 200 |       175.5µs |      10.255.0.2 | GET      /handler/Check
meshrpc-example_greeter.5.t900isxh3568    | [GIN] | 200 |       304.3µs |      10.255.0.2 | GET      /handler/Check
```

At this point our tutorial and example section is over. We kindly forwarding you to [example](https://github.com/astranet/meshRPC/tree/master/example) dir for reference implementation and a playground for starting your cluster.

### Benchmarks

[MeshRPC Benchmark Suite](https://github.com/astranet/meshRPC-benchmark)

**1.8 ms** per call is the current latency using `docker stack` on local machine and virtual network.
Tested on 2014 Macbook Pro (2,8 GHz Intel Core i5).

### Fixing templates

1) Edit `templates/XXX_go.tpl`
2) Run `go generate`
3) `go install`

### License

MIT
