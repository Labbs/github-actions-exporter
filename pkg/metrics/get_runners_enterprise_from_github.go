package metrics

import (
	"context"
	"github-actions-exporter/pkg/config"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	runnersEnterpriseGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_enterprise_status",
			Help: "runner status",
		},
		[]string{"os", "name", "status"},
	)
)

func getRunnersEnterpriseFromGithub() {
	runnersEnterpriseGauge.Reset()
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
				var status string
				if status = "idle"; *runner.Busy == true {
					status = "busy"
				}
				runnersEnterpriseGauge.WithLabelValues(*runner.OS, *runner.Name, status).Set(integerStatus)
			}
		}
		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
