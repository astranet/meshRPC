## Vroomy Plugin

This is a plugin that enables MeshRPC capabilities inside [Vroomy Server](https://github.com/vroomy/vroomy). It allows to expose handlers and RPC handlers through vroomy server. Follow this manual for usage instructions.

### Install

To install vroomy using Go:

```
go get github.com/vroomy/vpm
go get github.com/vroomy/vroomy
```

And get the greeter example service that has been generated using meshRPC:

```
go get github.com/astranet/meshRPC/example/greeter
```

Alternatively, run a prepared Greeter image for Docker:

```
docker run -it --rm docker.direct/meshrpc/example/greeter -h

Usage: greeter [OPTIONS]

A Greeter service server for meshRPC cluster.
```

### Usage

#### Step 1: define routes

This plugin has two methods: `Ping` and `Route(...)`. Ping can be used to check services for liveness, it returns the latency as a result. To enable Ping on a route:

```toml
[[route]]
method = "GET"
group = "services"
httpPath = "/ping/:service"
handlers = [
    "meshrpc.Ping",
]
```

The method Route allows to map any endpoint onto meshRPC handler, either exposed `handler` or generated `rpcHandler`. The method accepts 3 arguments:

    * Service name
    * Handler name (must match the interface name)
    * Method name

Example:

```toml
[[route]]
method = "GET"
group = "greeter"
httpPath = "/check"
handlers = [
    "meshrpc.Route(greeter,handler,Check)",
]

[[route]]
method = "POST"
group = "greeter"
httpPath = "/greet"
handlers = [
    "meshrpc.Route(greeter,rpcHandler,Greet)",
]
```

#### Step 2: check global config

Make sure to import the meshrpc plugin in your configuration, also specify additional parameters for the cluster config using `env` section.

```toml
plugins = [
    "github.com/astranet/meshRPC/plugins/vroomy as meshrpc",
]

[env]
# meshrpc_host = "localhost"
meshrpc_port = "11999"
meshrpc_debug = "true"
# meshrpc_cluster_name = "vroomy-meshrpc"
# meshrpc_cluster_nodes = []
```

By default this plugin will listen on `:11999` for incoming connections from meshRPC-compatible apps, and the namespace is set as `vroomy-meshrpc`.

#### Step 3: run Vroomy!

(it seems that `vpm update` is broken, use `cd ./plugins/vroomy && make plugin` meanwhile)

```
$ vroomy

Vroomy :: Hello there! One moment, initializing..
● Vroomy :: Starting service
● Plugins :: Initialized meshrpc (plugins/meshrpc.so)
INFO[0000] new astranet router created                   layer=cluster net_env=vroomy-meshrpc service=vroomy
INFO[0000] binding vroomy as meshrpc.vroomy              fn=ListenAndServe
INFO[0000] ListenAndServe on :11999                      fn=ListenAndServe
● Vroomy :: HTTP is now listening on port 8080 (HTTP)
```

Amazing! It runs both HTTP and meshRPC servers now.

#### Step 4: run Greeter service

```
$ docker run -it --rm docker.direct/meshrpc/example/greeter --nodes docker.for.mac.host.internal:11999 --tag vroomy-meshrpc

INFO[0000] new astranet router created                   layer=cluster net_env=vroomy-meshrpc service=greeter
INFO[0000] binding greeter as meshrpc.greeter            fn=ListenAndServe
INFO[0000] ListenAndServe on 0.0.0.0:0                   fn=ListenAndServe
```

Note that we don't care about exposing ports, however we do care about nodes list (must include listening address of vroomy) and the cluster tag name that should match plugin's configuration.

#### Step 5: the ultimate showdown

Ping stuff:

```
$ curl http://localhost:8080/services/ping/greeter

{"data":"869.899µs","errors":[]}
```

Get stuff:

```
$ curl http://localhost:8080/services/greeter/check

{"data":{"fingerprint":"5b73ad07bc72a4ce","status":"ok","timestamp":"2019-08-31T05:36:26+03:00"},"errors":[]}
```

Automatic round-robin when multiple Greeter services are running:

```
$ curl http://localhost:8080/services/greeter/check
{"data":{"fingerprint":"f7404fa8b2be412c","status":"ok","timestamp":"2019-08-31T02:39:21Z"},"errors":[]}

$ curl http://localhost:8080/services/greeter/check
{"data":{"fingerprint":"c0a21eb358cc3d20","status":"ok","timestamp":"2019-08-31T02:39:23Z"},"errors":[]}

$ curl http://localhost:8080/services/greeter/check
{"data":{"fingerprint":"7047b7b73c53cff3","status":"ok","timestamp":"2019-08-31T02:39:24Z"},"errors":[]}
```

Post stuff to RPC methods that have been generated and exposed from the service interface:

```
$ curl -d'{"name": "Max"}' http://localhost:8080/services/greeter/greet

{"data":{"message":"Hello, Max"},"errors":[]}
```
