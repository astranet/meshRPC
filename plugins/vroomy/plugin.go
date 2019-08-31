package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vroomy/plugins"

	"github.com/Hatch1fy/httpserve"
	"github.com/astranet/meshRPC/cluster"
)

var (
	cfg  *PluginConfig
	log  *logrus.Entry
	mesh cluster.Cluster
)

func init() {
	log = logrus.WithFields(logrus.Fields{
		"plugin": "meshrpc",
	})
}

func OnInit(p *plugins.Plugins, env map[string]string) error {
	cfg = checkPluginConfig(configFromEnv(env))
	mesh = cluster.NewAstraCluster("vroomy", &cluster.AstraOptions{
		Tags: []string{
			cfg.ClusterName,
		},
		Nodes: cfg.ClusterNodes,
		Debug: cfg.Debug,
	})
	// Listen on a TCP address, this address can be used
	// by other peers to discover each other in this cluster.
	if err := mesh.ListenAndServe(
		fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.ListenPort),
	); err != nil {
		return err
	}

	// TODO: implement optional waiting for service dependencies
	// How a plugin would know all routes it is used for?
	// Observation: OnInit is done before Route() closure is evaluated on start.
	return nil
}

// Backend will return the underlying instance of meshrpc cluster.
func Backend() interface{} {
	return mesh
}

// Ping will check if the service is alive.
func Ping(c *httpserve.Context) (res httpserve.Response) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()

	ts := time.Now()
	state := mesh.PingService(ctx, c.Param("service"))

	var code int
	switch state {
	case cluster.StateOK:
		code = http.StatusOK
		return httpserve.NewJSONResponse(code, time.Since(ts).String())
	case cluster.StateTimeout, cluster.StateCanceled:
		code = http.StatusGatewayTimeout
		return httpserve.NewJSONResponse(code, errors.New(time.Since(ts).String()))
	default:
		code = http.StatusBadGateway
		return httpserve.NewJSONResponse(code, errors.New(time.Since(ts).String()))
	}
}

var ErrInsufficientParams = errors.New("insufficient params count")

// Route will route the request to HTTP handler using mesh cluster instance.
func Route(params ...string) (handler httpserve.Handler, err error) {
	if len(params) < 3 {
		err = errors.Wrap(ErrInsufficientParams, "method requires 3 params")
		return nil, err
	}
	serviceName := params[0]
	if len(serviceName) == 0 {
		err = errors.Wrap(ErrInsufficientParams, "no service name provided")
		return nil, err
	}
	handlerName := params[1]
	if len(handlerName) == 0 {
		err = errors.Wrap(ErrInsufficientParams, "no handler name provided")
		return nil, err
	}
	methodName := params[2]
	if len(methodName) == 0 {
		err = errors.Wrap(ErrInsufficientParams, "no method name provided")
		return nil, err
	}
	// TODO: cache this client with given params?
	cli := mesh.NewClient(serviceName, handlerName).Use(methodName)

	handler = func(c *httpserve.Context) (res httpserve.Response) {
		cli.ServeHTTP(c.Writer, c.Request)
		return httpserve.NewAdoptResponse()
	}
	return handler, nil
}

// Close will close the plugin
func Close() error {
	log.Debugln("closing plugin")
	return nil
}
