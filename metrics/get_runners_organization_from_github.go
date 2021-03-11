package metrics

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"

	"github-actions-exporter/config"

	"github.com/prometheus/client_golang/prometheus"

	"golang.org/x/oauth2"
)

var (
	RunnersOrganizationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_organization_status",
			Help: "runner status",
		},
		[]string{"organization", "os", "name", "id"},
	)
)

// GetRunnersOrganizationFromGithub - get runners status
func GetRunnersOrganizationFromGithub() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Github.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.ListOptions{PerPage: 10}

	for {
		for _, orga := range config.Github.Organizations {
			for {
				runners, resp, err := client.Actions.ListOrganizationRunners(ctx, orga, opt)
				if err != nil {
					log.Fatal(err)
				}
			        for _, r := range runners.Runners {
					if strings.ToLower(r.GetStatus()) == "online" {
						RunnersOrganizationGauge.WithLabelValues(orga, r.GetOS(), r.GetName(), strconv.Itoa(int(r.GetID()))).Set(1)
					} else {
						RunnersOrganizationGauge.WithLabelValues(orga, r.GetOS(), r.GetName(), strconv.Itoa(int(r.GetID()))).Set(0)
					}
				}
				if resp.NextPage == 0 {
					break
			        }
				opt.Page = resp.NextPage
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
