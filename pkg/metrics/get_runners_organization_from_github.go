package metrics

import (
	"context"
	"github-actions-exporter/pkg/config"
	"log"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	runnersOrganizationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_organization_status",
			Help: "runner status",
		},
		[]string{"organization", "os", "status", "name", "id"},
	)
)

func getRunnersOrganizationFromGithub() {
	for {
		for _, orga := range config.Github.Organizations.Value() {
			resp, _, err := client.Actions.ListOrganizationRunners(context.Background(), orga, nil)
			if err != nil {
				log.Printf("ListOrganizationRunners error for %s: %s", orga, err.Error())
			} else {
				for _, runner := range resp.Runners {
					if runner.GetStatus() == "online" {
						runnersOrganizationGauge.WithLabelValues(orga, *runner.OS, *runner.Status, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(1)
					} else {
						runnersOrganizationGauge.WithLabelValues(orga, *runner.OS, *runner.Status, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(0)
					}
				}
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
