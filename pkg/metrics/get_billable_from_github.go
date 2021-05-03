package metrics

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Spendesk/github-actions-exporter/pkg/config"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	workflowBillGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_usage_seconds",
			Help: "Number of billable seconds used by a specific workflow during the current billing cycle. Any job re-runs are also included in the usage. Only apply to workflows in private repositories that use GitHub-hosted runners.",
		},
		[]string{"repo", "id", "node_id", "name", "state", "os"},
	)
)

// getBillableFromGithub - return billable informations for MACOS, WINDOWS and UBUNTU runners.
func getBillableFromGithub() {
	for {
		for _, repo := range config.Github.Repositories.Value() {
			for k, v := range workflows[repo] {
				r := strings.Split(repo, "/")
				resp, _, err := client.Actions.GetWorkflowUsageByID(context.Background(), r[0], r[1], k)
				if err != nil {
					log.Printf("GetWorkflowUsageByID error for %s: %s", repo, err.Error())
				} else {
					workflowBillGauge.WithLabelValues(repo, strconv.FormatInt(*v.ID, 10), *v.NodeID, *v.Name, *v.State, "MACOS").Set(float64(resp.GetBillable().MacOS.GetTotalMS()) / 1000)
					workflowBillGauge.WithLabelValues(repo, strconv.FormatInt(*v.ID, 10), *v.NodeID, *v.Name, *v.State, "WINDOWS").Set(float64(resp.GetBillable().Windows.GetTotalMS()) / 1000)
					workflowBillGauge.WithLabelValues(repo, strconv.FormatInt(*v.ID, 10), *v.NodeID, *v.Name, *v.State, "UBUNTU").Set(float64(resp.GetBillable().Ubuntu.GetTotalMS()) / 1000)
				}
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * 5 * time.Second)
	}
}
