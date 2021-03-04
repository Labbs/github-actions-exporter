/*
Package config get configuration from environment variables or/and flags
*/
package config

import (
	"github.com/urfave/cli"
)

var (
	Github struct {
		Token         string
		Refresh       int64
		Repositories  cli.StringSlice
		Organizations cli.StringSlice
		ApiUri	string
	}
	Port int
)

// NewContext => set configuration from env vars or command parameters
func NewContext() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "github_token, gt",
			Destination: &Github.Token,
			EnvVar:      "GITHUB_TOKEN",
			Usage:       "Github Personal Token",
		},
		cli.Int64Flag{
			Name:        "github_refresh, gr",
			Value:       30,
			Destination: &Github.Refresh,
			EnvVar:      "GITHUB_REFRESH",
			Usage:       "Refresh time Github Pipelines status in sec",
		},
		cli.StringSliceFlag{
			Name:   "github_orgas, go",
			Value:  &Github.Organizations,
			EnvVar: "GITHUB_ORGAS",
			Usage:  "List all organizations you want get informations. Format <orga>,<orga2>,<orga3> (like test,test2)",
		},
		cli.StringSliceFlag{
			Name:   "github_repos, grs",
			Value:  &Github.Repositories,
			EnvVar: "GITHUB_REPOS",
			Usage:  "List all repositories you want get informations. Format <orga>/<repo>,<orga>/<repo2>,<orga>/<repo3> (like test/test)",
		},
		cli.StringFlag{
			Name:        "github_api_uri, uri",
			Value:       "api.github.com",
			Destination: &Github.ApiUri,
			EnvVar:      "GITHUB_API_URI",
			Usage:       "Github API URI (primarily designed for Github Enterprise use cases)",
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       9999,
			Destination: &Port,
			EnvVar:      "PORT",
			Usage:       "Exporter port",
		},
	}
}
