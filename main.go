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

var version = "v1.0"

var (
	runners = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"repo", "os", "status", "name", "id"},
	)

	jobs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_job",
			Help: "job status",
		},
		[]string{"repo", "id", "node_id", "head_branch", "head_sha", "run_number", "event", "status"},
	)
)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/metrics")
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("starting exporter with port %v", config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

// init prometheus metrics
func init() {
	prometheus.MustRegister(runners)
	prometheus.MustRegister(jobs)
}

func getRunnersFromGithub() {
	client := &http.Client{}

	for {
		for _, repo := range config.Github.Repositories {
			var p []runner
			req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/runners", nil)
			req.Header.Set("Authorization", "token "+config.Github.Token)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			err = json.NewDecoder(resp.Body).Decode(&p)
			if err != nil {
				log.Fatal(err)
			}
			for _, r := range p {
				if r.Status == "online" {
					runners.WithLabelValues(repo, r.OS, r.Status, r.Name, strconv.Itoa(r.ID)).Set(1)
				} else {
					runners.WithLabelValues(repo, r.OS, r.Status, r.Name, strconv.Itoa(r.ID)).Set(0)
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
				}
				jobs.WithLabelValues(repo, strconv.Itoa(r.ID), r.NodeID, r.HeadBranch, r.HeadSha, strconv.Itoa(r.RunNumber), r.Event, r.Status).Set(s)
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
