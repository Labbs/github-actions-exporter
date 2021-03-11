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
	WorkflowRunStatusDeprecatedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_job",
			Help: "Workflow run status, old name and duplicate of github_workflow_run_status that will soon be deprecated",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"},
	)
	WorkflowRunStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_status",
			Help: "Workflow run status",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"},
	)
	WorkflowRunDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_duration_ms",
			Help: "Workflow run duration (in milliseconds)",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"},
	)
)

type workflowRunsReturn struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []workflowRun `json:"workflow_runs"`
}

type workflowRunDurationReturn struct {
	RunDurationMs float64 `json:"run_duration_ms"`
}

type workflowRun struct {
	ID         int    `json:"id"`
	NodeID     string `json:"node_id"`
	HeadBranch string `json:"head_branch"`
	HeadSha    string `json:"head_sha"`
	RunNumber  int    `json:"run_number"`
	Event      string `json:"event"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	UpdatedAt  string `json:"updated_at"`
	WorkflowID int    `json:"workflow_id"`
}

func GetWorkflowRunsFromGithub() {
	client := &http.Client{}

	for {
		for _, repo := range config.Github.Repositories {
			var runs workflowRunsReturn
			req, _ := http.NewRequest("GET", "https://"+config.Github.ApiUrl+"/repos/"+repo+"/actions/runs", nil)
			req.Header.Set("Authorization", "token "+config.Github.Token)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != 200 {
				log.Fatalf("the status code returned by the server for runs in repo %s is different from 200: %d", repo, resp.StatusCode)
			}
			err = json.NewDecoder(resp.Body).Decode(&runs)
			if err != nil {
				log.Fatal(err)
			}
			for _, r := range runs.WorkflowRuns {
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
				WorkflowRunStatusGauge.WithLabelValues(repo, strconv.Itoa(r.ID), r.NodeID, r.HeadBranch, r.HeadSha, strconv.Itoa(r.RunNumber), strconv.Itoa(r.WorkflowID), workflows[repo][r.WorkflowID].Name, r.Event, r.Status).Set(s)
				WorkflowRunStatusDeprecatedGauge.WithLabelValues(repo, strconv.Itoa(r.ID), r.NodeID, r.HeadBranch, r.HeadSha, strconv.Itoa(r.RunNumber), strconv.Itoa(r.WorkflowID), workflows[repo][r.WorkflowID].Name, r.Event, r.Status).Set(s)

				var duration workflowRunDurationReturn
				req, _ := http.NewRequest("GET", "https://"+config.Github.ApiUrl+"/repos/"+repo+"/actions/runs/"+strconv.Itoa(r.ID)+"/timing", nil)
				req.Header.Set("Authorization", "token "+config.Github.Token)
				resp, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
				}
				if resp.StatusCode != 200 {
					log.Fatalf("the status code returned by the server for duration of run #%d  in repo %s is different from 200: %d", r.ID, repo, resp.StatusCode)
				}
				err = json.NewDecoder(resp.Body).Decode(&duration)
				if err != nil {
					log.Fatal(err)
				}
				WorkflowRunDurationGauge.WithLabelValues(repo, strconv.Itoa(r.ID), r.NodeID, r.HeadBranch, r.HeadSha, strconv.Itoa(r.RunNumber), strconv.Itoa(r.WorkflowID), workflows[repo][r.WorkflowID].Name, r.Event, r.Status).Set(duration.RunDurationMs)
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
