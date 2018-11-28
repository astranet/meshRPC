package metrics

import (
	"fmt"
	"time"

	"github.com/alexcesaro/statsd"
	"github.com/astranet/galaxy/logging"
	bugsnag "github.com/bugsnag/bugsnag-go"
	log "github.com/sirupsen/logrus"
)

var client Statter
var config *StatterConfig

type StatterConfig struct {
	EnvName              string
	StuckFunctionTimeout time.Duration
	MockingEnabled       bool
}

func (m *StatterConfig) BaseTags() string {
	return ",env=" + config.EnvName
}

type Statter interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Histogram(bucket string, value interface{})
	Unique(bucket string, value string)
	Close()
}

func Close() {
	if client == nil {
		return
	}
	client.Close()
}

func Disable() {
	config = checkConfig(nil)
	client = newMockStatter(true)
}

func Init(addr string, prefix string, cfg *StatterConfig) error {
	config = checkConfig(cfg)
	if config.MockingEnabled {
		// init a mock statter instead of real statsd client
		client = newMockStatter(false)
		return nil
	}
	if c, err := statsd.New(statsd.Address(addr),
		statsd.Prefix(prefix),
		statsd.ErrorHandler(errHandler)); err != nil {
		return err
	} else {
		client = c
	}
	return nil
}

func checkConfig(cfg *StatterConfig) *StatterConfig {
	if cfg == nil {
		cfg = &StatterConfig{}
	}
	if cfg.StuckFunctionTimeout < time.Second {
		cfg.StuckFunctionTimeout = 5 * time.Minute
	}
	if len(cfg.EnvName) == 0 {
		cfg.EnvName = "local"
	}
	return cfg
}

func errHandler(err error) {
	bugsnag.Notify(fmt.Errorf("statsd error: %v", err))
}

type proxy struct {
	client Statter
}

func (p proxy) Timing(bucket string, d time.Duration) {
	if d < 0 {
		return
	}
	p.client.Timing(bucket, int(d/time.Millisecond))
}

func (p proxy) Gauge(bucket string, value int) {
	p.client.Gauge(bucket, value)
}

func newMockStatter(noop bool) Statter {
	return &mockStatter{
		noop: noop,
		fields: log.Fields{
			"module": "mock_statter",
		},
	}
}

type mockStatter struct {
	fields log.Fields
	noop   bool
}

func (s *mockStatter) Count(bucket string, n interface{}) {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, n)
}

func (s *mockStatter) Increment(bucket string) {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("Bucket %s", bucket)
}

func (s *mockStatter) Gauge(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Timing(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Histogram(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Unique(bucket string, value string) {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Close() {
	if s.noop {
		return
	}
	log.WithFields(logging.WithFn(s.fields)).Debugf("closed at %s", time.Now())
}
