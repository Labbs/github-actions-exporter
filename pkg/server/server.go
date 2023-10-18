package server

import (
	"strconv"

	"github.com/fasthttp/router"
	"github.com/urfave/cli/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/spendesk/github-actions-exporter/pkg/config"
	"github.com/spendesk/github-actions-exporter/pkg/metrics"
)

// RunServer - run http server for expose metrics
func RunServer(ctx *cli.Context, logger *zap.SugaredLogger) error {
	metrics.InitMetrics(logger)

	r := router.New()
	r.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.WriteString("/metrics")
	})
	r.GET("/metrics", prometheusHandler())

	if config.Debug {
		r.GET("/debug/pprof/", pprofHandlerIndex)
		r.GET("/debug/pprof/cmdline", pprofHandlerCmdline)
		r.GET("/debug/pprof/profile", pprofHandlerIndex)
		r.GET("/debug/pprof/trace", pprofHandlerTrace)
		r.GET("/debug/pprof/{profile}", pprofHandlerIndex)
	}

	logger.Info("exporter listening on 0.0.0.0:" + strconv.Itoa(config.Port))
	return fasthttp.ListenAndServe(":"+strconv.Itoa(config.Port), r.Handler)
}
