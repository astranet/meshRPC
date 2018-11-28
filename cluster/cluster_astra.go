package cluster

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/astranet/astranet"
	"github.com/astranet/galaxy/logging"
)

type AstraOptions struct {
	Tags  []string
	Nodes []string
	Debug bool
}

func checkAstraOptions(opt *AstraOptions) *AstraOptions {
	if opt == nil {
		opt = &AstraOptions{}
	}
	if len(opt.Tags) == 0 {
		opt.Tags = []string{"default", "local"}
	}
	return opt
}

func NewAstraCluster(service string, opt *AstraOptions) Cluster {
	opt = checkAstraOptions(opt)
	net := astranet.New().Router().WithEnv(opt.Tags...)
	fields := log.Fields{
		"layer":   "cluster",
		"service": service,
		"net_env": strings.Join(opt.Tags, "."),
	}
	log.WithFields(fields).Infoln("new astranet router created")

	if opt.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	c := &astraCluster{
		fields: fields,
		net:    net,
		tags:   opt.Tags,
		dbg:    opt.Debug,

		service: service,
		router:  gin.New(),
	}
	c.initRouter()
	if len(opt.Nodes) > 0 {
		go c.Join(opt.Nodes)
	}
	return c
}

func (a *astraCluster) initRouter() {
	a.router.Use(gin.Logger())
	a.router.GET("/ping", okLoopback())
	a.router.GET("/__heartbeat__", okLoopback())
	a.router.POST("/__error__", errLoopback())
}

type astraCluster struct {
	fields log.Fields
	net    astranet.AstraNet
	tags   []string
	dbg    bool

	service string
	router  *gin.Engine
}

const defaultAstraPort = "11999"

func (a *astraCluster) Join(nodes []string) error {
	var failed []string
	for _, nodeAddr := range nodes {
		if _, _, err := net.SplitHostPort(nodeAddr); err != nil {
			nodeAddr = nodeAddr + ":" + defaultAstraPort
		}
		if err := a.net.Join("tcp4", nodeAddr); err != nil {
			failed = append(failed, nodeAddr)
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("cluster: failed to join nodes: %v", failed)
	}
	a.net.Services()
	return nil
}

func (a *astraCluster) ListenAndServe(addr string) error {
	fnLog := log.WithFields(logging.WithFn(a.fields))
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return newError(err, fmt.Sprintf("cluster: wrong hostport to listen and serve: %v", err))
	}

	fnLog.Infoln("binding", a.service, "as", serviceFQDN(a.service))
	listener, err := a.net.Bind("", serviceFQDN(a.service))
	if err != nil {
		return err
	}
	fnLog.Infoln("ListenAndServe on", addr)
	// expose internal HTTP router to the net using custom listener
	go http.Serve(listener, a.router)

	err = a.net.ListenAndServe("tcp4", net.JoinHostPort(host, port))
	if err != nil {
		return newError(err, fmt.Sprintf("cluster: failed to listen: %v", err))
	}
	if len(host) == 0 || host == "0.0.0.0" {
		host = "localhost"
	}
	a.net.Join("tcp4", net.JoinHostPort(host, port))
	a.net.Services()
	return nil
}

func newHTTPTransport(aNet astranet.AstraNet) *http.Transport {
	return &http.Transport{
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Minute,
		Dial: func(network, addr string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			return aNet.Dial(network, host)
		},
	}
}

func (a *astraCluster) ListenAndServeHTTP(addr string) error {
	fnLog := log.WithFields(logging.WithFn(a.fields))
	fnLog.Infoln("ListenAndServeHTTP on", addr)
	return http.ListenAndServe(addr, &httputil.ReverseProxy{
		Transport:     newHTTPTransport(a.net),
		FlushInterval: time.Millisecond * 10,
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = serviceFQDN(a.service)
		},
	})
}

func (a *astraCluster) Publish(spec HandlerSpec) error {
	endpoints, err := reflectEndpoints(spec)
	if err != nil {
		err = newError(err, fmt.Sprintf("cluster: failed to inspect provided HandlerSpec: %v", err))
		return err
	} else if len(endpoints) == 0 {
		err := fmt.Errorf("cluster: provided HandlerSpec doesn't have any public http.HandlerFunc")
		return err
	}
	for _, e := range endpoints {
		target := *e
		if len(target.Methods) == 0 {
			a.router.Any(target.Path, func(c *gin.Context) {
				target.Handler(c)
			})
			continue
		}
		for _, m := range target.Methods {
			a.router.Handle(m, target.Path, func(c *gin.Context) {
				target.Handler(c)
			})
		}
	}
	return nil
}

func (a *astraCluster) Wait(ctx context.Context, specs ...HandlerSpec) error {
	fnLog := log.WithFields(logging.WithFn(a.fields))
	readyNames := make(map[string]struct{}, len(specs))
	readyMux := new(sync.RWMutex)
	allNames := make([]string, 0, len(specs))
	allMux := new(sync.RWMutex)

	doneC := make(chan struct{})
	wg := new(sync.WaitGroup)
	wg.Add(len(specs))
	go func() {
		wg.Wait()
		close(doneC)
	}()
	for _, spec := range specs {
		go func(spec HandlerSpec) {
			defer wg.Done()
			name, err := reflectServiceName(spec)
			if err != nil {
				fnLog.WithError(err).Warnln("skipping wait on %s", spec)
				return
			}
			allMux.Lock()
			allNames = append(allNames, name)
			allMux.Unlock()
			for {
				switch state := a.pingService(ctx, name); state {
				case stateCancelled, stateTimeout:
					return
				case stateReady:
					readyMux.Lock()
					readyNames[name] = struct{}{}
					readyMux.Unlock()
					return
				default: // retry upon error
					time.Sleep(time.Second)
					continue
				}
			}
		}(spec)
	}
	select {
	case <-doneC:
		// all ok
		return nil
	case <-ctx.Done():
		readyMux.RLock()
		allMux.RLock()
		notReady := make([]string, 0, len(allNames))
		for _, name := range allNames {
			if _, ok := readyNames[name]; !ok {
				notReady = append(notReady, name)
			}
		}
		allMux.RUnlock()
		readyMux.RUnlock()
		err := fmt.Errorf("wait error: services failed to respond in time: %s", strings.Join(notReady, ","))
		return err
	}
}

type serviceState int

const (
	stateReady     serviceState = 1
	stateTimeout   serviceState = 2
	stateCancelled serviceState = 3
	stateError     serviceState = 4
)

// pingService returns the serviceState for the provided service name.
func (a *astraCluster) pingService(ctx context.Context, serviceName string) serviceState {
	fnLog := log.WithFields(logging.WithFn(a.fields))
	cli := &http.Client{
		Transport: newHTTPTransport(a.net),
	}
	u := fmt.Sprintf("http://%s/ping", serviceFQDN(serviceName))
	req, _ := http.NewRequest("GET", u, nil)
	req = req.WithContext(ctx)
	resp, err := cli.Do(req)
	if err != nil {
		fnLog.Debugln("pingService:", serviceName, err)
		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				return stateCancelled
			}
			return stateTimeout
		default:
			return stateError
		}
	}
	if resp.StatusCode != http.StatusOK {
		return stateError
	}
	return stateReady
}

func (a *astraCluster) NewClient(spec HandlerSpec, nameOpt ...string) Client {
	var fnName string
	if len(nameOpt) > 0 {
		fnName = nameOpt[0]
	}
	endpoint, err := reflectEndpointInfo(spec, fnName)
	if err != nil {
		panic(fmt.Errorf("cluster: failed to reflect target http.HandlerFunc: %v", err))
	}
	cli := &astraClient{
		net:       a.net,
		endpoint:  endpoint,
		localhost: serviceFQDN(a.service),
		cli: &http.Client{
			Transport: newHTTPTransport(a.net),
		},
	}
	if len(fnName) > 0 {
		cli.enableReverseProxy()
	}
	return cli
}

func (a *astraClient) enableReverseProxy() {
	a.Handler = &httputil.ReverseProxy{
		Transport:     newHTTPTransport(a.net),
		FlushInterval: time.Millisecond * 10,
		Director: func(req *http.Request) {
			reportErr := func(status int, err error) {
				req.Method = "POST"
				req.URL, _ = url.ParseRequestURI("http://" + a.localhost + "/__error__")
				v, _ := json.Marshal(proxyError{
					Status:  status,
					Message: err.Error(),
				})
				req.ContentLength = int64(len(v))
				req.Body = ioutil.NopCloser(bytes.NewReader(v))
			}
			if !a.endpoint.MethodAllowed(req.Method) {
				reportErr(400,
					fmt.Errorf("cluster client: method %s not allowed for %s: must be %s",
						req.Method, req.URL.Path, strings.Join(a.endpoint.Methods, ",")))
				return
			}
			req.URL, _ = url.Parse("http://" + serviceFQDN(a.endpoint.Service) + a.endpoint.Path)
		},
	}
}

type astraClient struct {
	http.Handler

	net       astranet.AstraNet
	endpoint  *EndpointInfo
	localhost string
	cli       *http.Client
}

func (a *astraClient) Use(fnName string) Client {
	if len(fnName) == 0 || a.endpoint == nil {
		return a
	}
	endpoint := *a.endpoint
	endpoint.Path = rewritePath(fnName, endpoint.Path)
	cli := &astraClient{
		net:       a.net,
		Handler:   a.Handler,
		endpoint:  &endpoint,
		localhost: a.localhost,
		cli:       a.cli,
	}
	if cli.Handler == nil {
		cli.enableReverseProxy()
	}
	return cli
}

func (a *astraClient) Do(req *http.Request) (*http.Response, error) {
	if a.cli == nil {
		return nil, nil
	} else if a.endpoint == nil {
		return a.cli.Do(req)
	}
	path := a.endpoint.Path
	if fnName := req.URL.String(); len(fnName) > 0 {
		if !a.endpoint.IsValidHandler(fnName) {
			err := fmt.Errorf("cluster client: %s is not a valid http.HandlerFunc or not exists", fnName)
			return nil, err
		}
		path = rewritePath(fnName, path)
	}
	req.URL, _ = url.Parse("http://" + serviceFQDN(a.endpoint.Service) + path)
	if !a.endpoint.MethodAllowed(req.Method) {
		err := fmt.Errorf("cluster client: method %s not allowed for %s: must be %s",
			req.Method, path, strings.Join(a.endpoint.Methods, ","))
		return nil, err
	}
	return a.cli.Do(req)
}

func rewritePath(fnName string, path string) string {
	if len(fnName) == 0 || strings.ContainsAny(fnName, "/") {
		return path
	}
	parts := strings.Split(path, "/")
	parts[len(parts)-1] = fnName
	return strings.Join(parts, "/")
}

func okLoopback() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatus(200)
	}
}

func errLoopback() gin.HandlerFunc {
	return func(c *gin.Context) {
		var e proxyError
		if err := c.BindJSON(&e); err != nil {
			c.Status(500)
			return
		}
		if e.Status > 0 && e.Status != http.StatusOK {
			c.String(e.Status, "%s", e.Message)
			return
		}
	}
}

type proxyError struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

func randTag(n int) string {
	buf := make([]byte, n)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}

func serviceFQDN(service string) string {
	return "galaxy." + service
}
