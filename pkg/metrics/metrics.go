package metrics

import (
	"github-actions-exporter/pkg/config"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bradleyfalzon/ghinstallation"
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
	prometheus.MustRegister(runnersOrganizationGauge)
	prometheus.MustRegister(workflowRunStatusGauge)
	prometheus.MustRegister(workflowRunStatusDeprecatedGauge)
	prometheus.MustRegister(workflowRunDurationGauge)
	prometheus.MustRegister(workflowBillGauge)

	var t http.RoundTripper

	if len(config.Github.Token) > 0 {
		t = &oauth2.Transport{
			Source: oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: config.Github.Token},
			),
		}
	} else {
		var tr *ghinstallation.Transport

		if _, err := os.Stat(config.Github.AppPrivateKey); err == nil {
			tr, err = ghinstallation.NewKeyFromFile(http.DefaultTransport, config.Github.AppID, config.Github.AppInstallationID, config.Github.AppPrivateKey)
			if err != nil {
				log.Fatalf("authentication failed: using private key from file %s: %v", config.Github.AppPrivateKey, err)
			}
		} else {
			tr, err = ghinstallation.New(http.DefaultTransport, config.Github.AppID, config.Github.AppInstallationID, []byte(config.Github.AppPrivateKey))
			if err != nil {
				log.Fatalf("authentication failed: using private key of size %d (%s...): %v", len(config.Github.AppPrivateKey), strings.Split(config.Github.AppPrivateKey, "\n")[0], err)
			}
		}

		t = tr
	}

	transport := &httpcache.Transport{
		Transport:           t,
		Cache:               httpcache.NewMemoryCache(),
		MarkCachedResponses: true,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	if config.Github.APIURL == "api.github.com" {
		client = github.NewClient(httpClient)
	} else {
		client, err = github.NewEnterpriseClient(config.Github.APIURL, config.Github.APIURL, httpClient)
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
