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
	JobsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_job",
			Help: "job status",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "event", "status"},
	)
)

type jobsReturn struct {
	TotalCount   int   `json:"total_count"`
	WorkflowRuns []job `json:"workflow_runs"`
}

type job struct {
	ID         int    `json:"id"`
	NodeID     string `json:"node_id"`
	HeadBranch string `json:"head_branch"`
	HeadSha    string `json:"head_sha"`
	RunNumber  int    `json:"run_number"`
	Event      string `json:"event"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	UpdatedAt  string `json:"updated_at"`
}

func GetJobsFromGithub() {
	client := &http.Client{}

	for {
		for _, repo := range config.Github.Repositories {
			var p jobsReturn
			req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/runs", nil)
			req.Header.Set("Authorization", "token "+config.Github.Token)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != 200 {
				log.Fatalf("the status code returned by the server is different from 200: %d", resp.StatusCode)
			}
			err = json.NewDecoder(resp.Body).Decode(&p)
			if err != nil {
				log.Fatal(err)
			}
			for _, r := range p.WorkflowRuns {
				var s float64 = 0
				if r.Conclusion == "success" {
					s = 1
				} else if r.Conclusion == "skipped" {
					s = 2
				} else if r.Status == "in_progress" {
					s = 3
				} else if r.Status == "queued" {
					s = 4
				}
				JobsGauge.WithLabelValues(repo, strconv.Itoa(r.ID), r.NodeID, r.HeadBranch, r.HeadSha, strconv.Itoa(r.RunNumber), r.Event, r.Status).Set(s)
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
