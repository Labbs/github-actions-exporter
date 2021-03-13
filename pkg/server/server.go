package server

import (
	"log"
	"strconv"

	"github.com/fasthttp/router"
	"github.com/urfave/cli/v2"
	"github.com/valyala/fasthttp"

	"github-actions-exporter/pkg/config"
	"github-actions-exporter/pkg/metrics"
)

func RunServer(ctx *cli.Context) error {
	metrics.InitMetrics()

	r := router.New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("/metrics")
	})
	r.GET("/metrics", prometheusHandler())
	r.GET("/debug/pprof/", pprofHandlerIndex)
	r.GET("/debug/pprof/cmdline", pprofHandlerCmdline)
	r.GET("/debug/pprof/profile", pprofHandlerIndex)
	r.GET("/debug/pprof/trace", pprofHandlerTrace)
	r.GET("/debug/pprof/{profile}", pprofHandlerIndex)

	log.Print("exporter listening on 0.0.0.0:" + strconv.Itoa(config.Port))
	return fasthttp.ListenAndServe(":"+strconv.Itoa(config.Port), r.Handler)
}
