package requestdecorator

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/containous/alice"
	"github.com/traefik/traefik/v2/pkg/types"
	"github.com/valyala/fasthttp"
)

const (
	canonicalKey key = "canonical"
	flattenKey   key = "flatten"
)

type key string

// RequestDecorator is the struct for the middleware that adds the CanonicalDomain of the request Host into a context for later use.
type RequestDecorator struct {
	hostResolver *Resolver
}

// New creates a new request host middleware.
func New(hostResolverConfig *types.HostResolverConfig) *RequestDecorator {
	requestDecorator := &RequestDecorator{}
	if hostResolverConfig != nil {
		requestDecorator.hostResolver = &Resolver{
			CnameFlattening: hostResolverConfig.CnameFlattening,
			ResolvConfig:    hostResolverConfig.ResolvConfig,
			ResolvDepth:     hostResolverConfig.ResolvDepth,
		}
	}
	return requestDecorator
}

func (r *RequestDecorator) Serve(next fasthttp.RequestHandler) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		host := types.CanonicalDomain(parseHost(string(ctx.Host())))
		ctx.SetUserValue(string(canonicalKey), host)

		if r.hostResolver != nil && r.hostResolver.CnameFlattening {
			flatHost := r.hostResolver.CNAMEFlatten(host)
			ctx.SetUserValue(string(flattenKey), flatHost)
		}
		next(ctx)
	}
}

func (r *RequestDecorator) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// host := types.CanonicalDomain(parseHost(req.Host))
	// reqt := req.WithContext(context.WithValue(req.Context(), canonicalKey, host))
	//
	// if r.hostResolver != nil && r.hostResolver.CnameFlattening {
	// 	flatHost := r.hostResolver.CNAMEFlatten(reqt.Context(), host)
	// 	reqt = reqt.WithContext(context.WithValue(reqt.Context(), flattenKey, flatHost))
	// }
	//
	// next(rw, reqt)
}

func parseHost(addr string) string {
	if !strings.Contains(addr, ":") {
		return addr
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

// GetCanonizedHost retrieves the canonized host from the given context (previously stored in the request context by the middleware).
func GetCanonizedHost(ctx context.Context) string {
	if val, ok := ctx.Value(canonicalKey).(string); ok {
		return val
	}

	return ""
}

// GetCNAMEFlatten return the flat name if it is present in the context.
func GetCNAMEFlatten(ctx context.Context) string {
	if val, ok := ctx.Value(flattenKey).(string); ok {
		return val
	}

	return ""
}

// WrapHandler Wraps a ServeHTTP with next to an alice.Constructor.
func WrapHandler(handler *RequestDecorator) alice.Constructor {
	return func(next http.Handler) (http.Handler, error) {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(rw, req, next.ServeHTTP)
		}), nil
	}
}
