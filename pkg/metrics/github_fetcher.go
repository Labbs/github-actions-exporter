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
	repositories []string
	workflows    map[string]map[int64]github.Workflow
)

func getAllReposForOrg(orga string) []string {
	var all_repos []string

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 200,
			Page:    0,
		},
	}
	for {
		repos_page, resp, err := client.Repositories.ListByOrg(context.Background(), orga, opt)
		if err != nil {
			log.Printf("ListByOrg error for %s: %s", orga, err.Error())
			break
		}
		for _, repo := range repos_page {
			all_repos = append(all_repos, *repo.FullName)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return all_repos
}

func getAllWorkflowsForRepo(owner string, repo string) map[int64]github.Workflow {
	res := make(map[int64]github.Workflow)

	opt := &github.ListOptions{
		PerPage: 200,
		Page:    0,
	}

	for {
		workflows_page, resp, err := client.Actions.ListWorkflows(context.Background(), owner, repo, opt)
		if err != nil {
			log.Printf("ListWorkflows error for %s: %s", repo, err.Error())
			return res
		}
		for _, w := range workflows_page.Workflows {
			res[*w.ID] = *w
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return res
}

func periodicGithubFetcher() {
	for {

		// Fetch repositories (if dynamic)
		var repos_to_fetch []string
		if len(config.Github.Repositories.Value()) > 0 {
			repos_to_fetch = config.Github.Repositories.Value()
		} else {
			for _, orga := range config.Github.Organizations.Value() {
				repos_to_fetch = append(repos_to_fetch, getAllReposForOrg(orga)...)
			}
		}
		repositories = repos_to_fetch

		// Fetch workflows
		ww := make(map[string]map[int64]github.Workflow)
		for _, repo := range repositories {
			r := strings.Split(repo, "/")
			ww[repo] = getAllWorkflowsForRepo(r[0], r[1])
			log.Printf("Fetched %d workflows for repository %s", len(ww[repo]), repo)
		}
		workflows = ww

		time.Sleep(time.Duration(60) * time.Second)
	}
}
