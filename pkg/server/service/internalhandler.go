package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/valyala/fasthttp"
)

type serviceManager interface {
	BuildHTTP(rootCtx context.Context, serviceName string) (fasthttp.RequestHandler, error)
	LaunchHealthCheck()
}

// InternalHandlers is the internal HTTP handlers builder.
type InternalHandlers struct {
	api        fasthttp.RequestHandler
	dashboard  fasthttp.RequestHandler
	rest       http.Handler
	prometheus http.Handler
	ping       http.Handler
	acmeHTTP   http.Handler
	serviceManager
}

// NewInternalHandlers creates a new InternalHandlers.
func NewInternalHandlers(next serviceManager, apiHandler fasthttp.RequestHandler, rest, metricsHandler, pingHandler http.Handler, dashboard fasthttp.RequestHandler, acmeHTTP http.Handler) *InternalHandlers {
	return &InternalHandlers{
		api:            apiHandler,
		dashboard:      dashboard,
		rest:           rest,
		prometheus:     metricsHandler,
		ping:           pingHandler,
		acmeHTTP:       acmeHTTP,
		serviceManager: next,
	}
}

// BuildHTTP builds an HTTP handler.
func (m *InternalHandlers) BuildHTTP(rootCtx context.Context, serviceName string) (fasthttp.RequestHandler, error) {
	if !strings.HasSuffix(serviceName, "@internal") {
		return m.serviceManager.BuildHTTP(rootCtx, serviceName)
	}

	internalHandler, err := m.get(serviceName)
	if err != nil {
		return nil, errors.New("Wrong internal")
	}

	return internalHandler, nil
}

func (m *InternalHandlers) get(serviceName string) (fasthttp.RequestHandler, error) {
	switch serviceName {
	case "noop@internal":
		return func(ctx *fasthttp.RequestCtx) {
			ctx.SetStatusCode(http.StatusTeapot)
		}, nil

	// case "acme-http@internal":
	// 	if m.acmeHTTP == nil {
	// 		return nil, errors.New("HTTP challenge is not enabled")
	// 	}
	// 	return m.acmeHTTP, nil

	case "api@internal":
		if m.api == nil {
			return nil, errors.New("api is not enabled")
		}
		return m.api, nil

	// case "dashboard@internal":
	// 	if m.dashboard == nil {
	// 		return nil, errors.New("dashboard is not enabled")
	// 	}
	// 	return m.dashboard, nil
	//
	// case "rest@internal":
	// 	if m.rest == nil {
	// 		return nil, errors.New("rest is not enabled")
	// 	}
	// 	return m.rest, nil
	//
	// case "ping@internal":
	// 	if m.ping == nil {
	// 		return nil, errors.New("ping is not enabled")
	// 	}
	// 	return m.ping, nil
	//
	// case "prometheus@internal":
	// 	if m.prometheus == nil {
	// 		return nil, errors.New("prometheus is not enabled")
	// 	}
	// 	return m.prometheus, nil

	default:
		return nil, fmt.Errorf("unknown internal service %s", serviceName)
	}
}
