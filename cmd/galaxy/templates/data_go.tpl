package {{.PackageName}}

import (
	cache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/astranet/galaxy/data"
	"github.com/astranet/galaxy/metrics"
)

type {{.FeaturePrefix}}DataRepo interface {
	data.Repo
{{ if .ExampleFuncName }}
	{{.ExampleFuncName}}() error
{{- end}}
}

type {{.FeaturePrefix}}DataRepoOptions struct {
}

func check{{.FeaturePrefix}}DataRepoOptions(opt *{{.FeaturePrefix}}DataRepoOptions) *{{.FeaturePrefix}}DataRepoOptions {
	if opt == nil {
		opt = &{{.FeaturePrefix}}DataRepoOptions{}
	}
	return opt
}

func New{{.FeaturePrefix}}DataRepo(
	globalRepo data.Repo,
	dataCache *cache.Cache,
	opt *{{.FeaturePrefix}}DataRepoOptions,
) {{.FeaturePrefix}}DataRepo {
	return &{{.RepoPrivateName}}{
		Repo: globalRepo,

		opt: check{{.FeaturePrefix}}DataRepoOptions(opt),
		tags: metrics.Tags{
			"layer":   "data",
			"service": "{{.ServiceName}}",
		},
		fields: log.Fields{
			"layer":   "data",
			"service": "{{.ServiceName}}",
		},

		dataCache: dataCache,
	}
}

type {{.RepoPrivateName}} struct {
	data.Repo

	tags metrics.Tags
	fields log.Fields
	opt  *{{.FeaturePrefix}}DataRepoOptions

	dataCache *cache.Cache
}

{{ if .ExampleFuncName }}
func (r *{{.RepoPrivateName}}) {{.ExampleFuncName}}() error {
	metrics.ReportFuncCall(r.tags)
	statsFn := metrics.ReportFuncTiming(r.tags)
	defer statsFn()

	return nil
}
{{ end }}