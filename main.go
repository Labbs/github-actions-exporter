package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Spendesk/github-actions-exporter/pkg/config"
	"github.com/Spendesk/github-actions-exporter/pkg/server"
)

var (
	version = "development"
)

func main() {
	app := cli.NewApp()
	app.Name = "github-actions-exporter"
	app.Flags = config.InitConfiguration()
	app.Version = version
	app.Action = server.RunServer

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
