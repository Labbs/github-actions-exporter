package metrics

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v38/github"

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
			s := make(map[int64]github.Workflow)
			opt := &github.ListOptions{PerPage: 30}

			for {
				resp, rr, err := client.Actions.ListWorkflows(context.Background(), r[0], r[1], opt)
				if err != nil {
					log.Printf("ListWorkflows error for %s: %s", repo, err.Error())
					break
				}
				for _, w := range resp.Workflows {
					s[*w.ID] = *w
				}

				if rr.NextPage == 0 {
					break
				}
				opt.Page = rr.NextPage
			}

			ww[repo] = s

		}

		workflows = ww

		time.Sleep(time.Duration(60) * time.Second)
	}
}
