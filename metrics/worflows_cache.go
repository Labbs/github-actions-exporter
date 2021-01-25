package metrics

import (
	"encoding/json"
	"github-actions-exporter/config"
	"log"
	"net/http"
	"time"
)

var workflows = map[string]map[int]workflow{}

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

func WorkflowsCache() {
	client := &http.Client{}

	for {
		for _, repo := range config.Github.Repositories {
			var p workflowsReturn
			req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+repo+"/actions/workflows", nil)
			req.Header.Set("Authorization", "token "+config.Github.Token)
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
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

			s := make(map[int]workflow)
			for _, w := range p.Workflows {
				s[w.ID] = w
			}

			workflows[repo] = s
		}
		time.Sleep(time.Duration(60) * time.Second)
	}
}
