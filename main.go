package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/spendesk/github-actions-exporter/pkg/config"
	"github.com/spendesk/github-actions-exporter/pkg/logging"
	"github.com/spendesk/github-actions-exporter/pkg/server"
)

var (
	version = "development"
	logger  = logging.GetLogger()
)

func main() {
	app := cli.NewApp()
	app.Name = "github-actions-exporter"
	app.Flags = config.InitConfiguration()
	app.Version = version
	app.Action = server.RunServer

	err := app.Run(os.Args)

	if err != nil {
		logger.Fatal(err)
	}
}
