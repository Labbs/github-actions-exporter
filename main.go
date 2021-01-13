/*
Main application package
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"

	"github-actions-exporter/config"
)

var version = "v1.2"

var (
	runnersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"repo", "os", "status", "name", "id"},
	)

	jobsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_job",
			Help: "job status",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "event", "status"},
	)
	workflowBillGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_workflow_usage_seconds",
			Help: "Number of billable seconds used by a specific workflow during the current billing cycle. Any job re-runs are also included in the usage. Only apply to workflows in private repositories that use GitHub-hosted runners.",
		},
		[]string{"repo", "id", "node_id", "name", "state", "os"},
	)
)

type runners struct {
	TotalCount int      `json:"total_count"`
	Runners    []runner `json:"runners"`
}

type runner struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	OS     string `json:"os"`
	Status string `json:"status"`
}

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

type workflowsReturn struct {
	TotalCount int        `json:"total_count"`
	Workflows  []workflow `json:"workflows"`
}

type workflow struct {
	ID     int    `json:"id"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	State  string `json:"state"`
}

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

// main init configuration
func main() {
	app := cli.NewApp()
	app.Name = "github-actions-exporter"
	app.Flags = config.NewContext()
	app.Action = runWeb
	app.Version = version

	app.Run(os.Args)
}

// runWeb start http server
func runWeb(ctx *cli.Context) {
	go getRunnersFromGithub()
	go getJobsFromGithub()
	go getBillableFromGithub()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/metrics")
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("starting exporter with port %v", config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

// init prometheus metrics
func init() {
	prometheus.MustRegister(runnersGauge)
	prometheus.MustRegister(jobsGauge)
	prometheus.MustRegister(workflowBillGauge)
}

func getRunnersFromGithub() {
	client := &http.Client{}

	for {
		for _, repo := range config.Github.Repositories {
			var p runners
			req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/runners", nil)
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
			for _, r := range p.Runners {
				if r.Status == "online" {
					runnersGauge.WithLabelValues(repo, r.OS, r.Status, r.Name, strconv.Itoa(r.ID)).Set(1)
				} else {
					runnersGauge.WithLabelValues(repo, r.OS, r.Status, r.Name, strconv.Itoa(r.ID)).Set(0)
				}

			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}

func getJobsFromGithub() {
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
				jobsGauge.WithLabelValues(repo, strconv.Itoa(r.ID), r.NodeID, r.HeadBranch, r.HeadSha, strconv.Itoa(r.RunNumber), r.Event, r.Status).Set(s)
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}

func getBillableFromGithub() {
	client := &http.Client{}
	for {
		for _, repo := range config.Github.Repositories {
			var p workflowsReturn
			req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/workflows", nil)
			req.Header.Set("Authorization", "token "+config.Github.Token)
			resp, err := client.Do(req)
			defer resp.Body.Close()
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

			for _, w := range p.Workflows {
				var bill Bill
				req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/workflows/"+strconv.Itoa(w.ID)+"/timing", nil)
				req.Header.Set("Authorization", "token "+config.Github.Token)
				resp2, err := client.Do(req)
				defer resp2.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
				if resp.StatusCode != 200 {
					log.Fatalf("the status code returned by the server is different from 200: %d", resp.StatusCode)
				}
				err = json.NewDecoder(resp2.Body).Decode(&bill)
				if err != nil {
					log.Fatal(err)
				}
				workflowBillGauge.WithLabelValues(repo, strconv.Itoa(w.ID), w.NodeID, w.Name, w.State, "MACOS").Set(bill.Billable.MACOS.TotalMs / 1000)
				workflowBillGauge.WithLabelValues(repo, strconv.Itoa(w.ID), w.NodeID, w.Name, w.State, "WINDOWS").Set(bill.Billable.WINDOWS.TotalMs / 1000)
				workflowBillGauge.WithLabelValues(repo, strconv.Itoa(w.ID), w.NodeID, w.Name, w.State, "UBUNTU").Set(bill.Billable.UBUNTU.TotalMs / 1000)
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * 5 * time.Second)
	}
}
