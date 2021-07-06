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
	runnersEnterpriseGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_enterprise_status",
			Help: "runner status",
		},
		[]string{"os", "name", "id"},
	)
)

func getRunnersEnterpriseFromGithub() {
	for {
		runners, _, err := client.Enterprise.ListRunners(context.Background(), config.EnterpriseName, nil)
		if err != nil {
			log.Printf("Enterprise.ListRunners error: %s", err.Error())
		} else {
			for _, runner := range runners.Runners {
				var integerStatus float64
				if integerStatus = 0; runner.GetStatus() == "online" {
					integerStatus = 1
				}
				runnersEnterpriseGauge.WithLabelValues(*runner.OS, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(integerStatus)
			}
		}
		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
