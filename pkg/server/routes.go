package server

import (
	"net/http/pprof"
	rtp "runtime/pprof"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var (
	cmdline = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Cmdline)
	profile = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Profile)
	symbol  = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Symbol)
	trace   = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Trace)
	index   = fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Index)
)

func prometheusHandler() fasthttp.RequestHandler {
	return fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
}

func pprofHandlerIndex(ctx *fasthttp.RequestCtx) {
	for _, v := range rtp.Profiles() {
		ppName := v.Name()
		if strings.HasPrefix(string(ctx.Path()), "/debug/pprof/"+ppName) {
			namedHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(pprof.Handler(ppName).ServeHTTP)
			namedHandler(ctx)
			return
		}
	}
	index(ctx)
}

func pprofHandlerCmdline(ctx *fasthttp.RequestCtx) {
	cmdline(ctx)
}

func pprofHandlerProfile(ctx *fasthttp.RequestCtx) {
	profile(ctx)
}

func pprofHandlerSymbol(ctx *fasthttp.RequestCtx) {
	symbol(ctx)
}

func pprofHandlerTrace(ctx *fasthttp.RequestCtx) {
	trace(ctx)
}
