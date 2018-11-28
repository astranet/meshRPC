package metrics

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	bugsnag "github.com/bugsnag/bugsnag-go"

	"github.com/astranet/galaxy/metrics/memstatsd"
)

func RunMemstatsd(envName string, d time.Duration) {
	if client == nil {
		return
	}
	m := memstatsd.New("memstatsd.", envName, proxy{client})
	m.Run(d)
}

func ReportFuncError(tags ...Tags) {
	fn := funcName()
	reportFunc(fn, "error", tags...)
}

func ReportFuncStatus(tags ...Tags) {
	fn := funcName()
	reportFunc(fn, "status", tags...)
}

func ReportFuncCall(tags ...Tags) {
	fn := funcName()
	reportFunc(fn, "called", tags...)
}

func reportFunc(fn, action string, tags ...Tags) {
	if client == nil {
		return
	}
	tagSpec := config.BaseTags() + joinTags(tags...)
	tagSpec += ",func_name=" + fn
	client.Increment(fmt.Sprintf("func.%v", action) + tagSpec)
}

type StopTimerFunc func()

func ReportFuncTiming(tags ...Tags) StopTimerFunc {
	if client == nil {
		return func() {}
	}
	t := time.Now()
	fn := funcName()
	tagSpec := config.BaseTags() + joinTags(tags...)
	tagSpec += ",func_name=" + fn

	doneC := make(chan struct{})
	go func(name string, start time.Time) {
		select {
		case <-doneC:
			return
		case <-time.Tick(config.StuckFunctionTimeout):
			err := fmt.Errorf("detected stuck function: %s stuck for %v", name, time.Since(start))
			bugsnag.Notify(err)
			client.Increment("func.stuck" + tagSpec)
		}
	}(fn, t)

	return func() {
		d := time.Since(t)
		close(doneC)
		client.Timing("func.timing"+tagSpec, int(d/time.Millisecond))
	}
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")
	return nameParts[len(nameParts)-1]
}

type Tags map[string]string

func (t Tags) With(k, v string) Tags {
	if t == nil || len(t) == 0 {
		return map[string]string{
			k: v,
		}
	}
	t[k] = v
	return t
}

func joinTags(tags ...Tags) string {
	if len(tags) == 0 {
		return ""
	}
	var str string
	for k, v := range tags[0] {
		str += fmt.Sprintf(",%s=%s", k, v)
	}
	return str
}
