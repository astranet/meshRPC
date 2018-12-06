# Galaxy

A Microservices toolkit and service layers generator. Very opinionated.

### Install

```
$ go get -u github.com/astranet/galaxy
```

### Usage

```
$ galaxy -h

Usage: galaxy [OPTIONS] COMMAND [arg...]

A tool for generating boilerplate code for new services and components using astranet:Galaxy toolkit.

Options:
  -D, --dir    Sets the target project root. (default "$GOPATH/src/github.com/astranet/example_api")

Commands:
  new          Creates a new service in the core.
  add          Creates a component in some existing service.
  expose       Generates RPC handler and cluster client for a service.
  grafana      Generates Grafana rows for services from dashboard source.

Run 'galaxy COMMAND --help' for more information on a command.
```

### Example

```
$ galaxy new -h

Usage: galaxy new [OPTIONS]

Creates a new service in the core.

Options:
  -P, --pkg       Specifies an existing package name in core. (default "foo")
  -M, --module    Feature prefix to distinguish components in the package.
  -F, --func      Generate an example function attached to repo, service and handler.
  -S, --service   Specify service name for metrics and reporting. (default "foo")
  -I, --include   Generate particular files (r = repo, s = service, h = handler). (default "r,s,h")
  -E, --exclude   Skip particular files (r = repo, s = service, h = handler).
  -y, --yes       Agree to all prompts automatically.


$ galaxy --dir $GOPATH/src/foo_api new --pkg foo --func Foo
Actions to be committed
├── [1]  new dir [project]/core/foo if not exists
├── [2]  new file [project]/core/foo/data.go with 66 lines of content (no overwrite)
├── [3]  new file [project]/core/foo/service.go with 63 lines of content (no overwrite)
└── [4]  new file [project]/core/foo/handler.go with 72 lines of content (no overwrite)

Are you sure to apply these changes?: y
queue.go:44: Action#1: new dir [project]/core/foo if not exists
queue.go:44: Action#2: new file [project]/core/foo/data.go with 66 lines of content (no overwrite)
queue.go:44: Action#3: new file [project]/core/foo/service.go with 63 lines of content (no overwrite)
queue.go:44: Action#4: new file [project]/core/foo/handler.go with 72 lines of content (no overwrite)
main.go:96: Done in 10.344156ms
```

The expose generator, that can be triggered by running `go generate`, already placed in generated `service.go`.

```go
//go:generate galaxy -D ../.. expose -y -P greeter
```

See https://github.com/astranet/example_api for a full-featured example infrastructure.

### License

MIT
