# github-actions-exporter
github-actions-exporter for prometheus

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/spendeskplatform/github-actions-exporter)
![Docker Pulls](https://img.shields.io/docker/pulls/spendeskplatform/github-actions-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/Spendesk/github-actions-exporter)](https://goreportcard.com/report/github.com/Spendesk/github-actions-exporter)

Container image : https://hub.docker.com/repository/docker/spendeskplatform/github-actions-exporter

## Information
If you want to monitor a public repository, you must put the public_repo option in the repo scope of your github token.

## Options
| Name | Flag | Env vars | Default | Description |
|---|---|---|---|---|
| Github Token | github_token, gt | GITHUB_TOKEN | - | Personnal Access Token |
| Github Refresh | github_refresh, gr | GITHUB_REFRESH | 30 | Refresh time Github Actions status in sec |
| Github Organizations | github_orgas, go | GITHUB_ORGAS | - | List all organizations you want get informations. Format \<orga1>,\<orga2>,\<orga3> (like test1,test2) |
| Github Repos | github_repos, grs | GITHUB_REPOS | - | List all repositories you want get informations. Format \<orga>/\<repo>,\<orga>/\<repo2>,\<orga>/\<repo3> (like test/test) |
| Exporter port | port, p | PORT | 9999 | Exporter port |
| Github Api URL | github_api_url, url | GITHUB_API_URL | api.github.com | Github API URL (primarily for Github Enterprise usage) |

## Exported stats

### github_workflow_run_status
Gauge type

**Result possibility**

| ID | Description |
|---|---|
| 0 | Failure |
| 1 | Success |
| 2 | Skipped |
| 3 | In Progress |

**Fields**

| Name | Description |
|---|---|
| event | Event type like push/pull_request/...|
| head_branch | Branch name |
| head_sha | Commit ID |
| node_id | Node ID (github actions) (mandatory ??) |
| repo | Repository like \<org>/\<repo> |
| run_number | Build id for the repo (incremental id => 1/2/3/4/...) |
| workflow_id | Workflow ID |
| workflow | Workflow Name |
| status | Workflow status (completed/in_progress) |

### github_workflow_run_duration_ms
Gauge type

**Result possibility**

| Gauge | Description |
|---|---|
| milliseconds | Number of milliseconds that a specific workflow run took time to complete. |

**Fields**

| Name | Description |
|---|---|
| event | Event type like push/pull_request/...|
| head_branch | Branch name |
| head_sha | Commit ID |
| node_id | Node ID (github actions) (mandatory ??) |
| repo | Repository like \<org>/\<repo> |
| run_number | Build id for the repo (incremental id => 1/2/3/4/...) |
| workflow_id | Workflow ID |
| workflow | Workflow Name |
| status | Workflow status (completed/in_progress) |

### github_job
> :warning: **This is a duplicate of the `github_workflow_run_status` metric that will soon be deprecated, do not use anymore.**

### github_runner_status
Gauge type
(If you have self hosted runner)

**Result possibility**

| ID | Description |
|---|---|
| 0 | Offline |
| 1 | Online |

**Fields**

| Name | Description |
|---|---|
| id | Runner id (incremental id) |
| name | Runner name |
| os | Operating system (linux/macos/windows) |
| repo | Repository like \<org>/\<repo> |
| status | Runner status (online/offline) |

### github_runner_organization_status
Gauge type
(If you have self hosted runner for an organization)

**Result possibility**

| ID | Description |
|---|---|
| 0 | Offline |
| 1 | Online |

**Fields**

| Name | Description |
|---|---|
| id | Runner id (incremental id) |
| name | Runner name |
| os | Operating system (linux/macos/windows) |
| orga | Organization name |
| status | Runner status (online/offline) |

### github_workflow_usage_seconds
Gauge type
(If you have private repositories that use GitHub-hosted runners)

**Result possibility**

| Gauge | Description |
|---|---|
| seconds | Number of billable seconds used by a specific workflow during the current billing cycle. |

**Fields**

| Name | Description |
|---|---|
| id | Workflow id (incremental id) |
| node_id | Node ID (github actions) |
| name | workflow name |
| os | Operating system (linux/macos/windows) |
| repo | Repository like \<org>/\<repo> |
| status | Workflow status |

Example:

```
# HELP github_workflow_usage Number of billable seconds used by a specific workflow during the current billing cycle. Any job re-runs are also included in the usage. Only apply to workflows in private repositories that use GitHub-hosted runners.
# TYPE github_workflow_usage gauge
github_workflow_usage_seconds{id="2862037",name="Create Release",node_id="MDg6V29ya2Zsb3cyODYyMDM3",repo="xxx/xxx",state="active",os="UBUNTU"} 706.609
```


## Github Token configuration

Scopes needed configuration for the Github token

```
repo
  - repo:status
  - repo_deployment
  - public_repo

admin:org
  - write:org
  - read:org
```
