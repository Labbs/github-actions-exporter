## Change of organization
The project was migrated to another organization because it was no longer maintained.
In the coming weeks, changes will be made on various points (code and container image) and also an update of all the libraries.
If there are feature requests, you can make them in the issues. For a better understanding of the requested features, having a detailed description or example will be greatly appreciated.

# github-actions-exporter
github-actions-exporter for prometheus

![Release pipeline](https://github.com/Labbs/github-actions-exporter/actions/workflows/release.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Spendesk/github-actions-exporter)](https://goreportcard.com/report/github.com/Spendesk/github-actions-exporter)

Container image : https://github.com/Labbs/github-actions-exporter/pkgs/container/github-actions-exporter

Helm Chart :
```
 helm chart repo: https://labbs.github.io/helm-charts/
 chart: github-actions-exporter
 version: "0.1.4"
```

## Information
If you want to monitor a public repository, you must put the public_repo option in the repo scope of your github token or Github App Authentication.

## Authentication 

Authentication can either via a Github Token or the Github App Authentication 3 parameters. When installing via the Helm Chart the authentication is provided via a secret.



## Options
| Name | Flag | Env vars | Default | Description |
|---|---|---|---|---|
| Github Token | github_token, gt | GITHUB_TOKEN | - | Personnel Access Token |
| Github App Id | app_id, gai | GITHUB_APP_ID |  | Github App Authentication App Id |
| Github App Installation Id | app_installation_id, gii | GITHUB_APP_INSTALLATION_ID | - | Github App Authentication Installation Id |
| Github App Private Key | app_private_key, gpk | GITHUB_APP_PRIVATE_KEY | - | Github App Authentication Private Key |
| Github Refresh | github_refresh, gr | GITHUB_REFRESH | 30 | Refresh time Github Actions status in sec |
| Github Organizations | github_orgas, go | GITHUB_ORGAS | - | List all organizations you want get informations. Format \<orga1>,\<orga2>,\<orga3> (like test1,test2) |
| Github Repositories List File | repo_list_file | REPO_LIST_FILE | - | [Optional] List all repositories you want get informations. Multiline format, check `example_repo.txt`. Defaults to all repositories owned by the organizations. |
| Exporter port | port, p | PORT | 9999 | Exporter port |
| Github Api URL | github_api_url, url | GITHUB_API_URL | api.github.com | Github API URL (primarily for Github Enterprise usage) |
| Github Enterprise Name | enterprise_name | ENTERPRISE_NAME | "" | Enterprise name. Needed for enterprise endpoints (/enterprises/{ENTERPRISE_NAME}/*). Currently used to get Enterprise level tunners status |
| Fields to export | export_fields | EXPORT_FIELDS | repo,id,node_id,head_branch,head_sha,run_number,workflow_id,workflow,event,status | A comma separated list of fields for workflow metrics that should be exported |

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
| 4 | Queued |
| 5 | Cancelled |

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
| busy | Runner busy or not (true/false) |

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
| busy | Runner busy or not (true/false) |

### github_runner_enterprise_status
Gauge type
(If you have self hosted runner for an enterprise)

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

## Setting up authentication with GitHub API

There are two ways for github-actions-exporter to authenticate with the GitHub API (only 1 can be configured at a time however):

1. Using a GitHub App (not supported when you use Github Enterprise )
2. Using a Personal Access Token

Functionality wise, there isn't much of a difference between the 2 authentication methods. The primarily benefit of authenticating via a GitHub App is an [increased API quota](https://docs.github.com/en/developers/apps/rate-limits-for-github-apps).

If you are deploying the solution for a GitHub Enterprise Server environment you are able to [configure your rate limiting settings](https://docs.github.com/en/enterprise-server@3.0/admin/configuration/configuring-rate-limits) making the main benefit irrelevant. If you're deploying the solution for a GitHub Enterprise Cloud or regular GitHub environment and you run into rate limiting issues, consider deploying the solution using the GitHub App authentication method instead.

### Deploying using GitHub App Authentication

You can create a GitHub App for either your account or any organization. If you want to create a GitHub App for your account, open the following link to the creation page, enter any unique name in the "GitHub App name" field, and hit the "Create GitHub App" button at the bottom of the page.

- [Create GitHub Apps on your account](https://github.com/settings/apps/new?url=http://github.com/github-actions-exporter/github-actions-exporter&webhook_active=false&public=false&administration=write&actions=read)

If you want to create a GitHub App for your organization, replace the `:org` part of the following URL with your organization name before opening it. Then enter any unique name in the "GitHub App name" field, and hit the "Create GitHub App" button at the bottom of the page to create a GitHub App.

- [Create GitHub Apps on your organization](https://github.com/organizations/:org/settings/apps/new?url=http://github.com/github-actions-exporter/github-actions-exporter&webhook_active=false&public=false&administration=write&organization_self_hosted_runners=write&actions=read)

### Github Token configuration

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

### Authentication Errors

#### Invalid Github Token 
 if token is invalid then `401 Bad credentials` will be returned on github API error and displayed in an error message. 

#### Invalid Github App configuration 
 if the app id or app installation id value is incorrect then messages like the following are displayed:
 ```
 could not refresh installation id 12345678's token: request &{Method:POST URL:https://api.github.com/app/installations/12345678/access_tokens
 ``` 

 if the github_app_private_key is incorrect then errors like the following are displayed. 
 ```
  Error: Client creation failed.authentication failed: could not parse private key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private ke
 ```

###  Secret actions-exporter 

In the kubernetes deployment authentication is passed via a kubernetes secret: 

```
kind: Secret
apiVersion: v1
metadata:
  name: actions-exporter
  namespace: github-actions-exporter
type: Opaque
data:
  github_token: AAAAAA
#  github_app_id: BBBBBB
#  github_app_installation_id: CCCCCCCCC
#  github_app_private_key: DDDDDDD
```

Or more probably using an external secret manager. Here is an example of using External Secrets with the EKS Secret Manager to define the authentication in a secret: 

```
apiVersion: 'kubernetes-client.io/v1'
kind: ExternalSecret
metadata:
  name: actions-exporter
  namespace: github-actions-exporter
spec:
  backendType: secretsManager
  data:
 #   - key: MySecretManagerKey
 #     name: github_token
 #     property: github_token
    - key: MySecretManagerKey
      name: github_app_id
      property: github_app_id
    - key: MySecretManagerKey
      name: github_app_installation_id
      property: github_app_installation_id
  # separate plaintext aws secret needed for ssh key
    - key: MySecretManagerKeyPrivateKey
      name: github_app_private_key
```
