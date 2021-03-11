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
	RunnersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"repo", "os", "name", "id"},
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

// GetRunnersFromGithub - get runners status
func GetRunnersFromGithub() {
	client := &http.Client{}

	for {
		for _, repo := range config.Github.Repositories {
			var p runners
			req, _ := http.NewRequest("GET", "https://"+config.Github.ApiUrl+"/repos/"+repo+"/actions/runners", nil)
			req.Header.Set("Authorization", "token "+config.Github.Token)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				log.Fatalf("the status code returned by the server for runners in repo %s is different from 200: %d", repo, resp.StatusCode)
			}
			err = json.NewDecoder(resp.Body).Decode(&p)
			if err != nil {
				log.Fatal(err)
			}
			for _, r := range p.Runners {
				if r.Status == "online" {
					RunnersGauge.WithLabelValues(repo, r.OS, r.Name, strconv.Itoa(r.ID)).Set(1)
				} else {
					RunnersGauge.WithLabelValues(repo, r.OS, r.Name, strconv.Itoa(r.ID)).Set(0)
				}

			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
