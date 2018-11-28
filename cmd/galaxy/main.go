package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jawher/mow.cli"
)

var app = cli.App("galaxy", "A tool for generating boilerplate code for new services and components using astranet:Galaxy toolkit.")
var projectDir = app.StringOpt("D dir", defaultProjectDir(), "Sets the target project root.")

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	app.Command("new", "Creates a new service in the core.", newCmd)
	app.Command("add", "Creates a component in some existing service.", addCmd)
	app.Command("expose", "Generates RPC handler and cluster client for a service.", exposeCmd)
	app.Command("grafana", "Generates Grafana rows for services from dashboard source.", grafanaCmd)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func newCmd(c *cli.Cmd) {
	packageName := c.StringOpt("P pkg", "foo", "Specifies an existing package name in core.")
	featurePrefix := c.StringOpt("M module", "", "Feature prefix to distinguish components in the package.")
	exampleFuncName := c.StringOpt("F func", "", "Generate an example function attached to repo, service and handler.")
	serviceName := c.StringOpt("S service", "foo", "Specify service name for metrics and reporting.")
	includeFiles := c.StringOpt("I include", "r,s,h", "Generate particular files (r = repo, s = service, h = handler).")
	excludeFiles := c.StringOpt("E exclude", "", "Skip particular files (r = repo, s = service, h = handler).")
	agreeAll := c.BoolOpt("y yes", false, "Agree to all prompts automatically.")

	c.Action = func() {
		fileSet := decodeIncludesExcludes(*includeFiles, *excludeFiles)
		ctx := &TemplateContext{
			PackageName:     strings.ToLower(*packageName),
			FeaturePrefix:   strings.Title(*featurePrefix),
			ExampleFuncName: strings.Title(*exampleFuncName),
			ServiceName:     strings.ToLower(*serviceName),

			RepoPrivateName:    repoPrivateName(*featurePrefix),
			ServicePrivateName: servicePrivateName(*featurePrefix),
			HandlerPrivateName: handlerPrivateName(*featurePrefix),

			RepoEnabled:    fileSet["repo"],
			ServiceEnabled: fileSet["service"],
			HandlerEnabled: fileSet["handler"],
		}
		basePath := filepath.Join(*projectDir, "core", *packageName)
		filePrefix := strings.ToLower(*featurePrefix) + "_"
		if len(*featurePrefix) == 0 {
			filePrefix = ""
		}
		actionQueue := NewQueue(
			NewDirAction(basePath),
		)
		if ctx.RepoEnabled {
			actionQueue = append(actionQueue,
				CreateNewFileAction(filepath.Join(basePath, filePrefix+"data.go"), ctx.RenderInto(dataTemplate)),
			)
		}
		if ctx.ServiceEnabled {
			actionQueue = append(actionQueue,
				CreateNewFileAction(filepath.Join(basePath, filePrefix+"service.go"), ctx.RenderInto(serviceTemplate)),
			)
		}
		if ctx.ServiceEnabled && ctx.HandlerEnabled {
			actionQueue = append(actionQueue,
				CreateNewFileAction(filepath.Join(basePath, filePrefix+"handler.go"), ctx.RenderInto(handlerTemplate)),
			)
		}
		fmt.Println(actionQueue.Description())
		agree := *agreeAll
		if !agree {
			agree = cliConfirm("Are you sure to apply these changes?")
			if !agree {
				log.Println("Action cancelled.")
				return
			}
		}
		ts := time.Now()
		if !actionQueue.Exec() {
			log.Println("Failed in", time.Since(ts))
			return
		}
		log.Println("Done in", time.Since(ts))
	}
}

func addCmd(c *cli.Cmd) {
	packageName := c.StringOpt("P pkg", "foo", "Specifies an existing package name in core.")
	featurePrefix := c.StringOpt("M module", "Bar", "Mandatory feature prefix to distinguish components in the package.")
	exampleFuncName := c.StringOpt("F func", "", "Generate an example function attached to repo, service and handler.")
	serviceName := c.StringOpt("S service", "bar", "Specify service name for metrics and reporting of this module.")
	includeFiles := c.StringOpt("I include", "r,s,h", "Generate particular files (r = repo, s = service, h = handler).")
	excludeFiles := c.StringOpt("E exclude", "", "Skip particular files (r = repo, s = service, h = handler).")
	agreeAll := c.BoolOpt("y yes", false, "Agree to all prompts automatically.")

	c.Action = func() {
		fileSet := decodeIncludesExcludes(*includeFiles, *excludeFiles)
		ctx := &TemplateContext{
			PackageName:     strings.ToLower(*packageName),
			FeaturePrefix:   strings.Title(*featurePrefix),
			ExampleFuncName: strings.Title(*exampleFuncName),
			ServiceName:     strings.ToLower(*serviceName),

			RepoPrivateName:    repoPrivateName(*featurePrefix),
			ServicePrivateName: servicePrivateName(*featurePrefix),
			HandlerPrivateName: handlerPrivateName(*featurePrefix),

			RepoEnabled:    fileSet["repo"],
			ServiceEnabled: fileSet["service"],
			HandlerEnabled: fileSet["handler"],
		}
		basePath := filepath.Join(*projectDir, "core", *packageName)
		filePrefix := strings.ToLower(*featurePrefix) + "_"
		if len(*featurePrefix) == 0 {
			log.Println("Warning: no feature prefix, adding module files without prefix.")
			filePrefix = ""
		}
		actionQueue := NewQueue(
			CheckDirAction(basePath),
		)
		if ctx.RepoEnabled {
			actionQueue = append(actionQueue,
				CreateNewFileAction(filepath.Join(basePath, filePrefix+"data.go"), ctx.RenderInto(dataTemplate)),
			)
		}
		if ctx.ServiceEnabled {
			actionQueue = append(actionQueue,
				CreateNewFileAction(filepath.Join(basePath, filePrefix+"service.go"), ctx.RenderInto(serviceTemplate)),
			)
		}
		if ctx.ServiceEnabled && ctx.HandlerEnabled {
			actionQueue = append(actionQueue,
				CreateNewFileAction(filepath.Join(basePath, filePrefix+"handler.go"), ctx.RenderInto(handlerTemplate)),
			)
		}
		fmt.Println(actionQueue.Description())
		agree := *agreeAll
		if !agree {
			agree = cliConfirm("Are you sure to apply these changes?")
			if !agree {
				log.Println("Action cancelled.")
				return
			}
		}
		ts := time.Now()
		if !actionQueue.Exec() {
			log.Println("Failed in", time.Since(ts))
			return
		}
		log.Println("Done in", time.Since(ts))
	}
}

func exposeCmd(c *cli.Cmd) {
	packageName := c.StringOpt("P pkg", "foo", "Specifies an existing package name in core.")
	featurePrefix := c.StringOpt("M module", "", "Optional feature prefix to distinguish interfaces of the package.")
	serviceName := c.StringOpt("S service", "", "Specify service name for metrics and reporting.")
	includeFiles := c.StringOpt("I include", "h,c", "Generate particular files (h = handler, c = service client).")
	excludeFiles := c.StringOpt("E exclude", "", "Skip particular files (h = handler, c = service client).")
	agreeAll := c.BoolOpt("y yes", false, "Agree to all prompts automatically.")

	c.Action = func() {
		fileSet := decodeIncludesExcludes(*includeFiles, *excludeFiles)
		ctx := &TemplateContext{
			PackageName:   strings.ToLower(*packageName),
			FeaturePrefix: strings.Title(*featurePrefix),
			ServiceName:   strings.ToLower(*serviceName),

			RPCHandlerPrivateName: rpcHandlerPrivateName(*featurePrefix),
			ClientPrivateName:     clientPrivateName(*featurePrefix),

			HandlerEnabled: fileSet["handler"],
			ClientEnabled:  fileSet["client"],
		}
		basePath := filepath.Join(*projectDir, "core", *packageName)
		filePrefix := strings.ToLower(*featurePrefix) + "_"
		if len(*featurePrefix) == 0 {
			filePrefix = ""
		}
		if info, err := os.Stat(basePath); err != nil || !info.IsDir() {
			log.Fatalln("Failed to read directory:", basePath, err)
			return
		}
		ifaceName := fmt.Sprintf("%s.%sService", ctx.PackageName, ctx.FeaturePrefix)
		iface, err := NewMethodsCollection(ifaceName, basePath)
		if err != nil {
			log.Fatalf("Failed to locate %s interface: %v", ifaceName, err)
			return
		}
		if len(ctx.ServiceName) == 0 {
			set := make(map[string]map[string]bool)
			if err := scanMetricTagsFile(set, iface.SrcPath); err != nil {
				log.Fatalf("Failed to scan metrics.Tags from target file: %v", err)
				return
			}
			ctx.ServiceName = anyServiceIn(set)
		}
		actionQueue := NewQueue(
			CheckDirAction(basePath),
		)
		if ctx.HandlerEnabled {
			ctx.RPCHandlerInterfaceBody = genRPCHandlerInterface(iface)
			ctx.RPCHandlerImplementationBody = genRPCHandlerImplementation(ctx.RPCHandlerPrivateName, ctx.FeaturePrefix, iface)
			actionQueue = append(actionQueue,
				OverwriteFileAction(filepath.Join(basePath, filePrefix+"handler_rpc.go"), ctx.RenderInto(rpcHandlerTemplate)),
			)
		}
		if ctx.HandlerEnabled && ctx.ClientEnabled {
			ctx.ClientImplementationBody = genServiceClientImplementation(ctx.ClientPrivateName, ctx.FeaturePrefix, iface)
			actionQueue = append(actionQueue,
				OverwriteFileAction(filepath.Join(basePath, filePrefix+"client_rpc.go"), ctx.RenderInto(rpcClientTemplate)),
			)
		}
		fmt.Println(actionQueue.Description())
		agree := *agreeAll
		if !agree {
			agree = cliConfirm("Are you sure to apply these changes?")
			if !agree {
				log.Println("Action cancelled.")
				return
			}
		}
		ts := time.Now()
		if !actionQueue.Exec() {
			log.Println("Failed in", time.Since(ts))
			return
		}
		log.Println("Done in", time.Since(ts))
	}
}

func defaultProjectDir() string {
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		panic("no $GOPATH env var set")
	}
	return filepath.Join(gopath, "src", "github.com", "astranet", "example_api")
}

type TemplateContext struct {
	PackageName     string
	FeaturePrefix   string
	ExampleFuncName string
	ServiceName     string

	RepoPrivateName       string
	ServicePrivateName    string
	HandlerPrivateName    string
	RPCHandlerPrivateName string
	ClientPrivateName     string

	RepoEnabled    bool
	ServiceEnabled bool
	HandlerEnabled bool
	ClientEnabled  bool

	RPCHandlerInterfaceBody      string
	RPCHandlerImplementationBody string
	ClientInterfaceBody          string
	ClientImplementationBody     string
}

//go:generate go-bindata -o bindata.go -pkg main templates/
var (
	dataTemplate = template.Must(
		template.New("data.go").Parse(
			string(MustAsset("templates/data_go.tpl")),
		),
	)
	serviceTemplate = template.Must(
		template.New("service.go").Parse(
			string(MustAsset("templates/service_go.tpl")),
		),
	)
	handlerTemplate = template.Must(
		template.New("handler.go").Parse(
			string(MustAsset("templates/handler_go.tpl")),
		),
	)
	rpcHandlerTemplate = template.Must(
		template.New("handler_rpc.go").Parse(
			string(MustAsset("templates/handler_rpc_go.tpl")),
		),
	)
	rpcClientTemplate = template.Must(
		template.New("client_rpc.go").Parse(
			string(MustAsset("templates/client_rpc_go.tpl")),
		),
	)
)

func (t *TemplateContext) RenderInto(tpl *template.Template) []byte {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, tpl.Name(), t); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func repoPrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "repo"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "Repo"
}

func handlerPrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "handler"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "Handler"
}

func servicePrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "service"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "Service"
}

func rpcHandlerPrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "rpcHandler"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "RPCHandler"
}

func clientPrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "client"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "ServiceClient"
}

func decodeIncludesExcludes(includes, excludes string) map[string]bool {
	includes = strings.ToLower(includes)
	excludes = strings.ToLower(excludes)
	includeList := strings.Split(includes, ",")
	excludeList := strings.Split(excludes, ",")
	set := map[string]bool{}

	for _, name := range includeList {
		switch strings.TrimSpace(name) {
		case "r", "repo", "repos", "data", "d":
			set["repo"] = true
		case "s", "service", "services":
			set["service"] = true
		case "h", "handler", "handlers":
			set["handler"] = true
		case "c", "client", "clients":
			set["client"] = true
		}
	}
	if len(set) == 0 {
		// nothing explicitly included - include all
		set = map[string]bool{
			"repo":    true,
			"service": true,
			"handler": true,
			"client":  true,
		}
	}
	for _, name := range excludeList {
		switch strings.TrimSpace(name) {
		case "r", "repo", "repos", "data", "d":
			set["repo"] = false
		case "s", "service", "services":
			set["service"] = false
		case "h", "handler", "handlers":
			set["handler"] = false
		case "c", "client", "clients":
			set["client"] = true
		}
	}
	return set
}
