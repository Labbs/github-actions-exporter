package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/spendesk/github-actions-exporter/pkg/config"
	"github.com/spendesk/github-actions-exporter/pkg/logging"
	"github.com/spendesk/github-actions-exporter/pkg/server"
)

var (
	version = "development"
	logger  *zap.SugaredLogger
)

func main() {
	app := &cli.App{
		Name:  "github-actions-exporter",
		Flags: config.InitConfiguration(),
		Action: func(ctx *cli.Context) error {
			logger = logging.InitLogger()
			server.RunServer(ctx, logger)
			return nil
		},
		Version: version,
	}

	err := app.Run(os.Args)

	if err != nil {
		logger.Fatal(err)
	}
}
