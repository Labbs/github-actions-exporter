package metrics

import (
	"context"
	"github-actions-exporter/pkg/config"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	runnersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"repo", "os", "status", "name", "id"},
	)
)

func getRunnersFromGithub() {
	for {
		for _, repo := range config.Github.Repositories.Value() {
			r := strings.Split(repo, "/")
			resp, _, err := client.Actions.ListRunners(context.Background(), r[0], r[1], nil)
			if err != nil {
				log.Printf("ListRunners error for %s: %s", repo, err.Error())
			} else {
				for _, runner := range resp.Runners {
					if runner.GetStatus() == "online" {
						runnersGauge.WithLabelValues(repo, *runner.OS, *runner.Status, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(1)
					} else {
						runnersGauge.WithLabelValues(repo, *runner.OS, *runner.Status, *runner.Name, strconv.FormatInt(runner.GetID(), 10)).Set(0)
					}
				}
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
