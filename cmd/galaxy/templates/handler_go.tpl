package {{.PackageName}}

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bugsnag "github.com/bugsnag/bugsnag-go"
	log "github.com/sirupsen/logrus"

	"github.com/astranet/galaxy/logging"
	"github.com/astranet/galaxy/metrics"
)

type {{.FeaturePrefix}}Handler interface {
{{ if .ExampleFuncName }}
	{{.ExampleFuncName}}(c *gin.Context)
{{ end }}
}

var {{.FeaturePrefix}}HandlerSpec {{.FeaturePrefix}}Handler = &{{.HandlerPrivateName}}{}

type {{.FeaturePrefix}}HandlerOptions struct {
}

func check{{.FeaturePrefix}}HandlerOptions(opt *{{.FeaturePrefix}}HandlerOptions) *{{.FeaturePrefix}}HandlerOptions {
	if opt == nil {
		opt = &{{.FeaturePrefix}}HandlerOptions{}
	}
	return opt
}

func New{{.FeaturePrefix}}Handler(
	svc {{.FeaturePrefix}}Service,
	opt *{{.FeaturePrefix}}HandlerOptions,
) {{.FeaturePrefix}}Handler {
	return &{{.HandlerPrivateName}}{
		opt: check{{.FeaturePrefix}}HandlerOptions(opt),
		tags: metrics.Tags{
			"layer":   "handler",
			"service": "{{.ServiceName}}",
		},
		fields: log.Fields{
			"layer":   "handler",
			"service": "{{.ServiceName}}",
		}

		svc:     svc,
	}
}

type {{.HandlerPrivateName}} struct {
	svc  {{.FeaturePrefix}}Service
	tags metrics.Tags
	fields log.Fields
	opt  *{{.FeaturePrefix}}HandlerOptions
}

{{ if .ExampleFuncName }}
func (h *{{.HandlerPrivateName}} ) {{.ExampleFuncName}}(c *gin.Context) {
	metrics.ReportFuncCall(h.tags)
	fnLog := log.WithFields(logging.WithFn(h.fields))
	statsFn := metrics.ReportFuncTiming(h.tags)
	defer statsFn()

	err := h.svc.{{.ExampleFuncName}}()
	if err != nil {
		fnLog.Warn(err)
		bugsnag.Notify(err)
		c.String(http.StatusInternalServerError, "error: %v", err)
		return
	}
}
{{ end }}