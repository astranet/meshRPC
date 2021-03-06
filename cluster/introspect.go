package cluster

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astranet/httpserve"
)

type HandlerSpec interface{}

type HTTPMethodsSpec interface {
	HTTPMethodsMap() map[string][]string
}

func httpMethodsOf(spec HandlerSpec) map[string][]string {
	if s, ok := spec.(HTTPMethodsSpec); ok {
		return s.HTTPMethodsMap()
	}
	return nil
}

func reflectEndpoints(serviceName string, spec HandlerSpec) ([]*EndpointInfo, error) {
	if spec == nil {
		return nil, errors.New("reflectEndpoints: spec is nil")
	}
	specTyp := reflect.TypeOf(spec)
	specVal := reflect.ValueOf(spec)
	_, handlerName := ifaceToPkgName(specTyp)
	httpMethods := httpMethodsOf(spec)

	n := specTyp.NumMethod()
	endpoints := make([]*EndpointInfo, 0, n)
	for i := 0; i < n; i++ {
		m := specTyp.Method(i)
		if isHandlerFunc(m.Type) {
			handlerFn := specVal.MethodByName(m.Name).Interface().(func(c *httpserve.Context) httpserve.Response)
			endpoint := &EndpointInfo{
				Service: serviceName,
				Path:    fmt.Sprintf("/%s/%s", handlerName, m.Name),
				Handler: handlerFn,
				SpecTyp: specTyp,
			}
			if methods, ok := httpMethods["*"]; ok {
				// has a record that matches all endpoints
				endpoint.Methods = methods
			} else {
				// otherwise use specific methods
				endpoint.Methods = httpMethods[m.Name]
			}
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints, nil
}

func reflectEndpointInfo(serviceName string, spec HandlerSpec, fnName string) (*EndpointInfo, error) {
	if spec == nil {
		return nil, errors.New("reflectEndpointInfo: spec is nil")
	}
	if handlerName, ok := spec.(string); ok {
		endpoint := &EndpointInfo{
			Service: serviceName,
			Path:    fmt.Sprintf("/%s/%s", handlerName, fnName),
		}
		return endpoint, nil
	}
	specTyp := reflect.TypeOf(spec)
	if len(fnName) > 0 {
		m, ok := specTyp.MethodByName(fnName)
		if !ok {
			err := fmt.Errorf("reflectEndpointInfo: spec doesnt't have method %s", fnName)
			return nil, err
		}
		if !isHandlerFunc(m.Type) {
			err := fmt.Errorf("reflectEndpointInfo: method %s is not a http.HandlerFunc", fnName)
			return nil, err
		}
	}
	_, handlerName := ifaceToPkgName(specTyp)
	httpMethods := httpMethodsOf(spec)
	endpoint := &EndpointInfo{
		Service: serviceName,
		Path:    fmt.Sprintf("/%s/%s", handlerName, fnName),
		SpecTyp: specTyp,
	}
	if methods, ok := httpMethods["*"]; ok {
		// has a record that matches all endpoints
		endpoint.Methods = methods
	} else {
		// otherwise use specific methods
		endpoint.Methods = httpMethods[fnName]
	}
	return endpoint, nil
}

type EndpointInfo struct {
	Service string
	Path    string
	Methods []string
	SpecTyp reflect.Type
	Handler func(c *httpserve.Context) httpserve.Response
}

func (e *EndpointInfo) MethodAllowed(method string) bool {
	if e == nil || len(e.Methods) == 0 {
		return true // no constraints
	}
	method = strings.ToUpper(method)
	for _, m := range e.Methods {
		if m == method {
			return true
		}
	}
	return false
}

func (e *EndpointInfo) IsValidHandler(name string) bool {
	if e.SpecTyp == nil {
		return true
	}
	fn, exists := e.SpecTyp.MethodByName(name)
	if !exists {
		return false
	}
	return isHandlerFunc(fn.Type)
}

var httpContextTyp = reflect.TypeOf((*httpserve.Context)(nil))

// ifacePkgName returns package path for an interface type, and the type's name.
func ifaceToPkgName(typ reflect.Type) (pkgName string, typName string) {
	implTyp := typ.Elem()
	pkgName = implTyp.PkgPath()
	typName = implTyp.Name()
	return
}

// isHandlerFunc basically checks method to match httpserve.HandlerFunc
// func(*httpserve.Context) httpserve.Response
func isHandlerFunc(fn reflect.Type) bool {
	if fn.NumIn() != 2 {
		return false
	}
	if fn.NumOut() != 1 {
		return false
	}
	if fn.In(1) != httpContextTyp {
		return false
	}
	return true
}
