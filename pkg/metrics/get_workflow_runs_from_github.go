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
	workflowRunStatusDeprecatedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_job",
			Help: "Workflow run status, old name and duplicate of github_workflow_run_status that will soon be deprecated",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"},
	)
	workflowRunStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_status",
			Help: "Workflow run status",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"},
	)
	workflowRunDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_run_duration_ms",
			Help: "Workflow run duration (in milliseconds)",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "workflow_id", "workflow", "event", "status"},
	)
)

// getWorkflowRunsFromGithub - return informations and status about a worflow
func getWorkflowRunsFromGithub() {
	for {
		for _, repo := range config.Github.Repositories.Value() {
			r := strings.Split(repo, "/")
			resp, _, err := client.Actions.ListRepositoryWorkflowRuns(context.Background(), r[0], r[1], nil)
			if err != nil {
				log.Printf("ListRepositoryWorkflowRuns error for %s: %s", repo, err.Error())
			} else {
				for _, run := range resp.WorkflowRuns {
					var s float64 = 0
					if run.GetConclusion() == "success" {
						s = 1
					} else if run.GetConclusion() == "skipped" {
						s = 2
					} else if run.GetConclusion() == "in_progress" {
						s = 3
					} else if run.GetConclusion() == "queued" {
						s = 4
					}
					workflowRunStatusGauge.WithLabelValues(repo, strconv.FormatInt(*run.ID, 10), *run.NodeID, *run.HeadBranch, *run.HeadSHA, strconv.Itoa(*run.RunNumber), strconv.FormatInt(*run.WorkflowID, 10), *workflows[repo][*run.WorkflowID].Name, *run.Event, *run.Status).Set(s)
					workflowRunStatusDeprecatedGauge.WithLabelValues(repo, strconv.FormatInt(*run.ID, 10), *run.NodeID, *run.HeadBranch, *run.HeadSHA, strconv.Itoa(*run.RunNumber), strconv.FormatInt(*run.WorkflowID, 10), *workflows[repo][*run.WorkflowID].Name, *run.Event, *run.Status).Set(s)

					resp, _, err := client.Actions.GetWorkflowRunUsageByID(context.Background(), r[0], r[1], *run.ID)
					if err != nil {
						log.Printf("GetWorkflowRunUsageByID error for %s: %s", repo, err.Error())
					} else {
						workflowRunDurationGauge.WithLabelValues(repo, strconv.FormatInt(*run.ID, 10), *run.NodeID, *run.HeadBranch, *run.HeadSHA, strconv.Itoa(*run.RunNumber), strconv.FormatInt(*run.WorkflowID, 10), *workflows[repo][*run.WorkflowID].Name, *run.Event, *run.Status).Set(float64(resp.GetRunDurationMS()))
					}
				}
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
