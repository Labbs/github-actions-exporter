package metrics

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"

	"github-actions-exporter/pkg/config"
)

var (
	workflows map[string]map[int64]github.Workflow
)

// workflowCache - used for limit calls to github api
func workflowCache() {
	for {
		ww := make(map[string]map[int64]github.Workflow)

		for _, repo := range config.Github.Repositories.Value() {
			r := strings.Split(repo, "/")

			resp, _, err := client.Actions.ListWorkflows(context.Background(), r[0], r[1], nil)
			if err != nil {
				log.Printf("ListWorkflows error for %s: %s", repo, err.Error())
			} else {

				s := make(map[int64]github.Workflow)
				for _, w := range resp.Workflows {
					s[*w.ID] = *w
				}

				ww[repo] = s
			}	

		}

		workflows = ww

		time.Sleep(time.Duration(60) * time.Second)
	}
}
