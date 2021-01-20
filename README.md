# github-actions-exporter
github-actions-exporter for prometheus

## Information
If you want to monitor a public repository, you must put the public_repo option in the repo scope of your github token.

## Options
| Name | Flag | Env vars | Default | Description |
|---|---|---|---|---|
| Github Token | github_token, gt | GITHUB_TOKEN | - | Personnal Access Token |
| Github Refresh | github_refresh, gr | GITHUB_REFRESH | 30 | Refresh time Github Actions status in sec |
| Github Organizations | github_orgas, go | GIHUB_ORGAS | - | List all organizations you want get informations. Format \<orga1>,\<orga2>,\<orga3> (like test1,test2) |
| Github Repos | github_repos, grs | GIHUB_REPOS | - | List all repositories you want get informations. Format \<orga>/\<repo>,\<orga>/\<repo2>,\<orga>/\<repo3> (like test/test) |
| Exporter port | port, p | PORT | 9999 | Exporter port |

## Exported stats

### github_job
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
| status | Workflow status (completed/in_progress) |

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

Es:

```
# HELP github_workflow_usage Number of billable seconds used by a specific workflow during the current billing cycle. Any job re-runs are also included in the usage. Only apply to workflows in private repositories that use GitHub-hosted runners.
# TYPE github_workflow_usage gauge
github_workflow_usage_seconds{id="2862037",name="Create Release",node_id="MDg6V29ya2Zsb3cyODYyMDM3",repo="xxx/xxx",state="active",os="UBUNTU"} 706.609
```