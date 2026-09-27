package jira

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/kernel/plugin"
	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

type apiFactory func(jiraConnection, string) (jiraapi.JiraAPI, error)

type jiraConnection struct {
	URL string
}

type Plugin struct {
	environment env.Env
	newAPI      apiFactory
}

func New(environment env.Env) plugin.RemotePlugin {
	return NewService(environment)
}

func NewService(environment env.Env) *Plugin {
	return newPlugin(environment, defaultAPIFactory(environment))
}

func NewWithAPI(environment env.Env, api jiraapi.JiraAPI) *Plugin {
	return newPlugin(environment, func(jiraConnection, string) (jiraapi.JiraAPI, error) { return api, nil })
}

func newPlugin(environment env.Env, factory apiFactory) *Plugin {
	return &Plugin{environment: environment, newAPI: factory}
}

func (p *Plugin) Name() string {
	return "jira"
}

type CheckRequest struct {
	URL     string
	PAT     string
	Project string
}

type CheckResult struct {
	BaseURL           string
	Version           string
	Project           string
	ProjectIssueCount int
}

func (r CheckResult) HasProject() bool {
	return r.Project != ""
}

func (p *Plugin) Check(ctx context.Context, request CheckRequest) (CheckResult, error) {
	if request.PAT == "" {
		return CheckResult{}, fmt.Errorf("--pat is required")
	}
	connection, err := newJiraConnection(request.URL)
	if err != nil {
		return CheckResult{}, err
	}
	api, err := p.newAPI(connection, request.PAT)
	if err != nil {
		return CheckResult{}, err
	}
	if err := api.ValidateConnection(ctx); err != nil {
		return CheckResult{}, err
	}
	serverInfo, err := api.ServerInfo(ctx)
	if err != nil {
		return CheckResult{}, err
	}
	baseURL := serverInfo.BaseURL()
	result := CheckResult{
		BaseURL: baseURL.String(),
		Version: serverInfo.Version().String(),
		Project: request.Project,
	}
	if request.Project != "" {
		searchResult, err := api.SearchIssues(ctx, projectIssuesJQL(request.Project))
		if err != nil {
			return CheckResult{}, err
		}
		result.ProjectIssueCount = searchResult.Total
	}
	return result, nil
}

type PullRequest struct {
	URL       string
	PAT       string
	Project   string
	IssueKeys []string
}

type PullResult struct {
	IssueCount int
}

func (p *Plugin) Pull(ctx context.Context, request PullRequest) (PullResult, error) {
	if request.PAT == "" {
		return PullResult{}, fmt.Errorf("--pat is required")
	}
	if request.Project == "" {
		return PullResult{}, fmt.Errorf("--project is required")
	}
	if len(request.IssueKeys) == 0 {
		return PullResult{}, fmt.Errorf("--issues is required")
	}
	connection, err := newJiraConnection(request.URL)
	if err != nil {
		return PullResult{}, err
	}
	api, err := p.newAPI(connection, request.PAT)
	if err != nil {
		return PullResult{}, err
	}
	if err := api.ValidateConnection(ctx); err != nil {
		return PullResult{}, err
	}
	searchResult, err := api.SearchIssues(ctx, selectedProjectIssuesJQL(request.Project, request.IssueKeys))
	if err != nil {
		return PullResult{}, err
	}
	return PullResult{IssueCount: searchResult.Total}, nil
}

func projectIssuesJQL(project string) string {
	return fmt.Sprintf("project = %s", project)
}

func selectedProjectIssuesJQL(project string, issueKeys []string) string {
	return fmt.Sprintf("%s AND key in (%s)", projectIssuesJQL(project), strings.Join(issueKeys, ", "))
}

func newJiraConnection(jiraURL string) (jiraConnection, error) {
	if jiraURL == "" {
		return jiraConnection{}, fmt.Errorf("--url is required")
	}
	baseURL, err := url.Parse(jiraURL)
	if err != nil {
		return jiraConnection{}, fmt.Errorf("--url: %w", err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return jiraConnection{}, fmt.Errorf("--url must be absolute")
	}
	return jiraConnection{URL: jiraURL}, nil
}

func defaultAPIFactory(environment env.Env) apiFactory {
	return func(connection jiraConnection, pat string) (jiraapi.JiraAPI, error) {
		baseURL, err := url.Parse(connection.URL)
		if err != nil {
			return nil, err
		}
		return jiraapi.NewAuthenticatedRestClient(baseURL, environment.HTTPClient, pat)
	}
}
