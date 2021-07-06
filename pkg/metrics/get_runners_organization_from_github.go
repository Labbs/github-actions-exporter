package metrics

import (
	"context"
	"github-actions-exporter/pkg/config"
	"log"
	"strconv"
	"time"

	"github.com/google/go-github/v33/github"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	runnersOrganizationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_organization_status",
			Help: "runner status",
		},
		[]string{"organization", "os", "name", "id"},
	)
)

// getRunnersOrganizationFromGithub - return information about runners and their status for an organization
func getRunnersOrganizationFromGithub() {
	for {
		for _, orga := range config.Github.Organizations.Value() {
			opt := &github.ListOptions{PerPage: 10}
			for {
				resp, rr, err := client.Actions.ListOrganizationRunners(context.Background(), orga, opt)
				if err != nil {
					log.Printf("ListOrganizationRunners error for %s: %s", orga, err.Error())
				} else {
					for _, runner := range resp.Runners {
						if runner.GetStatus() == "online" {
							runnersOrganizationGauge.WithLabelValues(orga, *runner.OS, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(1)
						} else {
							runnersOrganizationGauge.WithLabelValues(orga, *runner.OS, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(0)
						}
					}
					if rr.NextPage == 0 {
						break
					}
					opt.Page = rr.NextPage
				}
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
		runnersOrganizationGauge.Reset()
	}
}
