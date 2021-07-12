package config

import "github.com/urfave/cli/v2"

var (
	// Github - github configuration
	Github struct {
		AppID             int64  `split_words:"true"`
		AppInstallationID int64  `split_words:"true"`
		AppPrivateKey     string `split_words:"true"`
		Token             string
		Refresh           int64
		Repositories      cli.StringSlice
		Organizations     cli.StringSlice
		APIURL            string			
	}
	Port  int
	Debug bool
	EnterpriseName    string
	WorkflowFields    string
)

// InitConfiguration - set configuration from env vars or command parameters
func InitConfiguration() []cli.Flag {
	return []cli.Flag{
		&cli.Int64Flag{
			Name:        "app_id",
			Aliases:     []string{"gai"},
			EnvVars:     []string{"GITHUB_APP_ID"},
			Usage:       "Github App Id",
			Destination: &Github.AppID,
		},
		&cli.Int64Flag{
			Name:        "app_installation_id",
			Aliases:     []string{"gii"},
			EnvVars:     []string{"GITHUB_APP_INSTALLATION_ID"},
			Usage:       "Github App Installation Id",
			Destination: &Github.AppInstallationID,
		},
		&cli.StringFlag{
			Name:        "app_private_key",
			Aliases:     []string{"gpk"},
			EnvVars:     []string{"GITHUB_APP_PRIVATE_KEY"},
			Usage:       "Github App Private Key",
			Destination: &Github.AppPrivateKey,
		},
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
		&cli.StringFlag{
			Name:        "enterprise_name",
			EnvVars:     []string{"ENTERPRISE_NAME"},
			Usage:       "Enterprise name. Needed for enterprise endpoints (/enterprises/{ENTERPRISE_NAME}/*)",
			Destination: &EnterpriseName,
			Value:       "",
		},
		&cli.StringFlag{
			Name:        "export_fields",
			EnvVars:     []string{"EXPORT_FIELDS"},
			Usage:       "A comma separated list of fields for workflow metrics that should be exported",
			Value:       "repo,id,node_id,head_branch,head_sha,run_number,workflow_id,workflow,event,status",
			Destination: &WorkflowFields,
		},
	}
}
