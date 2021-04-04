package chain

import (
	"context"

	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/server/router/alice"
	"github.com/valyala/fasthttp"
)

const (
	typeName = "Chain"
)

type chainBuilder interface {
	BuildChain(ctx context.Context, middlewares []string) *alice.Chain
}

// New creates a chain middleware.
func New(ctx context.Context, next fasthttp.RequestHandler, config dynamic.Chain, builder chainBuilder, name string) (fasthttp.RequestHandler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")

	middlewareChain := builder.BuildChain(ctx, config.Middlewares)
	return middlewareChain.Then(next)
}
