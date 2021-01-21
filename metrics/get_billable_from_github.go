package metrics

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github-actions-exporter/config"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	WorkflowBillGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_usage_seconds",
			Help: "Number of billable seconds used by a specific workflow during the current billing cycle. Any job re-runs are also included in the usage. Only apply to workflows in private repositories that use GitHub-hosted runners.",
		},
		[]string{"repo", "id", "node_id", "name", "state", "os"},
	)
)

type UBUNTU struct {
	TotalMs float64 `json:"total_ms"`
}
type MACOS struct {
	TotalMs float64 `json:"total_ms"`
}
type WINDOWS struct {
	TotalMs float64 `json:"total_ms"`
}

type Bill struct {
	Billable Billable `json:"billable"`
}

type Billable struct {
	UBUNTU  UBUNTU  `json:"UBUNTU"`
	MACOS   MACOS   `json:"MACOS"`
	WINDOWS WINDOWS `json:"WINDOWS"`
}

func GetBillableFromGithub() {
	client := &http.Client{}
	for {
		for _, repo := range config.Github.Repositories {
			for k, v := range workflows[repo] {
				var bill Bill
				req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/workflows/"+strconv.Itoa(k)+"/timing", nil)
				req.Header.Set("Authorization", "token "+config.Github.Token)
				resp, err := client.Do(req)
				defer resp.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
				if resp.StatusCode != 200 {
					log.Fatalf("the status code returned by the server is different from 200: %d", resp.StatusCode)
				}
				err = json.NewDecoder(resp.Body).Decode(&bill)
				if err != nil {
					log.Fatal(err)
				}
				WorkflowBillGauge.WithLabelValues(repo, strconv.Itoa(v.ID), v.NodeID, v.Name, v.State, "MACOS").Set(bill.Billable.MACOS.TotalMs / 1000)
				WorkflowBillGauge.WithLabelValues(repo, strconv.Itoa(v.ID), v.NodeID, v.Name, v.State, "WINDOWS").Set(bill.Billable.WINDOWS.TotalMs / 1000)
				WorkflowBillGauge.WithLabelValues(repo, strconv.Itoa(v.ID), v.NodeID, v.Name, v.State, "UBUNTU").Set(bill.Billable.UBUNTU.TotalMs / 1000)
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * 5 * time.Second)
	}
}
