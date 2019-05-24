package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jawher/mow.cli"
)

var app = cli.App("meshRPC", "Tool for generating an RPC handler and cluster client for any service interface.")
var projectDir = app.StringOpt("R project-root", filepath.Join(".."), "Sets the target project root. By default the upper level of working dir")

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	app.Command("expose", "Creates RPC handler/client that exposes provided service into a mesh cluster.", exposeCmd)
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func exposeCmd(c *cli.Cmd) {
	targetPath := c.StringArg("SRC", ".", "Target Go source file or a package with service definitions.")
	packageName := c.StringOpt("P pkg-name", "foo", "Must specify the package name.")
	featurePrefix := c.StringOpt("M module-prefix", "", "Optional feature prefix to distinguish multiple service interfaces in the same package.")
	agreeAll := c.BoolOpt("y yes", false, "Agree to all prompts automatically.")
	c.Spec = "-P [-M] [-y] [SRC]"

	c.Action = func() {
		ctx := &TemplateContext{
			PackageName:   strings.ToLower(*packageName),
			FeaturePrefix: strings.Title(*featurePrefix),

			RPCHandlerPrivateName: rpcHandlerPrivateName(*featurePrefix),
			RPCClientPrivateName:  rpcClientPrivateName(*featurePrefix),
		}
		var basePath string
		if info, err := os.Stat(*targetPath); err != nil {
			log.Fatalln("Failed to read SRC dir:", *targetPath)
		} else if !info.IsDir() {
			basePath = filepath.Dir(*targetPath)
		} else {
			basePath = *targetPath
		}
		basePath, _ = filepath.Abs(basePath)
		if len(*projectDir) == 0 {
			*projectDir = "."
		}
		*projectDir, _ = filepath.Abs(*projectDir)

		filePrefix := strings.ToLower(*featurePrefix) + "_"
		if len(*featurePrefix) == 0 {
			filePrefix = ""
		}
		ifaceName := fmt.Sprintf("%s.%sService", ctx.PackageName, ctx.FeaturePrefix)
		iface, err := NewMethodsCollection(ifaceName, basePath)
		if err != nil {
			log.Fatalf("Failed to locate %s interface: %v", ifaceName, err)
			return
		}
		actionQueue := NewQueue(
			CheckDirAction(basePath),
		)
		ctx.JsonHandlerInterfaceBody = genRPCHandlerInterface(iface)
		ctx.JsonHandlerImplementationBody = genRPCHandlerImplementation(ctx.RPCHandlerPrivateName, ctx.FeaturePrefix, iface)
		actionQueue = append(actionQueue,
			OverwriteFileAction(filepath.Join(basePath, filePrefix+"handler_gen.go"), ctx.RenderInto(rpcHandlerTemplate)),
		)
		ctx.JsonClientImplementationBody = genServiceClientImplementation(ctx.RPCClientPrivateName, ctx.FeaturePrefix, iface)
		actionQueue = append(actionQueue,
			OverwriteFileAction(filepath.Join(basePath, filePrefix+"client_gen.go"), ctx.RenderInto(rpcClientTemplate)),
		)
		fmt.Println(actionQueue.Description())
		agree := *agreeAll
		if !agree {
			agree = cliConfirm("Are you sure to apply these changes?")
			if !agree {
				log.Println("Action cancelled.")
				return
			}
		}
		if !actionQueue.Exec() {
			os.Exit(1)
			return
		}
	}
}

type TemplateContext struct {
	PackageName   string
	FeaturePrefix string

	RPCHandlerPrivateName string
	RPCClientPrivateName  string

	JsonHandlerInterfaceBody      string
	JsonHandlerImplementationBody string
	CapnHandlerInterfaceBody      string
	CapnHandlerImplementationBody string

	JsonClientInterfaceBody      string
	JsonClientImplementationBody string
	CapnClientInterfaceBody      string
	CapnClientImplementationBody string
}

//go:generate go-bindata -o bindata.go -pkg main templates/
var (
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

func rpcHandlerPrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "rpcHandler"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "RPCHandler"
}

func rpcClientPrivateName(featurePrefix string) string {
	if len(featurePrefix) == 0 {
		return "rpcClient"
	}
	return strings.ToLower(string(featurePrefix[0])) + featurePrefix[1:] + "RPCClient"
}
