package metrics

import (
	"log"
	"net/http"

	"github.com/Spendesk/github-actions-exporter/pkg/config"

	"github.com/google/go-github/v33/github"
	"github.com/gregjones/httpcache"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/oauth2"
)

var (
	client *github.Client
	err    error
)

// InitMetrics - register metrics in prometheus lib and start func for monitor
func InitMetrics() {
	prometheus.MustRegister(runnersGauge)
	prometheus.MustRegister(runnersBusyGauge)
	prometheus.MustRegister(runnersOrganizationGauge)
	prometheus.MustRegister(runnersOrganizationBusyGauge)
	prometheus.MustRegister(workflowRunStatusGauge)
	prometheus.MustRegister(workflowRunStatusDeprecatedGauge)
	prometheus.MustRegister(workflowRunDurationGauge)
	prometheus.MustRegister(workflowBillGauge)

	t := &oauth2.Transport{
		Source: oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Github.Token},
		),
	}

	transport := &httpcache.Transport{
		Transport:           t,
		Cache:               httpcache.NewMemoryCache(),
		MarkCachedResponses: true,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	if config.Github.ApiUrl == "api.github.com" {
		client = github.NewClient(httpClient)
	} else {
		client, err = github.NewEnterpriseClient(config.Github.ApiUrl, config.Github.ApiUrl, httpClient)
		if err != nil {
			log.Fatalln("Github enterprise init error: " + err.Error())
		}
	}

	go workflowCache()

	for {
		if workflows != nil {
			break
		}
	}

	go getBillableFromGithub()
	go getRunnersFromGithub()
	go getRunnersOrganizationFromGithub()
	go getWorkflowRunsFromGithub()
}
