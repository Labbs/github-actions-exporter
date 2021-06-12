package config

import "github.com/urfave/cli/v2"

var (
	// Github - github configuration
	Github struct {
		Token             string
		AppInstallationID int64
		AppID             int64
		AppPrivateKey     string
		Refresh           int64
		Repositories      cli.StringSlice
		Organizations     cli.StringSlice
		APIURL            string
	}
	// Port - http port used by fasthttp
	Port int
	// Debug - option for debug the exporter
	Debug bool
)

// InitConfiguration - set configuration from env vars or command parameters
func InitConfiguration() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:        "port",
			Aliases:     []string{"p"},
			EnvVars:     []string{"PORT"},
			Value:       9999,
			Usage:       "Exporter port",
			Destination: &Port,
		},
		&cli.StringFlag{
			Name:        "github_token",
			Aliases:     []string{"gt"},
			EnvVars:     []string{"GITHUB_TOKEN"},
			Usage:       "Github Personal Token",
			Destination: &Github.Token,
		},
		&cli.Int64Flag{
			Name:        "github_app_installation_id",
			Aliases:     []string{"gaiid"},
			EnvVars:     []string{"GITHUB_APP_INSTALLATION_ID"},
			Usage:       "Github App Installation ID",
			Destination: &Github.AppInstallationID,
		},
		&cli.Int64Flag{
			Name:        "github_app_id",
			Aliases:     []string{"gaid"},
			EnvVars:     []string{"GITHUB_APP_ID"},
			Usage:       "Github App ID",
			Destination: &Github.AppID,
		},
		&cli.StringFlag{
			Name:        "github_private_key",
			Aliases:     []string{"gapk"},
			EnvVars:     []string{"GITHUB_APP_PRIVATE_KEY"},
			Usage:       "Github App Private Key",
			Destination: &Github.AppPrivateKey,
		},
		&cli.Int64Flag{
			Name:        "github_refresh",
			Aliases:     []string{"gr"},
			EnvVars:     []string{"GITHUB_REFRESH"},
			Value:       30,
			Usage:       "Refresh time Github Pipelines status in sec",
			Destination: &Github.Refresh,
		},
		&cli.StringFlag{
			Name:        "github_api_url",
			Aliases:     []string{"url"},
			EnvVars:     []string{"GITHUB_API_URL"},
			Value:       "api.github.com",
			Usage:       "Github API URL (primarily designed for Github Enterprise use cases)",
			Destination: &Github.APIURL,
		},
		&cli.StringSliceFlag{
			Name:        "github_orgas",
			Aliases:     []string{"go"},
			EnvVars:     []string{"GITHUB_ORGAS"},
			Usage:       "List all organizations you want get informations. Format <orga>,<orga2>,<orga3> (like test,test2)",
			Destination: &Github.Organizations,
		},
		&cli.StringSliceFlag{
			Name:        "github_repos",
			Aliases:     []string{"grs"},
			EnvVars:     []string{"GITHUB_REPOS"},
			Usage:       "List all repositories you want get informations. Format <orga>/<repo>,<orga>/<repo2>,<orga>/<repo3> (like test/test)",
			Destination: &Github.Repositories,
		},
		&cli.BoolFlag{
			Name:        "debug_profile",
			EnvVars:     []string{"DEBUG_PROFILE"},
			Usage:       "Expose pprof information on /debug/pprof/",
			Destination: &Debug,
		},
	}
}
