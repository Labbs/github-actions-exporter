/*
Main application package
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"

	"github-actions-exporter/config"
	"github-actions-exporter/metrics"
)

var version = "v1.4"

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
	go metrics.WorkflowsCache()
	go metrics.GetRunnersFromGithub()
	go metrics.GetRunnersOrganizationFromGithub()
	go metrics.GetJobsFromGithub()
	go metrics.GetBillableFromGithub()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/metrics")
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("starting exporter with port %v", config.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), nil))
}

// init prometheus metrics
func init() {
	prometheus.MustRegister(metrics.RunnersGauge)
	prometheus.MustRegister(metrics.RunnersOrganizationGauge)
	prometheus.MustRegister(metrics.JobsGauge)
	prometheus.MustRegister(metrics.WorkflowBillGauge)
}
