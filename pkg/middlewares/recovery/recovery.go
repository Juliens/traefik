package recovery

import (
	"net/http"
	"runtime"

	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/valyala/fasthttp"
)

const (
	typeName       = "Recovery"
	middlewareName = "traefik-internal-recovery"
)

type recovery struct {
	next fasthttp.RequestHandler
}

// New creates recovery middleware.
func New(next fasthttp.RequestHandler) (fasthttp.RequestHandler, error) {
	// log.FromContext(middlewares.GetLoggerCtx(ctx, middlewareName, typeName)).Debug("Creating middleware")

	r := &recovery{
		next: next,
	}
	return r.Serve, nil
}

func (re *recovery) Serve(ctx *fasthttp.RequestCtx) {
	defer recoverFunc(ctx)
	re.next(ctx)
}
func (re *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// defer recoverFunc(rw, req)
	// re.next.ServeHTTP(rw, req)
}

func recoverFunc(ctx *fasthttp.RequestCtx) {
	if err := recover(); err != nil {
		logger := log.WithoutContext()
		if !shouldLogPanic(err) {
			logger.Debugf("Request has been aborted [%s - %s]: %v", ctx.RemoteAddr(), ctx.URI(), err)
			return
		}

		logger.Errorf("Recovered from panic in HTTP handler [%s - %s]: %+v", ctx.RemoteAddr(), ctx.URI(), err)
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		logger.Errorf("Stack: %s", buf)
		ctx.Error(http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// https://github.com/golang/go/blob/a0d6420d8be2ae7164797051ec74fa2a2df466a1/src/net/http/server.go#L1761-L1775
// https://github.com/golang/go/blob/c33153f7b416c03983324b3e8f869ce1116d84bc/src/net/http/httputil/reverseproxy.go#L284
func shouldLogPanic(panicValue interface{}) bool {
	//nolint:errorlint // false-positive because panicValue is an interface.
	return panicValue != nil && panicValue != http.ErrAbortHandler
}
