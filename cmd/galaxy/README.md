## galaxy

This utility allows to create new services and add modules into existing ones.

```
$ galaxy


Usage: galaxy [OPTIONS] COMMAND [arg...]

A tool for generating boilerplate code for new services and components using astranet:Galaxy toolkit.

Options:
  -D, --dir    Sets the target project root.

Commands:
  new          Creates a new service in the core.
  add          Creates a component in some existing service.
  expose       Generates RPC handler and cluster client for a service.
  grafana      Generates Grafana rows for services from dashboard source.

Run 'galaxy COMMAND --help' for more information on a command.
```

### Creating new

```
Usage: galaxy new [OPTIONS]

Creates a new service in the core.

Options:
  -P, --pkg       Specifies an existing package name in core. (default "foo")
  -M, --module    Feature prefix to distinguish components in the package.
  -F, --func      Generate an example function attached to repo, service and handler.
  -S, --service   Specify service name for metrics and reporting. (default "foo")
  -I, --include   Generate particular files (r = repo, s = service, h = handler). (default "r,s,h")
  -E, --exclude   Skip particular files (r = repo, s = service, h = handler).
```

For example, let's create a new service `something` with all three layers enabled and an example method attached named `FooFunc`. We won't use feature prefix to distinguish modules as we have only this main module.

```
$ galaxy new -P something -F FooThat -S something
.
├── [1]  new dir [project]/core/something if not exists
├── [2]  new file [project]/core/something/data.go with 60 lines of content (no overwrite)
├── [3]  new file [project]/core/something/service.go with 58 lines of content (no overwrite)
└── [4]  new file [project]/core/something/handler.go with 65 lines of content (no overwrite)

Are you sure to apply these changes?: y

queue.go:45: Action#1: new dir [project]/core/something if not exists
queue.go:45: Action#2: new file [project]/core/something/data.go with 60 lines of content (no overwrite)
queue.go:45: Action#3: new file [project]/core/something/service.go with 58 lines of content (no overwrite)
queue.go:45: Action#4: new file [project]/core/something/handler.go with 65 lines of content (no overwrite)
main.go:90: Done in 11.382682ms
```

You can include or exclude layers by listing them, e.g. `s,h` will enable only service and handler layers. Specifying `repo,service` will enable only repo and service layers. You can combine includes with excludes. 

### Adding modules

```
Usage: galaxy add [OPTIONS]

Creates a component in some existing service.

Options:
  -P, --pkg       Specifies an existing package name in core. (default "foo")
  -M, --module    Mandatory feature prefix to distinguish components in the package. (default "Bar")
  -F, --func      Generate an example function attached to repo, service and handler.
  -S, --service   Specify service name for metrics and reporting of this module. (default "bar")
  -I, --include   Generate particular files (r = repo, s = service, h = handler). (default "r,s,h")
  -E, --exclude   Skip particular files (r = repo, s = service, h = handler).
```

Let's add another module into this `something` service. At this time we would need to use feature prefix to avoid collision with main module of this service.

```
$ galaxy add -P something -F BarThis -S something_bar -M Bar
.
├── [1]  dir [project]/core/something must exist
├── [2]  new file [project]/core/something/bar_data.go with 60 lines of content (no overwrite)
├── [3]  new file [project]/core/something/bar_service.go with 58 lines of content (no overwrite)
└── [4]  new file [project]/core/something/bar_handler.go with 65 lines of content (no overwrite)

Are you sure to apply these changes?: y

queue.go:45: Action#1: dir [project]/core/something must exist
queue.go:45: Action#2: new file [project]/core/something/bar_data.go with 60 lines of content (no overwrite)
queue.go:45: Action#3: new file [project]/core/something/bar_service.go with 58 lines of content (no overwrite)
queue.go:45: Action#4: new file [project]/core/something/bar_handler.go with 65 lines of content (no overwrite)
main.go:153: Done in 16.024552ms
```

Now `something` package has two modules: `someting.Service` and `something.BarService`.

### Fixing templates

1) Edit `templates/XXX_go.tpl`
2) Run `go generate`
3) `go install`
