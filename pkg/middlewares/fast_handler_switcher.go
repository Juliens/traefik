package middlewares

import (
	"github.com/traefik/traefik/v2/pkg/safe"
	"github.com/valyala/fasthttp"
)

// FastHTTPHandlerSwitcher allows hot switching of http.ServeMux.
type FastHTTPHandlerSwitcher struct {
	handler *safe.Safe
}

// NewHandlerSwitcher builds a new instance of FastHTTPHandlerSwitcher.
func NewFastHandlerSwitcher(newHandler fasthttp.RequestHandler) (hs *FastHTTPHandlerSwitcher) {
	return &FastHTTPHandlerSwitcher{
		handler: safe.New(newHandler),
	}
}

func (h *FastHTTPHandlerSwitcher) Serve(ctx *fasthttp.RequestCtx) {
	handlerBackup := h.handler.Get().(fasthttp.RequestHandler)
	handlerBackup(ctx)
}

// GetHandler returns the current http.ServeMux.
func (h *FastHTTPHandlerSwitcher) GetHandler() (newHandler fasthttp.RequestHandler) {
	handler := h.handler.Get().(fasthttp.RequestHandler)
	return handler
}

// UpdateHandler safely updates the current http.ServeMux with a new one.
func (h *FastHTTPHandlerSwitcher) UpdateHandler(newHandler fasthttp.RequestHandler) {
	h.handler.Set(newHandler)
}
