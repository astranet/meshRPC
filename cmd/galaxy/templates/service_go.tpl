package {{.PackageName}}

import (
	cache "github.com/patrickmn/go-cache"

	bugsnag "github.com/bugsnag/bugsnag-go"
	log "github.com/sirupsen/logrus"

	"github.com/astranet/galaxy/metrics"
)

//go:generate galaxy expose -y -P {{.PackageName}}

type {{.FeaturePrefix}}Service interface {
{{if .ExampleFuncName}}
	{{.ExampleFuncName}}() error
{{end}}
}

type {{.FeaturePrefix}}ServiceOptions struct {
}

func check{{.FeaturePrefix}}ServiceOptions(opt *{{.FeaturePrefix}}ServiceOptions) *{{.FeaturePrefix}}ServiceOptions {
	if opt == nil {
		opt = &{{.FeaturePrefix}}ServiceOptions{}
	}
	return opt
}

func New{{.FeaturePrefix}}Service(
	{{if .RepoEnabled}}repo {{.FeaturePrefix}}DataRepo,{{end}}
	opt *{{.FeaturePrefix}}ServiceOptions,
) {{.FeaturePrefix}}Service {
	return &{{.ServicePrivateName}}{
		opt: check{{.FeaturePrefix}}ServiceOptions(opt),
		tags: metrics.Tags{
			"layer":   "service",
			"service": "{{.ServiceName}}",
		},
		fields: log.Fields{
			"layer":   "service",
			"service": "{{.ServiceName}}",
		}

		{{if .RepoEnabled}}repo: repo,{{end}}
	}
}

type {{.ServicePrivateName}} struct {
	{{if .RepoEnabled}}repo {{.FeaturePrefix}}DataRepo{{end}}
	tags metrics.Tags
	fields log.Fields
	opt  *{{.FeaturePrefix}}ServiceOptions
}

{{if .ExampleFuncName}}
func (s *{{.ServicePrivateName}}) {{.ExampleFuncName}}() error {
	metrics.ReportFuncCall(s.tags)
	statsFn := metrics.ReportFuncTiming(s.tags)
	defer statsFn()

	return nil
}
{{end}}