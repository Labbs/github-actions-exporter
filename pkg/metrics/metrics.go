package metrics

import (
	"github-actions-exporter/pkg/config"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/gregjones/httpcache"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/oauth2"
)

var (
	client                   *github.Client
	err                      error
	workflowRunStatusGauge   *prometheus.GaugeVec
	workflowRunDurationGauge *prometheus.GaugeVec
)

// InitMetrics - register metrics in prometheus lib and start func for monitor
func InitMetrics() {
	workflowRunStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_status",
			Help: "Workflow run status",
		},
		strings.Split(config.WorkflowFields, ","),
	)
	workflowRunDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_duration_ms",
			Help: "Workflow run duration (in milliseconds)",
		},
		strings.Split(config.WorkflowFields, ","),
	)
	prometheus.MustRegister(runnersGauge)
	prometheus.MustRegister(runnersOrganizationGauge)
	prometheus.MustRegister(workflowRunStatusGauge)
	prometheus.MustRegister(workflowRunDurationGauge)
	prometheus.MustRegister(workflowBillGauge)
	prometheus.MustRegister(runnersEnterpriseGauge)

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
	go getRunnersEnterpriseFromGithub()
}
