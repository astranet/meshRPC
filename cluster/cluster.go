package cluster

import (
	"context"
	"net/http"
)

type ErrHandlerFunc func(err error)

type Cluster interface {
	// ListenAndServe starts listening for TCP connections on the provided address
	// for interprocess communication. This is how services connect to each other,
	// this is actually the only port that needs to be open, if service is not
	// exposed outside the cluster. In case of astranet implementation, all TCP
	// connections between instances will be kept alive bound to the provided port.
	ListenAndServe(addr string) error
	// ListenAndServeHTTP allows to expose a debugging endpoint that is just an HTTP
	// reverse proxy into virtual network. Since we decided that our RPC proto and
	// routing is based on HTTP, you can use curl to make RPC calls using that
	// endpoint.
	//
	// Example: curl -X GET http://meshrpc.greeter/ping where
	// meshrpc.greeter is a service registered in the discovery table.
	// (Note: since it’s a virtual host, it may have different backends with a LB).
	ListenAndServeHTTP(addr string) error
	// Join is for joining to the service discovery cluster. In case of P2P
	// mechanics of astranet, it will try to connect to the listed nodes,
	// establish TCP connections, get routing and services info and share own
	// routing and services info.
	Join(nodes []string) error
	// Publish exposes provided Handler to the private net under own service name.
	// Eeach func that is http.HandlerFunc is being mapped to URI and connected
	// to the internal HTTP router (Gin). If HandlerSpec implements HTTPMethodsSpec,
	// a methods map is used, otherwise all methods are accepted and it's a client
	// responsibility to use the correct one.
	Publish(spec HandlerSpec) error
	// Wait blocks until all of the mapped services in spec become available or context
	// is cancelled. It’s useful to guarantee that all dependencies are up and
	// alive before running the service that uses them. Returns an error if not
	// services became available upon context cancellation. Specs contains a mapping
	// ServiceName -> HandlerSpec.
	Wait(ctx context.Context, specs map[string]HandlerSpec) error
	// NewClient returns a new cluster.Client for accessing published endpoints of
	// services. The mandatory HandlerSpec argument accepts any interface type has
	// http.HandlerFunc methods, like in Publish, but the its implementation is not
	// required. If provided with a specific function name, cluster.Client locks
	// itself to access that specific method from HandlerSpec and starts to conform
	// the http.Handler interface. Otherwise it may be used to access any method
	// from HandlerSpec, using the same Do(req *http.Request) method as with
	// http.Client.
	NewClient(serviceName string, spec HandlerSpec, fn ...string) Client
}

// Client is used for accessing published endpoints of services. It combines a
// fail-safe proxy for exposing remotely published http.HandlerFunc as local
// http.Handler, also a http.Client like function Do to send custom requests to
// remote methods directly.
type Client interface {
	http.Handler

	Do(req *http.Request) (*http.Response, error)
	Use(fnName string) Client
}

type Error struct {
	Desc string
	Err  error
}

func (c *Error) Error() string {
	if c == nil || c.Err == nil {
		return ""
	}
	if len(c.Desc) > 0 {
		return c.Desc
	}
	return c.Err.Error()
}

func newError(err error, desc string) error {
	if err == nil {
		return nil
	}
	return &Error{
		Desc: desc,
		Err:  err,
	}
}
