package jira

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

type apiFactory func(JiraURL, string) (jiraapi.JiraAPI, error)

type JiraURL struct {
	value url.URL
}

func (u JiraURL) String() string {
	return u.value.String()
}

type Plugin struct {
	environment env.Env
	newAPI      apiFactory
}

func New(environment env.Env) *Plugin {
	return newPlugin(environment, defaultAPIFactory(environment))
}

func NewWithAPI(environment env.Env, api jiraapi.JiraAPI) *Plugin {
	return newPlugin(environment, func(JiraURL, string) (jiraapi.JiraAPI, error) { return api, nil })
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
	if err := request.validate(); err != nil {
		return CheckResult{}, err
	}
	api, _, err := p.validatedAPI(ctx, request.URL, request.PAT)
	if err != nil {
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
		searchResult, err := api.SearchIssues(ctx, jiraapi.JiraIssueSearchRequest{JQL: projectIssuesJQL(request.Project)})
		if err != nil {
			return CheckResult{}, err
		}
		result.ProjectIssueCount = searchResult.Total
	}
	return result, nil
}

const defaultPullLimit = 5

type PullRequest struct {
	URL       string
	PAT       string
	Project   string
	IssueKeys []string
	Limit     int
}

type PullResult struct {
	WorkItems []kernel.WorkItem
}

func (p *Plugin) Pull(ctx context.Context, request PullRequest) (PullResult, error) {
	if err := request.validate(); err != nil {
		return PullResult{}, err
	}
	api, jiraURL, err := p.validatedAPI(ctx, request.URL, request.PAT)
	if err != nil {
		return PullResult{}, err
	}
	searchResult, err := api.SearchIssues(ctx, jiraapi.JiraIssueSearchRequest{
		JQL:        request.jql(),
		Fields:     []string{"summary", "issuetype", "project"},
		MaxResults: request.limit(),
	})
	if err != nil {
		return PullResult{}, err
	}
	return pullResultFromIssues(jiraURL.String(), searchResult.Issues), nil
}

func (r CheckRequest) validate() error {
	if r.PAT == "" {
		return fmt.Errorf("--pat is required")
	}
	return nil
}

func (r PullRequest) validate() error {
	if r.PAT == "" {
		return fmt.Errorf("--pat is required")
	}
	if r.Project == "" {
		return fmt.Errorf("--project is required")
	}
	if r.Limit < 0 {
		return fmt.Errorf("--limit must be positive")
	}
	if len(r.IssueKeys) > 0 && r.Limit > 0 && r.Limit < len(r.IssueKeys) {
		return fmt.Errorf("--limit must be at least the number of selected issues")
	}
	return nil
}

func (r PullRequest) jql() string {
	if len(r.IssueKeys) > 0 {
		return selectedProjectIssuesJQL(r.Project, r.IssueKeys)
	}
	return projectIssuesJQL(r.Project)
}

func (r PullRequest) limit() int {
	if r.Limit > 0 {
		return r.Limit
	}
	if len(r.IssueKeys) > 0 {
		return len(r.IssueKeys)
	}
	return defaultPullLimit
}

func (p *Plugin) validatedAPI(ctx context.Context, rawURL string, pat string) (jiraapi.JiraAPI, JiraURL, error) {
	jiraURL, err := parseJiraURL(rawURL)
	if err != nil {
		return nil, JiraURL{}, err
	}
	api, err := p.newAPI(jiraURL, pat)
	if err != nil {
		return nil, JiraURL{}, err
	}
	if err := api.ValidateConnection(ctx); err != nil {
		return nil, JiraURL{}, err
	}
	return api, jiraURL, nil
}

func pullResultFromIssues(remoteURL string, issues []jiraapi.JiraIssue) PullResult {
	workItems := make([]kernel.WorkItem, 0, len(issues))
	for _, issue := range issues {
		workItems = append(workItems, mapIssueWorkItem(remoteURL, issue))
	}
	return PullResult{WorkItems: workItems}
}

func mapIssueWorkItem(remoteURL string, issue jiraapi.JiraIssue) kernel.WorkItem {
	return kernel.WorkItem{
		Remote:   remoteURL,
		Provider: "jira",
		ID:       issue.ID,
		Key:      issue.Key,
		URL:      strings.TrimRight(remoteURL, "/") + "/browse/" + issue.Key,
		Title:    issue.Fields.Summary,
		Fields: []kernel.Field{
			{Name: "project", ProviderKey: "project", Type: "string", Value: issue.Fields.Project.Key, Editable: false},
			{Name: "issueType", ProviderKey: "issuetype", Type: "string", Value: issue.Fields.IssueType.Name, Editable: false},
		},
		Metadata: map[string]string{"jira.baseURL": remoteURL},
	}
}

func projectIssuesJQL(project string) string {
	return fmt.Sprintf("project = %s", project)
}

func selectedProjectIssuesJQL(project string, issueKeys []string) string {
	return fmt.Sprintf("%s AND key in (%s)", projectIssuesJQL(project), strings.Join(issueKeys, ", "))
}

func parseJiraURL(rawURL string) (JiraURL, error) {
	if rawURL == "" {
		return JiraURL{}, fmt.Errorf("--url is required")
	}
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return JiraURL{}, fmt.Errorf("--url: %w", err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return JiraURL{}, fmt.Errorf("--url must be absolute")
	}
	return JiraURL{value: *baseURL}, nil
}

func defaultAPIFactory(environment env.Env) apiFactory {
	return func(jiraURL JiraURL, pat string) (jiraapi.JiraAPI, error) {
		baseURL := jiraURL.value
		return jiraapi.NewAuthenticatedRestClient(&baseURL, environment.HTTPClient, pat)
	}
}
