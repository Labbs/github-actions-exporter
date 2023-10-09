package metrics

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spendesk/github-actions-exporter/pkg/config"

	"github.com/google/go-github/v45/github"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	runnersGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"repo", "os", "name", "id", "busy"},
	)
)

func getAllRepoRunners(owner string, repo string) []*github.Runner {
	var runners []*github.Runner
	opt := &github.ListOptions{PerPage: 200}

	for {
		resp, rr, err := client.Actions.ListRunners(context.Background(), owner, repo, opt)
		if rl_err, ok := err.(*github.RateLimitError); ok {
			log.Printf("ListRunners ratelimited. Pausing until %s", rl_err.Rate.Reset.Time.String())
			time.Sleep(time.Until(rl_err.Rate.Reset.Time))
			continue
		} else if err != nil {
			log.Printf("ListRunners error for repo %s: %s", repo, err.Error())
			return nil
		}

		runners = append(runners, resp.Runners...)
		if rr.NextPage == 0 {
			break
		}
		opt.Page = rr.NextPage
	}

	return runners
}

// getRunnersFromGithub - return information about runners and their status for a specific repo
func getRunnersFromGithub() {
	for {
		for _, repo := range repositories {
			r := strings.Split(repo, "/")

			runners := getAllRepoRunners(r[0], r[1])
			for _, runner := range runners {
				if runner.GetStatus() == "online" {
					runnersGauge.WithLabelValues(repo, *runner.OS, *runner.Name, strconv.FormatInt(runner.GetID(), 10), strconv.FormatBool(runner.GetBusy())).Set(1)
				} else {
					runnersGauge.WithLabelValues(repo, *runner.OS, *runner.Name, strconv.FormatInt(runner.GetID(), 10), strconv.FormatBool(runner.GetBusy())).Set(0)
				}
			}
		}

		time.Sleep(time.Duration(config.Github.Refresh) * time.Second)
		runnersGauge.Reset()
	}
}
