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
	RunnersOrganizationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_organization_status",
			Help: "runner status",
		},
		[]string{"organization", "os", "status", "name", "id"},
	)
)

type runnersOrganization struct {
	TotalCount int      `json:"total_count"`
	Runners    []runner `json:"runners"`
}

// GetRunnersOrganizationFromGithub - get runners status
func GetRunnersOrganizationFromGithub() {
	client := &http.Client{}

	for {
		for _, orga := range config.Github.Organizations {
			var p runnersOrganization
			req, _ := http.NewRequest("GET", "https://api.github.com/orgs/"+orga+"/actions/runners", nil)
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
					RunnersOrganizationGauge.WithLabelValues(orga, r.OS, r.Status, r.Name, strconv.Itoa(r.ID)).Set(1)
				} else {
					RunnersOrganizationGauge.WithLabelValues(orga, r.OS, r.Status, r.Name, strconv.Itoa(r.ID)).Set(0)
				}

			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
	}
}
