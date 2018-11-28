package metrics

import (
	"testing"
	"time"
)

func init() {
	client = newMockStatter(false)
	config = checkConfig(&MetricsConfig{
		EnvName:              "testing",
		StuckFunctionTimeout: time.Second,
	})
}

var testTags = Tags{
	"layer":   "service",
	"service": "users",
}

func TestReportFuncCall(t *testing.T) {
	ReportFuncCall(testTags)
}

func TestReportFuncTiming(t *testing.T) {
	stopFn := ReportFuncTiming(testTags)
	time.Sleep(500 * time.Millisecond)
	stopFn()
}

func TestReportFuncTimingStuck(t *testing.T) {
	stopFn := ReportFuncTiming(testTags)
	time.Sleep(2 * time.Second)
	stopFn()
}
